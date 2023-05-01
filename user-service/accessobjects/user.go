package accessobjects

import (
	"context"

	"github.com/leonsteinhaeuser/example-app/lib"
	"github.com/leonsteinhaeuser/example-app/lib/db"
	"github.com/leonsteinhaeuser/example-app/lib/log"
)

type User struct {
	db  db.Repository
	log log.Logger
}

func NewUser(db db.Repository, log log.Logger) *User {
	return &User{
		db:  db,
		log: log,
	}
}

func (a *User) Create(ctx context.Context, article *lib.User) error {
	return a.db.Create(ctx, article)
}

func (a *User) Update(ctx context.Context, article *lib.User) error {
	return a.db.Update(ctx, article, db.Selector{
		Field: "id",
		Value: article.ID,
	})
}

func (a *User) Delete(ctx context.Context, article *lib.User) error {
	return a.db.Delete(ctx, article, db.Selector{
		Field: "id",
		Value: article.ID,
	})
}

func (a *User) Get(ctx context.Context, id string) (*lib.User, error) {
	article := &lib.User{}
	err := a.db.Find(ctx, article, db.Selector{
		Field: "id",
		Value: id,
	})
	if err != nil {
		return nil, err
	}
	return article, err
}

func (a *User) List(ctx context.Context) ([]*lib.User, error) {
	var articles []*lib.User
	err := a.db.Find(ctx, &articles)
	if err != nil {
		return nil, err
	}
	return articles, err
}

func (a *User) Migrate(ctx context.Context) error {
	return a.db.Migrate(ctx, &lib.User{})
}
