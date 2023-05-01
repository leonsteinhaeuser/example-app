package accessobjects

import (
	"context"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/db"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

type Keyword struct {
	db  db.Repository
	log log.Logger
}

func NewKeyword(db db.Repository, log log.Logger) *Keyword {
	return &Keyword{
		db:  db,
		log: log,
	}
}

func (a *Keyword) Create(ctx context.Context, article *lib.Keyword) error {
	return a.db.Create(ctx, article)
}

func (a *Keyword) Update(ctx context.Context, article *lib.Keyword) error {
	return a.db.Update(ctx, article, db.Selector{
		Field: "id",
		Value: article.ID,
	})
}

func (a *Keyword) Delete(ctx context.Context, article *lib.Keyword) error {
	return a.db.Delete(ctx, article, db.Selector{
		Field: "id",
		Value: article.ID,
	})
}

func (a *Keyword) Get(ctx context.Context, id string) (*lib.Keyword, error) {
	article := &lib.Keyword{}
	err := a.db.Find(ctx, article, db.Selector{
		Field: "id",
		Value: id,
	})
	if err != nil {
		return nil, err
	}
	return article, err
}

func (a *Keyword) List(ctx context.Context) ([]*lib.Keyword, error) {
	var articles []*lib.Keyword
	err := a.db.Find(ctx, &articles)
	if err != nil {
		return nil, err
	}
	return articles, err
}

func (a *Keyword) Migrate(ctx context.Context) error {
	return a.db.Migrate(ctx, &lib.Keyword{})
}
