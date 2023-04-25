package views

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"

	_ "embed"
)

var (
	//go:embed number_view.html
	numberViewTemplate string

	indexTemplate = template.New("index.html")
)

type numberRouter struct {
	log log.Logger

	client        lib.Client[lib.NumberResponse]
	indexTemplate *template.Template
}

func NewNumberRouter(log log.Logger, client lib.Client[lib.NumberResponse]) lib.Router {
	numberTemplate, err := indexTemplate.Parse(numberViewTemplate)
	if err != nil {
		log.Panic(err).Log("failed to parse number_view.html template")
	}

	return &numberRouter{
		indexTemplate: numberTemplate,
		client:        client,
		log:           log,
	}
}

func (n *numberRouter) Router(r chi.Router) {
	r.Get("/number", n.numberEndpoint)
}

func (n *numberRouter) numberEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	number, err := n.client.GetByID(ctx, "")
	if err != nil {
		n.log.Error(err).Log("failed to get number")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = indexTemplate.Execute(w, map[string]interface{}{
		"number": number,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
