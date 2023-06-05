package keystore

import (
	"context"
	"time"
)

type KeyStore interface {
	Geter
	Seter
	Delete(ctx context.Context, key string) error
}

type Geter interface {
	Get(ctx context.Context, key string) ([]byte, error)
}

type Seter interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
}
