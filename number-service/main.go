package main

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/leonsteinhaeuser/example-app/lib"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(lib.NumberResponse{
			Number: rand.Int63(),
		})
	})
	http.ListenAndServe(":1111", nil)
}
