package accessobjects

import (
	"context"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/db"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

type ArticleComment struct {
	db  db.Repository
	log log.Logger
}

func NewArticleComment(db db.Repository, log log.Logger) *ArticleComment {
	return &ArticleComment{
		db:  db,
		log: log,
	}
}

func (a *ArticleComment) Create(ctx context.Context, article *lib.ArticleComment) error {
	return a.db.Create(ctx, article)
}

func (a *ArticleComment) Update(ctx context.Context, article *lib.ArticleComment) error {
	return a.db.Update(ctx, article, db.Selector{
		Field: "id",
		Value: article.ID,
	})
}

func (a *ArticleComment) Delete(ctx context.Context, article *lib.ArticleComment) error {
	return a.db.Delete(ctx, article, db.Selector{
		Field: "id",
		Value: article.ID,
	})
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
