package views

import (
	_ "embed"
	"html/template"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/article-service/client"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/rs/zerolog/log"
)

var (
	//go:embed list_article.html
	templateListArticles string
	//go:embed get_article.html
	templateGetArticle string

	// rendered templates
	listArticlesTemplate *template.Template
	getArticlesTemplate  *template.Template

	articleServiceClient lib.Client[lib.Article] = client.NewArticleClient(os.Getenv("ARTICLE_SERVICE_URL"))
)

func init() {
	// template for list articles
	t, err := template.New("list_article.html").Parse(templateListArticles)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse list_article.gohtml template")
	}
	listArticlesTemplate = t
	// template for get article
	t, err = template.New("get_article.html").Parse(templateGetArticle)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse list_article.gohtml template")
	}
	getArticlesTemplate = t
}

func ArticleRouter(rt chi.Router) {
	rt.Route("/article", func(r chi.Router) {
		r.Get("/list", listArticlesEndpoint)
		r.Get("/{id}", getArticleByIDEndpoint)
	})
}

func listArticlesEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	articles, err := articleServiceClient.List(ctx)
	if err != nil {
		return
	}

	log.Info().Any("articles", articles).Msg("articles")
	err = listArticlesTemplate.Execute(w, articles)
	if err != nil {
		log.Error().Err(err).Msg("failed to render template")
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func getArticleByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	articles, err := articleServiceClient.List(ctx)
	if err != nil {
		return
	}

	err = getArticlesTemplate.Execute(w, articles)
	if err != nil {
		log.Error().Err(err).Msg("failed to render template")
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}
