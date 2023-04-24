package accessobjects

import (
	"context"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/db"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

type Article struct {
	db  db.Repository
	log log.Logger
}

func NewArticle(db db.Repository, log log.Logger) *Article {
	return &Article{
		db:  db,
		log: log,
	}
}

func (a *Article) Create(ctx context.Context, article *lib.Article) error {
	return a.db.Create(ctx, article)
}

func (a *Article) Update(ctx context.Context, article *lib.Article) error {
	return a.db.Update(ctx, article, db.Selector{
		Field: "id",
		Value: article.ID,
	})
}

func (a *Article) Delete(ctx context.Context, article *lib.Article) error {
	return a.db.Delete(ctx, article, db.Selector{
		Field: "id",
		Value: article.ID,
	})
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
