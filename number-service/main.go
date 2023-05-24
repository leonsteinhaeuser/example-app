package main

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/leonsteinhaeuser/example-app/internal"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(internal.NumberResponse{
			Number: rand.Int63(),
		})
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	http.ListenAndServe(":1111", nil)
}
