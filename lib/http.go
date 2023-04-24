package lib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrContentTypeMismatch = errors.New("http content type mismatch")
	ErrStatusCodeMismatch  = errors.New("http status code mismatch")
)

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

type HttpError struct {
	Status  int
	Message string
	Err     string
	Value   any
}

func (h *HttpError) Error() string {
	return fmt.Sprintf("%d %s %q %v", h.Status, h.Message, h.Err, h.Value)
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
