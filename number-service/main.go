package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/leonsteinhaeuser/example-app/internal"
	"github.com/leonsteinhaeuser/example-app/internal/env"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/server"
	"github.com/leonsteinhaeuser/example-app/internal/utils"
	"github.com/leonsteinhaeuser/example-app/number-service/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var (
	// externalURL
	externalURL string = env.GetStringEnvOrDefault("EXTERNAL_URL", "http://localhost:1111")

	logr = log.NewZerlog()

	httpServer = server.NewDefaultServer(logr, env.GetStringEnvOrDefault("LISTEN_ADDRESS", ":1111"))
	httpRouter = server.NewGenericRouter()
)

// @title Number Service API
// @version 1.0
// @description This is the number-service.

// @contact.name API Support
// @contact.url http://github.com/leonsteinhaeuser/example-app
// @contact.email mail@example.localhost

// @license.name GNU GENERAL PUBLIC LICENSE 3.0
// @license.url https://github.com/leonsteinhaeuser/example-app/blob/main/LICENSE

// @host number-service
// @BasePath /
func main() {
	httpRouter.AddEndpoint("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(internal.NumberResponse{
			Number: rand.Int63(),
		})
	})
	httpRouter.AddEndpoint("GET", "/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	// override swagger docs host
	docs.SwaggerInfo.Host = utils.ReplaceAllToReplace(externalURL, "", "http://", "https://")
	// define swagger endpoint to serve swagger docs
	dfltrt := server.NewGenericRouter()
	dfltrt.AddEndpoint(http.MethodGet, "/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", externalURL)),
	))
	httpServer.AddRouter(httpRouter)
	err := httpServer.Start()
	if err != nil {
		panic(err)
	}
}
