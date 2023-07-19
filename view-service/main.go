package main

import (
	"html/template"
	"net/http"

	"github.com/leonsteinhaeuser/example-app/internal/env"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/server"
	"github.com/leonsteinhaeuser/example-app/number-service/client"
)

var (
	numberServiceAddress = env.GetStringEnvOrDefault("NUMBER_SERVICE_ADDRESS", "http://localhost:1111")

	logr       = log.NewZerlog()
	httpServer = server.NewDefaultServer(logr, env.GetStringEnvOrDefault("LISTEN_ADDRESS", ":1112"))
	httpRouter = server.NewGenericRouter()

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

	numberClient *client.RestClient
)

func init() {
	idxTpl, err := indexTemplate.Parse(stringTemplate)
	if err != nil {
		panic("failed to parse index.html template")
	}
	indexTemplate = idxTpl

	numberClient = client.NewRestClient(logr, 5, numberServiceAddress)
}

func main() {
	httpRouter.AddEndpoint("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		number, err := numberClient.GetNumber(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = indexTemplate.Execute(w, map[string]interface{}{
			"number": number.Number,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	httpRouter.AddEndpoint("GET", "/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	httpServer.AddRouter(httpRouter)
	err := httpServer.Start()
	if err != nil {
		panic(err)
	}
}
