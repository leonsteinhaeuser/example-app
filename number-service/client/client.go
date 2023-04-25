package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/leonsteinhaeuser/example-app/lib"
)

var (
	ErrNotSupported = errors.New("not supported")

	_ lib.Client[lib.NumberResponse] = &numberClient{}
)

type numberClient struct {
	serviceURL string
	client     http.Client
}

// NewNumberClient returns a new ArticleClient
// The client configures an http.Client with a 5 second timeout
// serviceURL should be in the format of http://host:port
func NewNumberClient(serviceURL string) lib.Client[lib.NumberResponse] {
	return &numberClient{
		serviceURL: serviceURL,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *numberClient) basePath(optionals ...string) (string, error) {
	paths := append([]string{"number"}, optionals...)
	path, err := url.JoinPath(c.serviceURL, paths...)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (c *numberClient) GetByID(ctx context.Context, id string) (*lib.NumberResponse, error) {
	path, err := c.basePath()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	number, err := lib.CheckAndParseResponse[lib.NumberResponse](rsp, 200)
	if err != nil {
		return nil, err
	}
	return number, nil
}

func (c *numberClient) List(ctx context.Context) ([]lib.NumberResponse, error) {
	return nil, fmt.Errorf("list endpoint %w", ErrNotSupported)
}

func (c *numberClient) Create(ctx context.Context, article lib.NumberResponse) (*lib.NumberResponse, error) {
	return nil, fmt.Errorf("create endpoint %w", ErrNotSupported)
}

func (c *numberClient) Update(ctx context.Context, article lib.NumberResponse) (*lib.NumberResponse, error) {
	return nil, fmt.Errorf("update endpoint %w", ErrNotSupported)
}

func (c *numberClient) Delete(ctx context.Context, article lib.NumberResponse) error {
	return fmt.Errorf("delete endpoint %w", ErrNotSupported)
}
