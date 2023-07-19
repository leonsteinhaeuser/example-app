package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leonsteinhaeuser/example-app/internal/client"
	"github.com/leonsteinhaeuser/example-app/internal/log"
	v1 "github.com/leonsteinhaeuser/example-app/number-service/api/v1"
)

type RestClient struct {
	log log.Logger

	client client.HttpClient
}

func NewRestClient(log log.Logger, timeoutSec int, address string) *RestClient {
	return &RestClient{
		log: log,
		client: client.HttpClient{
			Client:  &http.Client{},
			Address: address,
		},
	}
}

func (r *RestClient) GetNumber(ctx context.Context) (*v1.NumberResponse, error) {
	rsp, err := r.client.Request(ctx, http.MethodGet, "", nil)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	// check response
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("")
	}

	number := &v1.NumberResponse{}
	err = json.NewDecoder(rsp.Body).Decode(number)
	if err != nil {
		return nil, err
	}

	return number, nil
}
