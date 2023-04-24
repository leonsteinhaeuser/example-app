package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/leonsteinhaeuser/example-app/article-service/accessobjects"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

type articleRouter struct {
	articleAccessObject accessobjects.Article
	log                 log.Logger
}

func NewArticleRouter(articleAccessObject accessobjects.Article, log log.Logger) *articleRouter {
	return &articleRouter{
		articleAccessObject: articleAccessObject,
		log:                 log,
	}
}

func (a *articleRouter) Router(r chi.Router) {
	r.Route("/article", func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", a.get)
			r.Put("/", a.put)
			r.Delete("/", a.delete)
		})
		r.Get("/", a.list)
		r.Post("/", a.post)
	})
}

func (a *articleRouter) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	article, err := a.articleAccessObject.Get(ctx, id)
	if err != nil {
		a.log.Error(err).Log("failed to get article")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get article from database", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, article)
	if err != nil {
		a.log.Error(err).Log("failed to write article")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write article", err)
		return
	}
}

func (a *articleRouter) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	articles, err := a.articleAccessObject.List(ctx)
	if err != nil {
		a.log.Error(err).Log("failed to list articles")
		lib.WriteError(w, http.StatusInternalServerError, "failed to list articles", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, articles)
	if err != nil {
		a.log.Error(err).Log("failed to write articles")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write articles", err)
		return
	}
}

func (a *articleRouter) put(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	article, err := a.articleAccessObject.Get(ctx, id)
	if err != nil {
		a.log.Error(err).Log("failed to get article")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get article from database", err)
		return
	}

	err = lib.ReadJSON(r, article)
	if err != nil {
		a.log.Error(err).Log("failed to read article")
		lib.WriteError(w, http.StatusBadRequest, "failed to read article", err)
		return
	}

	if article.ID.String() != id {
		a.log.Warn().
			Field("urlID", id).
			Field("body", article.ID.String()).
			Log("id in path and body do not match")
		lib.WriteError(w, http.StatusBadRequest, "id in path and body do not match", nil)
		return
	}

	err = a.articleAccessObject.Update(ctx, article)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to update article")
		lib.WriteError(w, http.StatusInternalServerError, "failed to update article", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, article)
	if err != nil {
		a.log.Error(err).Log("failed to write article")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write article", err)
		return
	}
}

func (a *articleRouter) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	uid, err := uuid.Parse(id)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to parse id")
		lib.WriteError(w, http.StatusBadRequest, "failed to parse id", err)
		return
	}

	article := &lib.Article{ID: uid}
	err = a.articleAccessObject.Delete(ctx, article)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to delete article")
		lib.WriteError(w, http.StatusInternalServerError, "failed to delete article", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *articleRouter) post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	article := &lib.Article{}
	err := lib.ReadJSON(r, article)
	if err != nil {
		a.log.Error(err).Log("failed to read article")
		lib.WriteError(w, http.StatusBadRequest, "failed to read article", err)
		return
	}

	err = a.articleAccessObject.Create(ctx, article)
	if err != nil {
		a.log.Error(err).Log("failed to create article")
		lib.WriteError(w, http.StatusInternalServerError, "failed to create article", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusCreated, article)
	if err != nil {
		a.log.Error(err).Log("failed to write article")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write article", err)
		return
	}
}
