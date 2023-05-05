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

var (
	_ lib.Client[lib.ArticleComment] = &articleCommentClient{}
)

type articleCommentClient struct {
	log        log.Logger
	serviceURL string
	client     http.Client
}

// NewArticleCommentClient returns a new ArticleComment client
// The client configures an http.Client with a 5 second timeout
// serviceURL should be in the format of http://host:port
func NewArticleCommentClient(log log.Logger, serviceURL string) lib.CustomArticleClient[lib.ArticleComment] {
	log.Info().Logf("creating articleComment client for %s", serviceURL)
	return &articleCommentClient{
		log:        log,
		serviceURL: serviceURL,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *articleCommentClient) basePath(optionals ...string) (string, error) {
	paths := append([]string{"comment"}, optionals...)
	path, err := url.JoinPath(c.serviceURL, paths...)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (c *articleCommentClient) GetByID(ctx context.Context, id string) (*lib.ArticleComment, error) {
	path, err := c.basePath(id)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	articleComment, err := lib.CheckAndParseResponse[lib.ArticleComment](rsp, 200)
	if err != nil {
		return nil, err
	}
	return articleComment, nil
}

func (c *articleCommentClient) List(ctx context.Context) ([]lib.ArticleComment, error) {
	path, err := c.basePath()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	users, err := lib.CheckAndParseResponse[[]lib.ArticleComment](rsp, 200)
	if err != nil {
		return nil, err
	}
	return *users, nil
}

func (c *articleCommentClient) Create(ctx context.Context, articleComment lib.ArticleComment) (*lib.ArticleComment, error) {
	path, err := c.basePath()
	if err != nil {
		return nil, err
	}

	bts, err := json.Marshal(articleComment)
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

	artc, err := lib.CheckAndParseResponse[lib.ArticleComment](rsp, 200)
	if err != nil {
		return nil, err
	}
	return artc, nil
}

func (c *articleCommentClient) Update(ctx context.Context, articleComment lib.ArticleComment) (*lib.ArticleComment, error) {
	path, err := c.basePath()
	if err != nil {
		return nil, err
	}

	bts, err := json.Marshal(articleComment)
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

	artc, err := lib.CheckAndParseResponse[lib.ArticleComment](rsp, 200)
	if err != nil {
		return nil, err
	}
	return artc, nil
}

func (c *articleCommentClient) Delete(ctx context.Context, articleComment lib.ArticleComment) error {
	path, err := c.basePath(articleComment.ID.String())
	if err != nil {
		return err
	}

	bts, err := json.Marshal(articleComment)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", path, bytes.NewBuffer(bts))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	_, err = lib.CheckAndParseResponse[lib.ArticleComment](rsp, 200)
	if err != nil {
		return err
	}
	return nil
}

func (c *articleCommentClient) DeleteByArticleID(ctx context.Context, articleComment lib.ArticleComment) error {
	path, err := c.basePath("article", articleComment.ArticleID.String())
	if err != nil {
		return err
	}

	bts, err := json.Marshal(articleComment)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", path, bytes.NewBuffer(bts))
	if err != nil {
		c.log.Trace().Field("error", err).Log("error during forming request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.client.Do(req)
	if err != nil {
		c.log.Trace().Field("error", err).Log("error deleting articleComment")
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("server returned non empty reply status code: %d", rsp.StatusCode)
	}

	return nil
}
