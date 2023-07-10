package main

import (
	"net/http"

	"github.com/leonsteinhaeuser/example-app/internal/env"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/server"
	v1 "github.com/leonsteinhaeuser/example-app/number-service/api/v1"
)

var (
	logr = log.NewZerlog()

	httpServer = server.NewDefaultServer(logr, env.GetStringEnvOrDefault("LISTEN_ADDRESS", ":1111"))
	httpRouter = server.NewGenericRouter()
)

func main() {
	httpServer.AddRouter(v1.NewNumberRouter(logr))
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
