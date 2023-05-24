package main

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/leonsteinhaeuser/example-app/internal"
	"github.com/leonsteinhaeuser/example-app/internal/env"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/server"
)

var (
	logr = log.NewZerlog()

	httpServer = server.NewDefaultServer(logr, env.GetStringEnvOrDefault("LISTEN_ADDRESS", ":1111"))
	httpRouter = server.NewGenericRouter()
)

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
	httpServer.AddRouter(httpRouter)
	err := httpServer.Start()
	if err != nil {
		panic(err)
	}
}
