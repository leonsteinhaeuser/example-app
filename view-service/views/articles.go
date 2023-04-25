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
	//go:embed list_article.html
	templateListArticles string
	//go:embed get_article.html
	templateGetArticle string
)

type articleRouter struct {
	log log.Logger

	listArticlesTemplate *template.Template
	getArticlesTemplate  *template.Template

	client lib.Client[lib.Article]
}

func NewArticleRouter(log log.Logger, articleClient lib.Client[lib.Article]) lib.Router {
	// template for list articles
	listArticleTemplate, err := template.New("list_article.html").Parse(templateListArticles)
	if err != nil {
		log.Panic(err).Log("failed to parse list_article.html template")
	}
	// template for get article
	getArticlesTemplate, err := template.New("get_article.html").Parse(templateGetArticle)
	if err != nil {
		log.Panic(err).Log("failed to parse list_article.html template")
	}
	return &articleRouter{
		listArticlesTemplate: listArticleTemplate,
		getArticlesTemplate:  getArticlesTemplate,
		log:                  log,
		client:               articleClient,
	}
}

func (a *articleRouter) Router(r chi.Router) {
	r.Route("/article", func(r chi.Router) {
		r.Get("/list", a.listArticlesEndpoint)
		r.Get("/{id}", a.getArticleByIDEndpoint)
	})
}

func (a *articleRouter) listArticlesEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	articles, err := a.client.List(ctx)
	if err != nil {
		a.log.Error(err).Log("failed to get articles")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get articles", err)
		return
	}

	err = a.listArticlesTemplate.Execute(w, articles)
	if err != nil {
		a.log.Error(err).Log("failed to render template")
		lib.WriteError(w, http.StatusInternalServerError, "failed to render template", err)
		return
	}
}

func (a *articleRouter) getArticleByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := lib.GetStringParam(r, "id")

	article, err := a.client.GetByID(ctx, id)
	if err != nil {
		a.log.Error(err).Log("failed to get articles")
		lib.WriteError(w, http.StatusInternalServerError, "failed to get articles", err)
		return
	}

	err = a.getArticlesTemplate.Execute(w, article)
	if err != nil {
		a.log.Error(err).Log("failed to render template")
		lib.WriteError(w, http.StatusInternalServerError, "failed to render template", err)
		return
	}
}
