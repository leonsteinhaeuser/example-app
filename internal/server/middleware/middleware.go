package middleware

import (
	"context"
	"net/http"
	"os"

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
			// get request ID from header if present
			// otherwise generate a new one
			reqID := RequestIDFromHeader(r)
			if reqID == "" {
				host, err := os.Hostname()
				if err != nil {
					reqID = uuid.NewString()
				} else {
					reqID = host + "/" + uuid.NewString()
				}
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

// RequestIDFromHeader returns the request ID from the header.
func RequestIDFromHeader(r *http.Request) string {
	return r.Header.Get(HeaderRequestID)
}
