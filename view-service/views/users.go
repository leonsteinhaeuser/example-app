package views

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

var (
	//go:embed list_user.html
	templateListUsers string
	//go:embed get_user.html
	templateGetUser string
)

type userRouter struct {
	log log.Logger

	listUsersTemplate *template.Template
	getUsersTemplate  *template.Template

	client lib.Client[lib.User]
}

func NewUserRouter(log log.Logger, userClient lib.Client[lib.User]) lib.Router {
	// template for list users
	listUserTemplate, err := template.New("list_user.html").Parse(templateListUsers)
	if err != nil {
		log.Panic(err).Log("failed to parse list_user.html template")
	}
	// template for get user
	getUsersTemplate, err := template.New("get_user.html").Parse(templateGetUser)
	if err != nil {
		log.Panic(err).Log("failed to parse list_user.html template")
	}
	return &userRouter{
		listUsersTemplate: listUserTemplate,
		getUsersTemplate:  getUsersTemplate,
		log:               log,
		client:            userClient,
	}
}

func (a *userRouter) Router(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Get("/list", a.listUsersEndpoint)
		r.Get("/{id}", a.getUserByIDEndpoint)
	})
}

func (a *userRouter) listUsersEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := a.client.List(ctx)
	if err != nil {
		a.log.Error(err).Log("failed to get users")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get users", err)
		return
	}

	err = a.listUsersTemplate.Execute(w, users)
	if err != nil {
		a.log.Error(err).Log("failed to render template")
		lib.WriteError(w, http.StatusInternalServerError, "failed to render template", err)
		return
	}
}

func (a *userRouter) getUserByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	user, err := a.client.GetByID(ctx, id)
	if err != nil {
		a.log.Error(err).Log("failed to get users")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get users", err)
		return
	}

	err = a.getUsersTemplate.Execute(w, user)
	if err != nil {
		a.log.Error(err).Log("failed to render template")
		lib.WriteError(w, http.StatusInternalServerError, "failed to render template", err)
		return
	}
}
