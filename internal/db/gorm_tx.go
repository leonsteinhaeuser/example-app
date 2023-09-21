package db

import (
	"context"

	"gorm.io/gorm"
)

var (
	_ TX = (*gormTX)(nil)
)

type gormAction int

const (
	gormActionFind gormAction = iota
	gormActionUpdate
	gormActionDelete
)

type gormTX struct {
	tx         *gorm.DB
	receiver   any
	gormAction gormAction
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
	switch g.gormAction {
	case gormActionFind:
		return g.tx.Find(g.receiver).Error
	case gormActionUpdate:
		return g.tx.Updates(g.receiver).Error
	case gormActionDelete:
		return g.tx.Delete(g.receiver).Error
	}
	return nil
}
