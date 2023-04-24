package main

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/rs/zerolog/log"
)

func main() {
	mux := chi.NewRouter()
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(lib.NumberResponse{
			Number: rand.Int63(),
		})
	})
	mux.Get("/healthz", lib.Healthz)

	chi.Walk(mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Debug().Str("method", method).Str("route", route).Msg("registered route")
		return nil
	})

	log.Info().Msg("starting number-service with address: 0.0.0.0:1111")
	http.ListenAndServe(":1111", mux)
}
