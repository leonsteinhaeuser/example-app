package user

import (
	"context"
	"fmt"

	"github.com/leonsteinhaeuser/example-app/internal/db"
)

// Store represents the storage interface to the user store
type Store struct {
	DB db.Repository
}

func (u Store) GetByID(ctx context.Context, id string) (*User, error) {
	user := &User{}
	err := u.DB.Find(user).Where("id", id).Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id %q: %w", id, err)
	}
	return user, nil
}

func (u Store) GetByLogin(ctx context.Context, login string) (*User, error) {
	user := &User{}
	err := u.DB.Find(user).Where("username", login).Or("email", login).Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by login %q: %w", login, err)
	}
	return user, nil
}

func (u Store) Create(ctx context.Context, user *User) error {
	err := u.DB.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (u Store) Update(ctx context.Context, user *User) error {
	err := u.DB.Update(user).Where("id", user.ID).Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (u Store) Delete(ctx context.Context, id string) error {
	err := u.DB.Delete(&User{}).Where("id", id).Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (u Store) List(ctx context.Context) ([]*User, error) {
	users := []*User{}
	err := u.DB.Find(&users).Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}
