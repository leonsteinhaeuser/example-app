package public

import "context"

// Client is the interface for the REST client implementations.
type Client[T any] interface {
	Get(ctx context.Context, id string) (*T, error)
	List(ctx context.Context) ([]*T, error)
	Create(ctx context.Context, t *T) error
	Update(ctx context.Context, t *T) error
	Delete(ctx context.Context, t *T) error
}
