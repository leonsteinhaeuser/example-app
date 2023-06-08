package api

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/internal/log"
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
	log log.Logger
}

func NewViewRouter(log log.Logger) *ViewRouter {
	return &ViewRouter{
		log: log,
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
	title := r.FormValue("title")
	body := r.FormValue("body")
	v.log.Info().Field("title", title).Field("body", body).Log("creating thread")

	// TODO: send request to thread service
}

func (v *ViewRouter) getByIdForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(threadGetByIdForm))
}

func (v *ViewRouter) getByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	v.log.Info().Field("id", id).Log("getting thread by id")
	// TODO: send request to thread service

	err := templateGetByID.Execute(w, thread.Thread{})
	if err != nil {
		v.log.Error(err).Log("error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (v *ViewRouter) list(w http.ResponseWriter, r *http.Request) {
	// TODO: get all threads from thread service

	err := templateList.Execute(w, []thread.Thread{})
	if err != nil {
		v.log.Error(err).Log("error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
