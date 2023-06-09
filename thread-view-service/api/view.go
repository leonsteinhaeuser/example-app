package api

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/public"
	"github.com/leonsteinhaeuser/example-app/thread-service/thread"

	_ "embed"
)

var (
	//go:embed create_form.html
	threadCreateForm string
	//go:embed get_by_id_form.html
	threadGetByIdForm string
	//go:embed get_by_id.html
	threadGetByIdTemplate string
	//go:embed list.html
	threadListTemplate string

	templateGetByID *template.Template
	templateList    *template.Template
)

func init() {
	tgbid, err := template.New("get_by_id").Parse(threadGetByIdTemplate)
	if err != nil {
		panic(err)
	}
	templateGetByID = tgbid

	tl, err := template.New("list").Parse(threadListTemplate)
	if err != nil {
		panic(err)
	}
	templateList = tl
}

type ViewRouter struct {
	threadClient public.Client[thread.Thread]
	log          log.Logger
}

func NewViewRouter(log log.Logger, client public.Client[thread.Thread]) *ViewRouter {
	return &ViewRouter{
		threadClient: client,
		log:          log,
	}
}

func (v *ViewRouter) Router(r chi.Router) {
	r.Route("/thread", func(r chi.Router) {
		r.Get("/create", v.createForm)
		r.Post("/create", v.createHandler)
		r.Get("/get", v.getByIdForm)
		r.Get("/getid", v.getByIdHandler)

		r.Get("/list", v.list)
	})
}

func (v *ViewRouter) createForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(threadCreateForm))
}

func (v *ViewRouter) createHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := v.threadClient.Create(ctx, &thread.Thread{
		Title:    r.FormValue("title"),
		Body:     r.FormValue("body"),
		AuthorID: uuid.New(),
		KeywordIDs: []string{
			uuid.New().String(),
			uuid.New().String(),
			uuid.New().String(),
			uuid.New().String(),
			uuid.New().String(),
			uuid.New().String(),
		},
	})
	if err != nil {
		v.log.Error(err).Log("error creating thread")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (v *ViewRouter) getByIdForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(threadGetByIdForm))
}

func (v *ViewRouter) getByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.FormValue("id")
	// execute get call to thread service
	thread, err := v.threadClient.Get(ctx, id)
	if err != nil {
		v.log.Error(err).Log("error getting thread")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// render template
	err = templateGetByID.Execute(w, thread)
	if err != nil {
		v.log.Error(err).Log("error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (v *ViewRouter) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// execute list call to thread service
	threads, err := v.threadClient.List(ctx)
	if err != nil {
		v.log.Error(err).Log("error listing threads")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templateList.Execute(w, threads)
	if err != nil {
		v.log.Error(err).Log("error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
