package accessobjects

import (
	"context"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/db"
	"github.com/leonsteinhaeuser/example-app/lib/log"
	"github.com/leonsteinhaeuser/example-app/lib/pubsub"
)

type Keyword struct {
	pubsubClient pubsub.Client
	db           db.Repository
	log          log.Logger
}

func NewKeyword(db db.Repository, log log.Logger, pubsub pubsub.Client) *Keyword {
	return &Keyword{
		db:           db,
		log:          log,
		pubsubClient: pubsub.SetTopic("keyword"),
	}
}

func (a *Keyword) Create(ctx context.Context, keyword *lib.Keyword) error {
	err := a.db.Create(ctx, keyword)
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ResourceID: keyword.ID,
		ActionType: pubsub.ActionTypeCreate,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *Keyword) Update(ctx context.Context, keyword *lib.Keyword) error {
	err := a.db.Update(ctx, keyword, db.Selector{
		Field: "id",
		Value: keyword.ID,
	})
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ResourceID: keyword.ID,
		ActionType: pubsub.ActionTypeUpdate,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *Keyword) Delete(ctx context.Context, keyword *lib.Keyword) error {
	err := a.db.Delete(ctx, keyword, db.Selector{
		Field: "id",
		Value: keyword.ID,
	})
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ResourceID: keyword.ID,
		ActionType: pubsub.ActionTypeDelete,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *Keyword) Get(ctx context.Context, id string) (*lib.Keyword, error) {
	keyword := &lib.Keyword{}
	err := a.db.Find(ctx, keyword, db.Selector{
		Field: "id",
		Value: id,
	})
	if err != nil {
		return nil, err
	}
	return keyword, err
}

func (a *Keyword) List(ctx context.Context) ([]*lib.Keyword, error) {
	var keywords []*lib.Keyword
	err := a.db.Find(ctx, &keywords)
	if err != nil {
		return nil, err
	}
	return keywords, err
}

func (a *Keyword) Migrate(ctx context.Context) error {
	return a.db.Migrate(ctx, &lib.Keyword{})
}
