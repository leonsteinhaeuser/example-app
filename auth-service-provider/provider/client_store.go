package provider

import (
	"context"
	"fmt"

	"github.com/leonsteinhaeuser/example-app/internal/db"
)

type ClientStore struct {
	DB db.Repository
}

func (s *ClientStore) GetClientByID(ctx context.Context, id string) (*Client, error) {
	client := &Client{}
	err := s.DB.Find(client).Where("id = ?", id).Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find client by ID %q: %w", id, err)
	}
	return client, nil
}

func (s *ClientStore) Create(ctx context.Context, client *Client) error {
	err := s.DB.Create(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	return nil
}

func (s *ClientStore) Update(ctx context.Context, client *Client) error {
	err := s.DB.Update(client).Where("id = ?", client.ID).Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}
	return nil
}

func (s *ClientStore) Delete(ctx context.Context, id string) error {
	err := s.DB.Delete(&Client{}).Where("id = ?", id).Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}
	return nil
}
