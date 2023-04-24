package views

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/lib"
)

var (
	numberServiceAddress = os.Getenv("NUMBER_SERVICE_ADDRESS")

	stringTemplate = `
	<html>
		<head>
			<title>View Service</title>
		</head>
		<body>
			<h1>View Service</h1>
			<p>Number: {{ .number }}</p>
		</body>
	</html>
		`
	indexTemplate = template.New("index.html")
)

func init() {
	idxTpl, err := indexTemplate.Parse(stringTemplate)
	if err != nil {
		panic("failed to parse index.html template")
	}
	indexTemplate = idxTpl
}

func NumberRouter(rt chi.Router) {
	rt.Get("/number", numberEndpoint)
}

func getNumberFromNumberService(ctx context.Context) (int64, error) {
	req, err := http.NewRequest("GET", numberServiceAddress, nil)
	if err != nil {
		return 0, err
	}
	// add context to request
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	// parse response
	data := lib.NumberResponse{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, err
	}
	return data.Number, nil
}

func numberEndpoint(w http.ResponseWriter, r *http.Request) {
	number, err := getNumberFromNumberService(r.Context())
	if err != nil {
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
