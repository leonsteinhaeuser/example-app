package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
	"github.com/leonsteinhaeuser/example-app/view-service/views"
)

var (
	clog log.Logger = log.NewZerlog()

	pl lib.ProcessLifecycle = lib.NewProcessLifecycle([]os.Signal{os.Interrupt, os.Kill})
)

func main() {
	ctx, cf := context.WithTimeout(context.Background(), 30*time.Second)
	defer cf()

	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.NoCache)
	mux.Use(middleware.CleanPath)
	mux.Use(middleware.Logger)
	mux.Use(middleware.AllowContentType("application/json"))
	mux.Use(middleware.Recoverer)
	mux.Get("/healthz", lib.Healthz)
	views.ArticleRouter(mux)
	views.NumberRouter(mux)
	views.HomeRouter(mux)

	lib.WalkRoutes(mux, clog)

	go func() {
		err := http.ListenAndServe(":2222", mux)
		if err != nil {
			clog.Panic(err).Log("something went wrong with the server")
		}
	}()

	pl.Wait()
	pl.Shutdown(ctx)
}
