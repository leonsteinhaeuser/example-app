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
	_ lib.Client[lib.User] = &userClient{}
)

type userClient struct {
	log        log.Logger
	serviceURL string
	client     http.Client
}

// NewUserClient returns a new UserClient
// The client configures an http.Client with a 5 second timeout
// serviceURL should be in the format of http://host:port
func NewUserClient(log log.Logger, serviceURL string) lib.Client[lib.User] {
	log.Info().Logf("creating user client for %s", serviceURL)
	return &userClient{
		log:        log,
		serviceURL: serviceURL,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *userClient) basePath(optionals ...string) (string, error) {
	paths := append([]string{"user"}, optionals...)
	path, err := url.JoinPath(c.serviceURL, paths...)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (c *userClient) GetByID(ctx context.Context, id string) (*lib.User, error) {
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

	user, err := lib.CheckAndParseResponse[lib.User](rsp, 200)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *userClient) List(ctx context.Context) ([]lib.User, error) {
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

	users, err := lib.CheckAndParseResponse[[]lib.User](rsp, 200)
	if err != nil {
		return nil, err
	}
	return *users, nil
}

func (c *userClient) Create(ctx context.Context, user lib.User) (*lib.User, error) {
	path, err := c.basePath()
	if err != nil {
		return nil, err
	}

	bts, err := json.Marshal(user)
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

	artc, err := lib.CheckAndParseResponse[lib.User](rsp, 200)
	if err != nil {
		return nil, err
	}
	return artc, nil
}

func (c *userClient) Update(ctx context.Context, user lib.User) (*lib.User, error) {
	path, err := c.basePath()
	if err != nil {
		return nil, err
	}

	bts, err := json.Marshal(user)
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

	artc, err := lib.CheckAndParseResponse[lib.User](rsp, 200)
	if err != nil {
		return nil, err
	}
	return artc, nil
}

func (c *userClient) Delete(ctx context.Context, user lib.User) error {
	path, err := c.basePath(user.ID.String())
	if err != nil {
		return err
	}

	bts, err := json.Marshal(user)
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

	_, err = lib.CheckAndParseResponse[lib.User](rsp, 200)
	if err != nil {
		return err
	}
	return nil
}
