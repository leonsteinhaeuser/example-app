package accessobjects

import (
	"context"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/db"
	"github.com/leonsteinhaeuser/example-app/lib/log"
	"github.com/leonsteinhaeuser/example-app/lib/pubsub"
)

type File struct {
	pubsubClient pubsub.Client
	db           db.Repository
	log          log.Logger
}

func NewFile(db db.Repository, log log.Logger, pubsub pubsub.Client) *File {
	return &File{
		db:           db,
		log:          log,
		pubsubClient: pubsub.SetTopic("file"),
	}
}

func (a *File) Create(ctx context.Context, file *lib.File) error {
	err := a.db.Create(ctx, file)
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ResourceID: file.ID,
		ActionType: pubsub.ActionTypeCreate,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *File) Update(ctx context.Context, file *lib.File) error {
	err := a.db.Update(ctx, file, db.Selector{
		Field: "id",
		Value: file.ID,
	})
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ResourceID: file.ID,
		ActionType: pubsub.ActionTypeUpdate,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *File) Delete(ctx context.Context, file *lib.File) error {
	err := a.db.Delete(ctx, file, db.Selector{
		Field: "id",
		Value: file.ID,
	})
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ResourceID: file.ID,
		ActionType: pubsub.ActionTypeDelete,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *File) Get(ctx context.Context, id string) (*lib.File, error) {
	file := &lib.File{}
	err := a.db.Find(ctx, file, db.Selector{
		Field: "id",
		Value: id,
	})
	if err != nil {
		return nil, err
	}
	return file, err
}

func (a *File) List(ctx context.Context) ([]*lib.File, error) {
	var files []*lib.File
	err := a.db.Find(ctx, &files)
	if err != nil {
		return nil, err
	}
	return files, err
}

func (a *File) Migrate(ctx context.Context) error {
	return a.db.Migrate(ctx, &lib.File{})
}
