package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

type ClientError struct {
	StatusCode int
	Message    string
	Article    *lib.Article
	Articles   []lib.Article
}

func (c *ClientError) Error() string {
	return fmt.Sprintf("client error: %d - %s", c.StatusCode, c.Message)
}

var (
	_ lib.Client[lib.Article] = &articleClient{}
)

type articleClient struct {
	log        log.Logger
	serviceURL string
	client     http.Client
}

// NewArticleClient returns a new ArticleClient
// The client configures an http.Client with a 5 second timeout
// serviceURL should be in the format of http://host:port
func NewArticleClient(log log.Logger, serviceURL string) lib.Client[lib.Article] {
	log.Info().Logf("creating article client for %s", serviceURL)
	return &articleClient{
		log:        log,
		serviceURL: serviceURL,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *articleClient) basePath(optionals ...string) (string, error) {
	paths := append([]string{"article"}, optionals...)
	path, err := url.JoinPath(c.serviceURL, paths...)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (c *articleClient) GetByID(ctx context.Context, id string) (*lib.Article, error) {
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

	article, err := lib.CheckAndParseResponse[lib.Article](rsp, 200)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (c *articleClient) List(ctx context.Context) ([]lib.Article, error) {
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

	articles, err := lib.CheckAndParseResponse[[]lib.Article](rsp, 200)
	if err != nil {
		return nil, err
	}
	return *articles, nil
}

func (c *articleClient) Create(ctx context.Context, article lib.Article) (*lib.Article, error) {
	path, err := c.basePath()
	if err != nil {
		return nil, err
	}

	bts, err := json.Marshal(article)
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

	artc, err := lib.CheckAndParseResponse[lib.Article](rsp, 200)
	if err != nil {
		return nil, err
	}
	return artc, nil
}

func (c *articleClient) Update(ctx context.Context, article lib.Article) (*lib.Article, error) {
	path, err := c.basePath()
	if err != nil {
		return nil, err
	}

	bts, err := json.Marshal(article)
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

	artc, err := lib.CheckAndParseResponse[lib.Article](rsp, 200)
	if err != nil {
		return nil, err
	}
	return artc, nil
}

func (c *articleClient) Delete(ctx context.Context, article lib.Article) error {
	path, err := c.basePath(article.ID.String())
	if err != nil {
		return err
	}

	bts, err := json.Marshal(article)
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

	_, err = lib.CheckAndParseResponse[lib.Article](rsp, 200)
	if err != nil {
		return err
	}
	return nil
}
