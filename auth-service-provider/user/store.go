package user

import (
	"context"
	"fmt"

	"github.com/leonsteinhaeuser/example-app/internal/db"
)

// Store represents the storage interface to the user store
type Store interface {
	GetUserByID(context.Context, string) (*User, error)
	GetUserByLogin(context.Context, string) (*User, error)
}

var (
	_ Store = (*UserStoreDB)(nil)
)

type UserStoreDB struct {
	DB db.Repository
}

func (u UserStoreDB) GetUserByID(ctx context.Context, id string) (*User, error) {
	user := &User{}
	err := u.DB.Find(user).Where("id", id).Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id %q: %w", id, err)
	}
	return user, nil
}

func (u UserStoreDB) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	user := &User{}
	err := u.DB.Find(user).Where("username", login).Or("email", login).Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by login %q: %w", login, err)
	}
	return user, nil
}
