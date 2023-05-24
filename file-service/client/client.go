package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

var (
	_ lib.Client[lib.File] = &fileClient{}
)

type fileClient struct {
	log        log.Logger
	serviceURL string
	client     http.Client
}

// NewFileClient returns a new FileClient
// The client configures an http.Client with a 5 second timeout
// serviceURL should be in the format of http://host:port
func NewFileClient(log log.Logger, serviceURL string) lib.Client[lib.File] {
	log.Info().Logf("creating file client for %s", serviceURL)
	return &fileClient{
		log:        log,
		serviceURL: serviceURL,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *fileClient) basePath(optionals ...string) (string, error) {
	paths := append([]string{"file"}, optionals...)
	path, err := url.JoinPath(c.serviceURL, paths...)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (c *fileClient) GetByID(ctx context.Context, id string) (*lib.File, error) {
	path, err := c.basePath(id)
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

	file, err := lib.CheckAndParseResponse[lib.File](rsp, 200)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (c *fileClient) List(ctx context.Context) ([]lib.File, error) {
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

	users, err := lib.CheckAndParseResponse[[]lib.File](rsp, 200)
	if err != nil {
		return nil, err
	}
	return *users, nil
}

func (c *fileClient) Create(ctx context.Context, file lib.File) (*lib.File, error) {
	path, err := c.basePath()
	if err != nil {
		return nil, err
	}

	bts, err := json.Marshal(file)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(bts))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	artc, err := lib.CheckAndParseResponse[lib.File](rsp, 200)
	if err != nil {
		return nil, err
	}
	return artc, nil
}

func (c *fileClient) Update(ctx context.Context, file lib.File) (*lib.File, error) {
	path, err := c.basePath()
	if err != nil {
		return nil, err
	}

	bts, err := json.Marshal(file)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(bts))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	artc, err := lib.CheckAndParseResponse[lib.File](rsp, 200)
	if err != nil {
		return nil, err
	}
	return artc, nil
}

func (c *fileClient) Delete(ctx context.Context, file lib.File) error {
	path, err := c.basePath(file.ID.String())
	if err != nil {
		return err
	}

	bts, err := json.Marshal(file)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(bts))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	_, err = lib.CheckAndParseResponse[lib.File](rsp, 200)
	if err != nil {
		return err
	}
	return nil
}
