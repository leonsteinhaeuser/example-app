package db

import "context"

type TX interface {
	// Where adds a where clause to the query.
	Where(field string, value any) TX
	// Or adds an or clause to the query.
	Or(field string, value any) TX
	// Not adds a not clause to the query.
	Not(field string, value any) TX
	// Limit adds a limit clause to the query.
	Limit(limit int) TX
	// Commit executes the query.
	Commit(ctx context.Context) error
}

// Repository represents an interface between the application and the database.
type Repository interface {
	Create(ctx context.Context, data any) error
	Find(data any) TX
	Update(data any) TX
	Delete(data any) TX
	Raw(ctx context.Context, query string, args ...any) error
	Migrate(ctx context.Context, model any) error
	Close(context.Context) error
}
