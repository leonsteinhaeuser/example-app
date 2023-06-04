package db

import (
	"context"

	"gorm.io/gorm"
)

var (
	_ TX = (*gormTX)(nil)
)

type gormTX struct {
	tx *gorm.DB
}

func (g *gormTX) Where(field string, value any) TX {
	g.tx = g.tx.Where(field, value)
	return g
}

func (g *gormTX) Or(field string, value any) TX {
	g.tx = g.tx.Or(field, value)
	return g
}

func (g *gormTX) Not(field string, value any) TX {
	g.tx = g.tx.Not(field, value)
	return g
}

func (g *gormTX) Limit(limit int) TX {
	g.tx = g.tx.Limit(limit)
	return g
}

func (g *gormTX) Commit(ctx context.Context) error {
	return g.tx.Error
}
