package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/leonsteinhaeuser/example-app/public"
	"github.com/leonsteinhaeuser/example-app/thread-service/thread"
)

var (
	_ public.Client[thread.Thread] = (*DefaultClient)(nil)
)

type DefaultClient struct {
	serviceURI string
	httpClient *http.Client
}

// NewDefaultClient creates a new default client.
// The service address is the address of the thread service.
// The service address must be provided with the protocol and port (e.g. http://localhost:8080).
func NewDefaultClient(serviceAddress string) (*DefaultClient, error) {
	serviceURI, err := url.JoinPath(serviceAddress, "/thread")
	if err != nil {
		return nil, fmt.Errorf("failed to join url path: %w", err)
	}
	// create a copy of the default http client
	clnt := *http.DefaultClient
	return &DefaultClient{
		serviceURI: serviceURI,
		httpClient: &clnt,
	}, nil
}

// WithHTTPClient allows to set a custom http client.
func (d *DefaultClient) WithHTTPClient(clnt *http.Client) {
	d.httpClient = clnt
}

// doRequest prepares a request with the given method, id and body and executes it.
func (d *DefaultClient) doRequest(ctx context.Context, method string, id string, body any) (*http.Response, error) {
	// create the body buffer
	bodyBuf := bytes.NewBuffer([]byte{})
	if body != nil {
		bts, err := json.Marshal(body)
		if err != nil {
			return nil, &ClientMarshalError{
				Message: "failed to marshal body",
				Err:     err,
			}
		}
		*bodyBuf = *bytes.NewBuffer(bts)
	}
	// join the service address and the path
	urlPath, err := url.JoinPath(d.serviceURI, id)
	if err != nil {
		return nil, &ClientRequestError{
			Message: "failed to join url path",
			Err:     err,
		}
	}
	// create the request
	req, err := http.NewRequest(method, urlPath, bodyBuf)
	if err != nil {
		return nil, &ClientRequestError{
			Message: "failed to create new http request",
			Err:     err,
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	// fire the request
	rsp, err := d.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, &ClientRequestError{
			Message: "failed to execute request",
			Err:     err,
		}
	}
	return rsp, nil
}

func (d *DefaultClient) Get(ctx context.Context, id string) (*thread.Thread, error) {
	rsp, err := d.doRequest(ctx, http.MethodGet, id, nil)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return nil, &ClientRequestError{
			Message: "unexpected status code: " + rsp.Status,
			Err:     err,
		}
	}
	t := &thread.Thread{}
	if err := json.NewDecoder(rsp.Body).Decode(t); err != nil {
		return nil, &ClientMarshalError{
			Message: "failed to unmarshal response body",
			Err:     err,
		}
	}
	return t, nil
}

func (d *DefaultClient) List(ctx context.Context) ([]*thread.Thread, error) {
	rsp, err := d.doRequest(ctx, http.MethodGet, "", nil)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return nil, &ClientRequestError{
			Message: "unexpected status code: " + rsp.Status,
			Err:     err,
		}
	}
	t := &[]*thread.Thread{}
	if err := json.NewDecoder(rsp.Body).Decode(t); err != nil {
		return nil, &ClientMarshalError{
			Message: "failed to unmarshal response body",
			Err:     err,
		}
	}
	return *t, nil
}

func (d *DefaultClient) Create(ctx context.Context, t *thread.Thread) error {
	rsp, err := d.doRequest(ctx, http.MethodPost, "", t)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusCreated {
		return &ClientRequestError{
			Message: "unexpected status code: " + rsp.Status,
			Err:     err,
		}
	}
	return nil
}

func (d *DefaultClient) Update(ctx context.Context, t *thread.Thread) error {
	rsp, err := d.doRequest(ctx, http.MethodPut, t.ID.String(), t)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return &ClientRequestError{
			Message: "unexpected status code: " + rsp.Status,
			Err:     err,
		}
	}
	return nil
}

func (d *DefaultClient) Delete(ctx context.Context, t *thread.Thread) error {
	rsp, err := d.doRequest(ctx, http.MethodDelete, t.ID.String(), t)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusNoContent {
		return &ClientRequestError{
			Message: "unexpected status code: " + rsp.Status,
			Err:     err,
		}
	}
	return nil
}
