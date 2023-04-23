package main

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/rs/zerolog/log"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(lib.NumberResponse{
			Number: rand.Int63(),
		})
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	log.Info().Msg("starting number-service with address: 0.0.0.0:1111")
	http.ListenAndServe(":1111", nil)
}
