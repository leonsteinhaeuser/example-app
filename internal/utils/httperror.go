package utils

import "net/http"

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func NewHTTPError(code int, message string, err error) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
		Error:   err.Error(),
	}
}

// Write writes the HTTPError to the ResponseWriter.
func (e *HTTPError) Write(w http.ResponseWriter) error {
	return WriteJSON(w, e.Code, e)
}
