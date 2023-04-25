package log

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type Logger interface {
	SetModule(module string) Logger
	Trace() Field
	Debug() Field
	Info() Field
	Warn() Field
	Error(error) Field
	Panic(error) Field
}

type Field interface {
	Field(key string, value interface{}) Field
	Log(message string)
	Logf(format string, args ...interface{})
}

func LoggerMiddleware(log Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug().
				Field("method", r.Method).
				Field("path", r.URL.Path).
				Field("request_id", middleware.GetReqID(r.Context())).
				Log("incoming request")
			next.ServeHTTP(w, r)
		})
	}
}
