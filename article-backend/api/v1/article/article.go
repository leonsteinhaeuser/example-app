package article

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/leonsteinhaeuser/example-app/internal/db"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/server"
	"github.com/leonsteinhaeuser/example-app/internal/utils"
)

type articleRouter struct {
	log log.Logger

	db db.Repository
}

func NewArticleRouter(log log.Logger, db db.Repository) *articleRouter {
	return &articleRouter{
		log: log,
		db:  db,
	}
}

func (t *articleRouter) Router(rt chi.Router) {
	rt.Route("/articles", func(rt chi.Router) {
		rt.Route("/{id}", func(rt chi.Router) {
			rt.Get("/", t.getArticle)
			rt.Put("/", t.updateArticle)
			rt.Delete("/", t.deleteArticle)
		})
		rt.Get("/", t.getArticles)
		rt.Post("/", t.createArticle)
	})
}

func (t *articleRouter) createArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	article := &Article{}
	err := utils.ReadJSON(r, article)
	if err != nil {
		t.log.Error(err).Log("failed to parse JSON body")
		utils.WriteJSON(w, http.StatusBadRequest, server.Error{
			Status:  http.StatusBadRequest,
			Message: "failed to parse JSON body",
			Error:   err.Error(),
		})
		return
	}

	err = t.db.Create(ctx, article)
	if err != nil {
		t.log.Error(err).Log("failed to create article")
		utils.WriteJSON(w, http.StatusInternalServerError, server.Error{
			Status:  http.StatusInternalServerError,
			Message: "failed to create article",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, article)
}

// getArticles returns a list of articles.
// Optional query parameters:
// - published: bool
// - author_id: uuid
// - limit: int
// - published_before: timestamp
// - published_after: timestamp
func (t *articleRouter) getArticles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	articles := &[]*Article{}
	dbtx := t.db.Find(articles)

	t.log.Debug().Field("query", r.URL.Query()).Log("query")

	// apply optional filters to transaction

	// filter by "published"
	if publshd, err := strconv.ParseBool(r.URL.Query().Get("published")); err == nil {
		dbtx = dbtx.Where("published = ?", publshd)
	}
	// filter by "author"
	if author := r.URL.Query().Get("author_id"); author != "" {
		dbtx = dbtx.Where("author_id = ?", author)
	}
	// filter by "limit"
	if limit, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		dbtx = dbtx.Limit(limit)
	}
	// filter by published_before
	if publishedBefore := r.URL.Query().Get("published_before"); publishedBefore != "" {
		dbtx = dbtx.Where("published_at < ?", publishedBefore)
	}
	// filter by published_after
	if publishedAfter := r.URL.Query().Get("published_after"); publishedAfter != "" {
		dbtx = dbtx.Where("published_at > ?", publishedAfter)
	}

	err := dbtx.Commit(ctx)
	if err != nil {
		t.log.Error(err).Log("failed to list articles")
		utils.WriteJSON(w, http.StatusInternalServerError, server.Error{
			Status:  http.StatusInternalServerError,
			Message: "failed to list articles",
			Error:   err.Error(),
		})
		return
	}

	t.log.Debug().Field("articles", articles).Log("articles")

	utils.WriteJSON(w, http.StatusOK, articles)
}

func (t *articleRouter) getArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	article := &Article{}
	err := t.db.Find(article).Where("id = ?", id).Commit(ctx)
	if err != nil {
		t.log.Error(err).Log("failed to get article")
		utils.WriteJSON(w, http.StatusInternalServerError, server.Error{
			Status:  http.StatusInternalServerError,
			Message: "failed to get article",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, article)
}

func (t *articleRouter) updateArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	article := &Article{}
	err := utils.ReadJSON(r, article)
	if err != nil {
		t.log.Error(err).Log("failed to parse JSON body")
		utils.WriteJSON(w, http.StatusBadRequest, server.Error{
			Status:  http.StatusBadRequest,
			Message: "failed to parse JSON body",
			Error:   err.Error(),
		})
		return
	}

	err = t.db.Update(article).Where("id = ?", id).Commit(ctx)
	if err != nil {
		t.log.Error(err).Log("failed to update article")
		utils.WriteJSON(w, http.StatusInternalServerError, server.Error{
			Status:  http.StatusInternalServerError,
			Message: "failed to update article",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}

func (t *articleRouter) deleteArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	t.log.Debug().Field("id", id).Log("delete article")

	err := t.db.Delete(&Article{ID: uuid.MustParse(id)}).Where("id = ?", id).Commit(ctx)
	if err != nil {
		t.log.Error(err).Log("failed to delete article")
		utils.WriteJSON(w, http.StatusInternalServerError, server.Error{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete article",
			Error:   err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, map[string]any{})
}
