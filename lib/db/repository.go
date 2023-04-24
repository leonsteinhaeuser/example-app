package db

import "context"

// Selector represents a field and a value to select a specific row in the database.
// The underlying database implementation should use the field and value to create a WHERE clause.
//
// Example:
//
//		selector := Selector{
//			Field: "id",
//			Value: 1,
//		}
//
//	 The above selector should create a WHERE clause like this:
//		WHERE id = 1
type Selector struct {
	Field string
	Value any
}

// Repository represents an interface between the application and the database.
type Repository interface {
	Create(ctx context.Context, data any) error
	Find(ctx context.Context, data any, selectors ...Selector) error
	Update(ctx context.Context, data any, selectors ...Selector) error
	Delete(ctx context.Context, data any, selectors ...Selector) error
	Migrate(ctx context.Context, model any) error
	Close(context.Context) error
}
