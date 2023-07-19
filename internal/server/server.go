package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	customMiddleware "github.com/leonsteinhaeuser/example-app/internal/server/middleware"
)

// Router is an interface used to define the routes of the server
type Router interface {
	Router(chi.Router)
}

type Server struct {
	logger log.Logger

	server http.Server
	router chi.Router
}

func NewDefaultServer(logger log.Logger, listen string) *Server {
	rt := chi.NewRouter()
	rt.Use(customMiddleware.OpenTelemetryMiddleware)
	rt.Use(customMiddleware.RequestID())
	rt.Use(middleware.RealIP)
	rt.Use(middleware.NoCache)
	rt.Use(middleware.CleanPath)
	rt.Use(customMiddleware.LoggerMiddleware(logger))
	rt.Use(middleware.AllowContentType("application/json"))
	rt.Use(middleware.Recoverer)
	return &Server{
		logger: logger,
		server: http.Server{
			Addr:         listen,
			Handler:      rt,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		router: rt,
	}
}

// AddMiddleware adds a middleware to the server
// The middleware is applied to all routes
// Middlewares must be added before adding routes
func (s *Server) AddMiddleware(m func(http.Handler) http.Handler) {
	s.router.Use(m)
}

// AddRouter adds a router to the server
func (s *Server) AddRouter(r Router) {
	r.Router(s.router)
}

// Start starts the server
func (s *Server) Start() error {
	chi.Walk(s.router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		s.logger.Trace().Field("method", method).Field("route", route).Log("Adding route")
		return nil
	})
	s.logger.Info().Field("listen", s.server.Addr).Log("Starting server")
	return s.server.ListenAndServe()
}

// Stop stops the server
func (s *Server) Stop() error {
	return s.server.Close()
}
