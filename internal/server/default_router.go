package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GenericRouter struct {
	Endpoints []Endpoint
}

type Endpoint struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

func NewGenericRouter() *GenericRouter {
	return &GenericRouter{
		Endpoints: []Endpoint{},
	}
}

func (g *GenericRouter) AddEndpoint(method, path string, handler http.HandlerFunc) {
	g.Endpoints = append(g.Endpoints, Endpoint{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

func (g *GenericRouter) Router(r chi.Router) {
	for _, e := range g.Endpoints {
		r.MethodFunc(e.Method, e.Path, e.Handler.ServeHTTP)
	}
}
