package main

import (
	"fmt"
	"net/http"

	"github.com/leonsteinhaeuser/example-app/internal/env"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/server"
	"github.com/leonsteinhaeuser/example-app/internal/utils"
	"github.com/leonsteinhaeuser/example-app/thread-service/client"
	"github.com/leonsteinhaeuser/example-app/thread-service/docs"
	"github.com/leonsteinhaeuser/example-app/thread-view-service/api"
	_ "github.com/leonsteinhaeuser/example-app/thread-view-service/docs"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var (
	threadServiceAddr = env.GetStringEnvOrDefault("THREAD_SERVICE_ADDRESS", "http://localhost:8080")
	listenAddr        = env.GetStringEnvOrDefault("LISTEN_ADDR", ":8080")
	// externalURL
	externalURL string = env.GetStringEnvOrDefault("EXTERNAL_URL", "http://localhost:1112")

	// logger settings
	logLevel            = env.GetIntEnvOrDefault("LOG_LEVEL", 0)
	logr     log.Logger = log.NewZerlog(func(l *zerolog.Logger) {
		*l = l.With().Logger().Level(zerolog.Level(logLevel))
	})
)

// @title Thread View Service
// @version 1.0
// @description This is the thread-view-service.

// @contact.name API Support
// @contact.url http://github.com/leonsteinhaeuser/example-app
// @contact.email mail@example.localhost

// @license.name GNU GENERAL PUBLIC LICENSE 3.0
// @license.url https://github.com/leonsteinhaeuser/example-app/blob/main/LICENSE

// @host thread-view-service
// @BasePath /
func main() {
	srvr := server.NewDefaultServer(logr, listenAddr, "application/x-www-form-urlencoded")

	client, err := client.NewDefaultClient(threadServiceAddr)
	if err != nil {
		logr.Panic(err).Log("error creating thread service client")
	}
	// override swagger docs host
	docs.SwaggerInfo.Host = utils.ReplaceAllToReplace(externalURL, "", "http://", "https://")
	// define swagger endpoint to serve swagger docs
	dfltrt := server.NewGenericRouter()
	dfltrt.AddEndpoint(http.MethodGet, "/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", externalURL)),
	))

	srvr.AddRouter(api.NewViewRouter(logr, client))

	defer srvr.Stop()
	err = srvr.Start()
	if err != nil {
		logr.Panic(err).Log("error starting server")
	}
}
