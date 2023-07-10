package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/leonsteinhaeuser/example-app/internal/log"
)

const (
	RequestIDKey contextKey = 0

	HeaderRequestID = "X-Request-ID"
)

type contextKey int

// LoggerMiddleware is a middleware that logs incoming requests.
func LoggerMiddleware(log log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info().Field("method", r.Method).Field("path", r.URL.Path).Field("request-id", RequestIDFromContext(r.Context())).Log("received incoming request")
			next.ServeHTTP(w, r)
		})
	}
}

// RequestID is a middleware that adds a request ID to the context.
func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req := r.Context().Value(RequestIDKey)
			reqID := uuid.NewString()
			if req != nil {
				reqID = req.(string)
			}
			ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequestIDFromContext returns the request ID from the context.
func RequestIDFromContext(ctx context.Context) string {
	return ctx.Value(RequestIDKey).(string)
}
