package user

import (
	"context"

	"github.com/leonsteinhaeuser/example-app/internal/db"
)

// UserStore represents the storage interface to the user store
type UserStore interface {
	GetUserByID(context.Context, string) (*User, error)
	GetUserByLogin(context.Context, string) (*User, error)
}

var (
	_ UserStore = (*UserStoreDB)(nil)
)

type UserStoreDB struct {
	DB db.Repository
}

func (u UserStoreDB) GetUserByID(ctx context.Context, id string) (*User, error) {
	return nil, nil
}

func (u UserStoreDB) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	return nil, nil
}
