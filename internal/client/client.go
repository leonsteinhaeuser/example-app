package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/leonsteinhaeuser/example-app/internal/server/middleware"
)

type HttpClient struct {
	// Address is the base url of the service to talk to
	Address string
	Client  *http.Client
}

func (r *HttpClient) Request(ctx context.Context, method, path string, data any) (*http.Response, error) {
	dataBTS, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	// join url paths properly
	rslt, err := url.JoinPath(r.Address, path)
	if err != nil {
		return nil, err
	}
	// prepare request
	req, err := http.NewRequest(method, rslt, bytes.NewBuffer(dataBTS))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set(middleware.HeaderRequestID, middleware.RequestIDFromContext(ctx))
	req.Header.Set("Content-Type", "application/json")
	// send request
	rsp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}
