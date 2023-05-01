package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
	"github.com/leonsteinhaeuser/example-app/user-service/accessobjects"
)

type userRouter struct {
	userAccessObject accessobjects.User
	log              log.Logger
}

func NewUserRouter(userAccessObject accessobjects.User, log log.Logger) *userRouter {
	return &userRouter{
		userAccessObject: userAccessObject,
		log:              log,
	}
}

func (a *userRouter) Router(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", a.get)
			r.Put("/", a.put)
			r.Delete("/", a.delete)
		})
		r.Get("/", a.list)
		r.Post("/", a.post)
	})
}

func (a *userRouter) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	user, err := a.userAccessObject.Get(ctx, id)
	if err != nil {
		a.log.Error(err).Log("failed to get user")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get user from database", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, user)
	if err != nil {
		a.log.Error(err).Log("failed to write user")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write user", err)
		return
	}
}

func (a *userRouter) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := a.userAccessObject.List(ctx)
	if err != nil {
		a.log.Error(err).Log("failed to list users")
		lib.WriteError(w, http.StatusInternalServerError, "failed to list users", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, users)
	if err != nil {
		a.log.Error(err).Log("failed to write users")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write users", err)
		return
	}
}

func (a *userRouter) put(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	user, err := a.userAccessObject.Get(ctx, id)
	if err != nil {
		a.log.Error(err).Log("failed to get user")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get user from database", err)
		return
	}

	err = lib.ReadJSON(r, user)
	if err != nil {
		a.log.Error(err).Log("failed to read user")
		lib.WriteError(w, http.StatusBadRequest, "failed to read user", err)
		return
	}

	if user.ID.String() != id {
		a.log.Warn().
			Field("urlID", id).
			Field("body", user.ID.String()).
			Log("id in path and body do not match")
		lib.WriteError(w, http.StatusBadRequest, "id in path and body do not match", nil)
		return
	}

	err = a.userAccessObject.Update(ctx, user)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to update user")
		lib.WriteError(w, http.StatusInternalServerError, "failed to update user", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, user)
	if err != nil {
		a.log.Error(err).Log("failed to write user")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write user", err)
		return
	}
}

func (a *userRouter) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	uid, err := uuid.Parse(id)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to parse id")
		lib.WriteError(w, http.StatusBadRequest, "failed to parse id", err)
		return
	}

	user := &lib.User{ID: uid}
	err = a.userAccessObject.Delete(ctx, user)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to delete user")
		lib.WriteError(w, http.StatusInternalServerError, "failed to delete user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *userRouter) post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := &lib.User{}
	err := lib.ReadJSON(r, user)
	if err != nil {
		a.log.Error(err).Log("failed to read user")
		lib.WriteError(w, http.StatusBadRequest, "failed to read user", err)
		return
	}

	err = a.userAccessObject.Create(ctx, user)
	if err != nil {
		a.log.Error(err).Log("failed to create user")
		lib.WriteError(w, http.StatusInternalServerError, "failed to create user", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusCreated, user)
	if err != nil {
		a.log.Error(err).Log("failed to write user")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write user", err)
		return
	}
}
