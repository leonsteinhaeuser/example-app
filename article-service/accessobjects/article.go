package accessobjects

import (
	"context"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/db"
	"github.com/leonsteinhaeuser/example-app/lib/log"
	"github.com/leonsteinhaeuser/example-app/lib/pubsub"
)

type Article struct {
	pubsubClient pubsub.Client
	db           db.Repository
	log          log.Logger
}

func NewArticle(db db.Repository, log log.Logger, pubsub pubsub.Client) *Article {
	return &Article{
		db:           db,
		log:          log,
		pubsubClient: pubsub.SetTopic("article"),
	}
}

func (a *Article) Create(ctx context.Context, article *lib.Article) error {
	err := a.db.Create(ctx, article)
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(pubsub.Event{
		ID:     article.ID.String(),
		Action: pubsub.ActionTypeCreate,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *Article) Update(ctx context.Context, article *lib.Article) error {
	err := a.db.Update(ctx, article, db.Selector{
		Field: "id",
		Value: article.ID,
	})
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(pubsub.Event{
		ID:     article.ID.String(),
		Action: pubsub.ActionTypeUpdate,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *Article) Delete(ctx context.Context, article *lib.Article) error {
	err := a.db.Delete(ctx, article, db.Selector{
		Field: "id",
		Value: article.ID,
	})
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(pubsub.Event{
		ID:     article.ID.String(),
		Action: pubsub.ActionTypeDelete,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *Article) Get(ctx context.Context, id string) (*lib.Article, error) {
	article := &lib.Article{}
	err := a.db.Find(ctx, article, db.Selector{
		Field: "id",
		Value: id,
	})
	if err != nil {
		return nil, err
	}
	return article, err
}

func (a *Article) List(ctx context.Context) ([]*lib.Article, error) {
	var articles []*lib.Article
	err := a.db.Find(ctx, &articles)
	if err != nil {
		return nil, err
	}
	return articles, err
}

func (a *Article) Migrate(ctx context.Context) error {
	return a.db.Migrate(ctx, &lib.Article{})
}
