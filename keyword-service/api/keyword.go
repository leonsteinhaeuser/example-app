package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/leonsteinhaeuser/example-app/keyword-service/accessobjects"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

type keywordRouter struct {
	keywordAccessObject accessobjects.Keyword
	log                 log.Logger
}

func NewKeywordRouter(keywordAccessObject accessobjects.Keyword, log log.Logger) *keywordRouter {
	return &keywordRouter{
		keywordAccessObject: keywordAccessObject,
		log:                 log,
	}
}

func (a *keywordRouter) Router(r chi.Router) {
	r.Route("/keyword", func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", a.get)
			r.Put("/", a.put)
			r.Delete("/", a.delete)
		})
		r.Get("/", a.list)
		r.Post("/", a.post)
	})
}

func (a *keywordRouter) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	keyword, err := a.keywordAccessObject.Get(ctx, id)
	if err != nil {
		a.log.Error(err).Log("failed to get keyword")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get keywordfrom database", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, keyword)
	if err != nil {
		a.log.Error(err).Log("failed to write keyword")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write keyword", err)
		return
	}
}

func (a *keywordRouter) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	keywords, err := a.keywordAccessObject.List(ctx)
	if err != nil {
		a.log.Error(err).Log("failed to list keywords")
		lib.WriteError(w, http.StatusInternalServerError, "failed to list keywords", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, keywords)
	if err != nil {
		a.log.Error(err).Log("failed to write keywords")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write keywords", err)
		return
	}
}

func (a *keywordRouter) put(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	keyword, err := a.keywordAccessObject.Get(ctx, id)
	if err != nil {
		a.log.Error(err).Log("failed to get keyword")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get keywordfrom database", err)
		return
	}

	err = lib.ReadJSON(r, keyword)
	if err != nil {
		a.log.Error(err).Log("failed to read keyword")
		lib.WriteError(w, http.StatusBadRequest, "failed to read keyword", err)
		return
	}

	if keyword.ID.String() != id {
		a.log.Warn().
			Field("urlID", id).
			Field("body", keyword.ID.String()).
			Log("id in path and body do not match")
		lib.WriteError(w, http.StatusBadRequest, "id in path and body do not match", nil)
		return
	}

	err = a.keywordAccessObject.Update(ctx, keyword)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to update keyword")
		lib.WriteError(w, http.StatusInternalServerError, "failed to update keyword", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, keyword)
	if err != nil {
		a.log.Error(err).Log("failed to write keyword")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write keyword", err)
		return
	}
}

func (a *keywordRouter) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	uid, err := uuid.Parse(id)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to parse id")
		lib.WriteError(w, http.StatusBadRequest, "failed to parse id", err)
		return
	}

	keyword := &lib.Keyword{ID: uid}
	err = a.keywordAccessObject.Delete(ctx, keyword)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to delete keyword")
		lib.WriteError(w, http.StatusInternalServerError, "failed to delete keyword", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *keywordRouter) post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	keyword := &lib.Keyword{}
	err := lib.ReadJSON(r, keyword)
	if err != nil {
		a.log.Error(err).Log("failed to read keyword")
		lib.WriteError(w, http.StatusBadRequest, "failed to read keyword", err)
		return
	}

	err = a.keywordAccessObject.Create(ctx, keyword)
	if err != nil {
		a.log.Error(err).Log("failed to create keyword")
		lib.WriteError(w, http.StatusInternalServerError, "failed to create keyword", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusCreated, keyword)
	if err != nil {
		a.log.Error(err).Log("failed to write keyword")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write keyword", err)
		return
	}
}
