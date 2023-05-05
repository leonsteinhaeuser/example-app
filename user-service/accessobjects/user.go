package accessobjects

import (
	"context"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/db"
	"github.com/leonsteinhaeuser/example-app/lib/log"
	"github.com/leonsteinhaeuser/example-app/lib/pubsub"
)

type User struct {
	pubsubClient pubsub.Client
	db           db.Repository
	log          log.Logger
}

func NewUser(db db.Repository, log log.Logger, pubsub pubsub.Client) *User {
	return &User{
		db:           db,
		log:          log,
		pubsubClient: pubsub.SetTopic("user"),
	}
}

func (a *User) Create(ctx context.Context, user *lib.User) error {
	err := a.db.Create(ctx, user)
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ResourceID: user.ID,
		ActionType: pubsub.ActionTypeCreate,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *User) Update(ctx context.Context, user *lib.User) error {
	err := a.db.Update(ctx, user, db.Selector{
		Field: "id",
		Value: user.ID,
	})
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ResourceID: user.ID,
		ActionType: pubsub.ActionTypeUpdate,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *User) Delete(ctx context.Context, user *lib.User) error {
	err := a.db.Delete(ctx, user, db.Selector{
		Field: "id",
		Value: user.ID,
	})
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ResourceID: user.ID,
		ActionType: pubsub.ActionTypeDelete,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *User) Get(ctx context.Context, id string) (*lib.User, error) {
	user := &lib.User{}
	err := a.db.Find(ctx, user, db.Selector{
		Field: "id",
		Value: id,
	})
	if err != nil {
		return nil, err
	}
	return user, err
}

func (a *User) List(ctx context.Context) ([]*lib.User, error) {
	var users []*lib.User
	err := a.db.Find(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, err
}

func (a *User) Migrate(ctx context.Context) error {
	return a.db.Migrate(ctx, &lib.User{})
}
