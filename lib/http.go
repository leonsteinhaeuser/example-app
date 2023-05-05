package lib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

var (
	ErrContentTypeMismatch = errors.New("http content type mismatch")
	ErrStatusCodeMismatch  = errors.New("http status code mismatch")
)

// CheckAndParseResponse checks the given response for the expected status code and content type and parses the body
func CheckAndParseResponse[T any](rsp *http.Response, expectedStatus int) (*T, error) {
	data := new(T)
	// check for content type
	if ct := rsp.Header.Get("Content-Type"); ct != "application/json" {
		return nil, fmt.Errorf("%w; expected %q got %q", ErrContentTypeMismatch, "application/json", ct)
	}
	// check status code
	if rsp.StatusCode != expectedStatus {
		// if status code > 400 the server returned an error
		if rsp.StatusCode >= 400 && rsp.StatusCode < 600 {
			httpErr := &HttpError{}
			err := json.NewDecoder(rsp.Body).Decode(httpErr)
			if err != nil {
				return nil, err
			}
			return nil, httpErr
		}
		return nil, fmt.Errorf("%w: got %d want %d", ErrStatusCodeMismatch, rsp.StatusCode, expectedStatus)
	}
	// decode body
	err := json.NewDecoder(rsp.Body).Decode(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Client[T any] interface {
	GetByID(ctx context.Context, id string) (*T, error)
	List(ctx context.Context) ([]T, error)
	Create(ctx context.Context, data T) (*T, error)
	Update(ctx context.Context, data T) (*T, error)
	Delete(ctx context.Context, data T) error
}

type CustomArticleClient[T any] interface {
	Client[T]
	DeleteByArticleID(ctx context.Context, data T) error
}

type HttpError struct {
	Status  int
	Message string
	Reason  string
}

// Error implements the error interface
func (h *HttpError) Error() string {
	return fmt.Sprintf("%d %s %q", h.Status, h.Message, h.Reason)
}

// Healthz is a simple health check handler
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Router is an interface that can be implemented by a struct to register
// a set of routes with a chi router.
type Router interface {
	Router(chi.Router)
}

// GetIntParam returns the value of the given param as int
func GetStringParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}

// WriteJSON parses the given data to json and writes it to the response writer
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	JSONHeaderStatus(w, statusCode)
	return json.NewEncoder(w).Encode(data)
}

// JSONHeaderStatus sets the content type header to application/json and the status code
func JSONHeaderStatus(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}

// WriteError parses the given error to json and writes it to the response writer
func WriteError(w http.ResponseWriter, statusCode int, msg string, err error) error {
	return WriteJSON(w, statusCode, HttpError{
		Status:  statusCode,
		Message: msg,
		Reason:  err.Error(),
	})
}

// ReadJSON parses the request body to the given data struct
func ReadJSON(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}

// WalkRoutes logs all registered routes
func WalkRoutes(r chi.Router, log log.Logger) {
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Debug().
			Field("method", method).
			Field("route", route).
			Log("registered route")
		return nil
	})
}
