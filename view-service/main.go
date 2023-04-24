package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/view-service/views"
	"github.com/rs/zerolog/log"
)

func main() {
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

	chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		views.EndpointList = append(views.EndpointList, route)
		log.Debug().Str("method", method).Str("route", route).Msg("registered route")
		return nil
	})

	log.Info().Msg("starting view-service with address: 0.0.0.0:2222")
	http.ListenAndServe(":2222", mux)
}
