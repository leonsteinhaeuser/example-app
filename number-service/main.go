package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
	"github.com/leonsteinhaeuser/example-app/number-service/api"
)

var (
	logr log.Logger           = log.NewZerlog()
	pl   lib.ProcessLifecycle = lib.NewProcessLifecycle([]os.Signal{os.Interrupt, os.Kill})
)

func main() {
	ctx, cf := context.WithTimeout(context.Background(), 30*time.Second)
	defer cf()

	mux := chi.NewRouter()
	mux.Use(log.LoggerMiddleware(logr))
	mux.Get("/healthz", lib.Healthz)
	api.NewNumberRouter(logr).Router(mux)
	lib.WalkRoutes(mux, logr)

	logr.Info().Log("starting number-service with address: 0.0.0.0:1111")
	go func() {
		err := http.ListenAndServe(":1111", mux)
		if err != nil {
			logr.Panic(err).Log("something went wrong with the server")
		}
	}()

	pl.Wait()
	pl.Shutdown(ctx)
}
