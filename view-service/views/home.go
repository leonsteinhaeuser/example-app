package views

import (
	"html/template"
	"net/http"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"

	_ "embed"
)

var (
	//go:embed home.html
	templateHome string
	//go:embed main.css
	mainCSS []byte
)

type homeRouter struct {
	log log.Logger

	endpoints    *[]string
	homeTemplate *template.Template
}

func NewHomeRouter(log log.Logger, endpoints *[]string) lib.Router {
	log.Info().Field("endpoints", *endpoints).Log("creating home router")

	// template for list articles
	homeTemplate, err := template.New("home.html").Funcs(sprig.FuncMap()).Parse(templateHome)
	if err != nil {
		log.Panic(err).Log("failed to parse home.html template")
	}

	return &homeRouter{
		log:          log,
		endpoints:    endpoints,
		homeTemplate: homeTemplate,
	}
}

func (h *homeRouter) Router(r chi.Router) {
	r.Get("/", h.homeEndpoint)
	r.Get("/main.css", h.cssEndpoint)
}

func (h *homeRouter) homeEndpoint(w http.ResponseWriter, r *http.Request) {
	endpointList := make([]string, len(*h.endpoints))
	copy(endpointList, *h.endpoints)

	for idx, endpoint := range *h.endpoints {
		endpointList[idx] = lib.URLSchemeOrDefault(r.URL, "http") + "://" + r.Host + endpoint
	}

	err := h.homeTemplate.Execute(w, endpointList)
	if err != nil {
		h.log.Error(err).Log("failed to render template")
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *homeRouter) cssEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(mainCSS)
}
