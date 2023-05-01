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
	_ lib.Client[lib.Keyword] = &keywordClient{}
)

type keywordClient struct {
	log        log.Logger
	serviceURL string
	client     http.Client
}

// NewKeywordClient returns a new KeywordClient
// The client configures an http.Client with a 5 second timeout
// serviceURL should be in the format of http://host:port
func NewKeywordClient(log log.Logger, serviceURL string) lib.Client[lib.Keyword] {
	log.Info().Logf("creating keyword client for %s", serviceURL)
	return &keywordClient{
		log:        log,
		serviceURL: serviceURL,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *keywordClient) basePath(optionals ...string) (string, error) {
	paths := append([]string{"keyword"}, optionals...)
	path, err := url.JoinPath(c.serviceURL, paths...)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (c *keywordClient) GetByID(ctx context.Context, id string) (*lib.Keyword, error) {
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

	keyword, err := lib.CheckAndParseResponse[lib.Keyword](rsp, 200)
	if err != nil {
		return nil, err
	}
	return keyword, nil
}

func (c *keywordClient) List(ctx context.Context) ([]lib.Keyword, error) {
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

	users, err := lib.CheckAndParseResponse[[]lib.Keyword](rsp, 200)
	if err != nil {
		return nil, err
	}
	return *users, nil
}

func (c *keywordClient) Create(ctx context.Context, keyword lib.Keyword) (*lib.Keyword, error) {
	path, err := c.basePath()
	if err != nil {
		return nil, err
	}

	bts, err := json.Marshal(keyword)
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

	artc, err := lib.CheckAndParseResponse[lib.Keyword](rsp, 200)
	if err != nil {
		return nil, err
	}
	return artc, nil
}

func (c *keywordClient) Update(ctx context.Context, keyword lib.Keyword) (*lib.Keyword, error) {
	path, err := c.basePath()
	if err != nil {
		return nil, err
	}

	bts, err := json.Marshal(keyword)
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

	artc, err := lib.CheckAndParseResponse[lib.Keyword](rsp, 200)
	if err != nil {
		return nil, err
	}
	return artc, nil
}

func (c *keywordClient) Delete(ctx context.Context, keyword lib.Keyword) error {
	path, err := c.basePath(keyword.ID.String())
	if err != nil {
		return err
	}

	bts, err := json.Marshal(keyword)
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

	_, err = lib.CheckAndParseResponse[lib.Keyword](rsp, 200)
	if err != nil {
		return err
	}
	return nil
}
