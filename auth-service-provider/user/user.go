package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type User struct {
	// ID is the unique identifier of the user
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	// Username is the unique username of the user
	Username string
	// Password is the password of the user
	Password string
	// Email is the email address of the user it can be used as an alternative login
	Email         string
	EmailVerified bool

	// Optional fields
	FirstName string
	LastName  string

	PreferredLanguage language.Tag
}
