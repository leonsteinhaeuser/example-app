package api

import (
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

type numberRouter struct {
	log log.Logger
}

func NewNumberRouter(log log.Logger) lib.Router {
	return &numberRouter{
		log: log,
	}
}

func (n *numberRouter) Router(r chi.Router) {
	r.Get("/", n.numberEndpoint)
}

func (n *numberRouter) numberEndpoint(w http.ResponseWriter, r *http.Request) {
	err := lib.WriteJSON(w, http.StatusOK, lib.NumberResponse{
		Number: rand.Int63(),
	})
	if err != nil {
		n.log.Error(err).Log("failed to write json")
		lib.WriteError(w, http.StatusInternalServerError, "failed to write json", err)
		return
	}
}
