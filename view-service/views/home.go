package views

import (
	"html/template"
	"net/http"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/rs/zerolog/log"

	_ "embed"
)

var (
	//go:embed home.html
	templateHome string

	// rendered templates
	homeTemplate *template.Template

	EndpointList []string
)

func init() {
	// template for list articles
	t, err := template.New("home.html").Funcs(sprig.FuncMap()).Parse(templateHome)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse home.html template")
	}
	homeTemplate = t
}

func init() {
	idxTpl, err := indexTemplate.Parse(stringTemplate)
	if err != nil {
		panic("failed to parse index.html template")
	}
	indexTemplate = idxTpl
}

func HomeRouter(rt chi.Router) {
	rt.Get("/", homeEndpoint)
}

func homeEndpoint(w http.ResponseWriter, r *http.Request) {
	endpointList := make([]string, len(EndpointList))
	copy(endpointList, EndpointList)

	for idx, endpoint := range EndpointList {
		endpointList[idx] = lib.URLSchemeOrDefault(r.URL, "http") + "://" + r.Host + endpoint
	}

	err := homeTemplate.Execute(w, endpointList)
	if err != nil {
		log.Error().Err(err).Msg("failed to render template")
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}
