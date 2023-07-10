package v1

import (
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	"github.com/leonsteinhaeuser/example-app/internal/server"
	"github.com/leonsteinhaeuser/example-app/internal/utils"
)

var (
	_ server.Router = (*NumberRouter)(nil)
)

type NumberResponse struct {
	Number int64 `json:"number"`
}

// NumberRouter is a router for the number service.
type NumberRouter struct {
	log log.Logger
}

func NewNumberRouter(log log.Logger) *NumberRouter {
	return &NumberRouter{
		log: log,
	}
}

func (t *NumberRouter) Router(rt chi.Router) {
	rt.Get("/", t.getNumber)
}

func (t *NumberRouter) getNumber(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, NumberResponse{
		Number: rand.Int63(),
	})
}
