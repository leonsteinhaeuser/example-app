package accessobjects

import (
	"context"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/db"
	"github.com/leonsteinhaeuser/example-app/lib/log"
	"github.com/leonsteinhaeuser/example-app/lib/pubsub"
)

type ArticleComment struct {
	pubsubClient pubsub.Client
	db           db.Repository
	log          log.Logger
}

func NewArticleComment(db db.Repository, log log.Logger, pubsub pubsub.Client) *ArticleComment {
	return &ArticleComment{
		db:           db,
		log:          log,
		pubsubClient: pubsub.SetTopic("articlecomment"),
	}
}

func (a *ArticleComment) Create(ctx context.Context, article *lib.ArticleComment) error {
	err := a.db.Create(ctx, article)
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ResourceID: article.ID,
		ActionType: pubsub.ActionTypeCreate,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticleComment) Update(ctx context.Context, article *lib.ArticleComment) error {
	err := a.db.Update(ctx, article, db.Selector{
		Field: "id",
		Value: article.ID,
	})
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ResourceID: article.ID,
		ActionType: pubsub.ActionTypeUpdate,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticleComment) Delete(ctx context.Context, article *lib.ArticleComment) error {
	err := a.db.Delete(ctx, article, db.Selector{
		Field: "id",
		Value: article.ID,
	})
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ResourceID: article.ID,
		ActionType: pubsub.ActionTypeDelete,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticleComment) DeleteByArticleID(ctx context.Context, article *lib.ArticleComment) error {
	err := a.db.Delete(ctx, article, db.Selector{
		Field: "article_id",
		Value: article.ArticleID,
	})
	if err != nil {
		return err
	}
	err = a.pubsubClient.Publish(&pubsub.DefaultEvent{
		ActionType: pubsub.ActionTypeDelete,
		AdditionalData: map[string]interface{}{
			"article_id": article.ArticleID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticleComment) Get(ctx context.Context, id string) (*lib.ArticleComment, error) {
	article := &lib.ArticleComment{}
	err := a.db.Find(ctx, article, db.Selector{
		Field: "id",
		Value: id,
	})
	if err != nil {
		return nil, err
	}
	return article, err
}

func (a *ArticleComment) List(ctx context.Context) ([]*lib.ArticleComment, error) {
	var articles []*lib.ArticleComment
	err := a.db.Find(ctx, &articles)
	if err != nil {
		return nil, err
	}
	return articles, err
}

func (a *ArticleComment) Migrate(ctx context.Context) error {
	return a.db.Migrate(ctx, &lib.ArticleComment{})
}
