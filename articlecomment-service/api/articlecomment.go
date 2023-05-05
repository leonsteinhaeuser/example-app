package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/leonsteinhaeuser/example-app/articlecomment-service/accessobjects"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

type articleCommentRouter struct {
	articleCommentAccessObject accessobjects.ArticleComment
	log                        log.Logger
}

func NewArticleCommentRouter(articleCommentAccessObject accessobjects.ArticleComment, log log.Logger) *articleCommentRouter {
	return &articleCommentRouter{
		articleCommentAccessObject: articleCommentAccessObject,
		log:                        log,
	}
}

func (a *articleCommentRouter) Router(r chi.Router) {
	r.Route("/comment", func(r chi.Router) {
		r.Delete("/article/{id}", a.deleteByArticleID)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", a.get)
			r.Put("/", a.put)
			r.Delete("/", a.delete)
		})
		r.Get("/", a.list)
		r.Post("/", a.post)
	})
}

func (a *articleCommentRouter) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	articleComment, err := a.articleCommentAccessObject.Get(ctx, id)
	if err != nil {
		a.log.Error(err).Log("failed to get articleComment")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get articleCommentfrom database", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, articleComment)
	if err != nil {
		a.log.Error(err).Log("failed to write articleComment")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write articleComment", err)
		return
	}
}

func (a *articleCommentRouter) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	articleComments, err := a.articleCommentAccessObject.List(ctx)
	if err != nil {
		a.log.Error(err).Log("failed to list articleComments")
		lib.WriteError(w, http.StatusInternalServerError, "failed to list articleComments", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, articleComments)
	if err != nil {
		a.log.Error(err).Log("failed to write articleComments")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write articleComments", err)
		return
	}
}

func (a *articleCommentRouter) put(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	articleComment, err := a.articleCommentAccessObject.Get(ctx, id)
	if err != nil {
		a.log.Error(err).Log("failed to get articleComment")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get articleCommentfrom database", err)
		return
	}

	err = lib.ReadJSON(r, articleComment)
	if err != nil {
		a.log.Error(err).Log("failed to read articleComment")
		lib.WriteError(w, http.StatusBadRequest, "failed to read articleComment", err)
		return
	}

	if articleComment.ID.String() != id {
		a.log.Warn().
			Field("urlID", id).
			Field("body", articleComment.ID.String()).
			Log("id in path and body do not match")
		lib.WriteError(w, http.StatusBadRequest, "id in path and body do not match", nil)
		return
	}

	err = a.articleCommentAccessObject.Update(ctx, articleComment)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to update articleComment")
		lib.WriteError(w, http.StatusInternalServerError, "failed to update articleComment", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusOK, articleComment)
	if err != nil {
		a.log.Error(err).Log("failed to write articleComment")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write articleComment", err)
		return
	}
}

func (a *articleCommentRouter) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	uid, err := uuid.Parse(id)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to parse id")
		lib.WriteError(w, http.StatusBadRequest, "failed to parse id", err)
		return
	}

	articleComment := &lib.ArticleComment{ID: uid}
	err = a.articleCommentAccessObject.Delete(ctx, articleComment)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to delete articleComment")
		lib.WriteError(w, http.StatusInternalServerError, "failed to delete articleComment", err)
		return
	}

	lib.JSONHeaderStatus(w, http.StatusNoContent)
}

func (a *articleCommentRouter) deleteByArticleID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")
	a.log.Debug().Field("id", id).Log("deleting article comment by id")

	uid, err := uuid.Parse(id)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to parse id")
		lib.WriteError(w, http.StatusBadRequest, "failed to parse id", err)
		return
	}

	articleComment := &lib.ArticleComment{ArticleID: uid}
	err = a.articleCommentAccessObject.DeleteByArticleID(ctx, articleComment)
	if err != nil {
		a.log.Error(err).Field("id", id).Log("failed to delete articleComment")
		lib.WriteError(w, http.StatusInternalServerError, "failed to delete articleComment", err)
		return
	}

	lib.JSONHeaderStatus(w, http.StatusNoContent)
}

func (a *articleCommentRouter) post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	articleComment := &lib.ArticleComment{}
	err := lib.ReadJSON(r, articleComment)
	if err != nil {
		a.log.Error(err).Log("failed to read articleComment")
		lib.WriteError(w, http.StatusBadRequest, "failed to read articleComment", err)
		return
	}

	err = a.articleCommentAccessObject.Create(ctx, articleComment)
	if err != nil {
		a.log.Error(err).Log("failed to create articleComment")
		lib.WriteError(w, http.StatusInternalServerError, "failed to create articleComment", err)
		return
	}

	err = lib.WriteJSON(w, http.StatusCreated, articleComment)
	if err != nil {
		a.log.Error(err).Log("failed to write articleComment")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write articleComment", err)
		return
	}
}
