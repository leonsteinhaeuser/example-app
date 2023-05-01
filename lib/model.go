package lib

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type NumberResponse struct {
	Number int64 `json:"number"`
}

type Article struct {
	// identifier and state fields
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// model fields

	// Title reprGormDBDataTypeesents the title of the article
	Title string `json:"title,omitempty"`
	// Description represents a short description what the article is about
	Description string `json:"description,omitempty"`
	// Author represents the author of the article
	AuthorID uuid.UUID `json:"author_id,omitempty" gorm:"type:uuid"`
	// Text is the actual article text
	Text string `json:"text,omitempty"`

	KeywordIDs pq.StringArray `json:"keyword_ids,omitempty" gorm:"type:uuid[]"`
}

type User struct {
	// identifier and state fields
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// model fields
	Username    string `json:"username,omitempty" gorm:"not null;unique"`
	DisplayName string `json:"display_name,omitempty"`
	Password    string `json:"password,omitempty" gorm:"not null"`
	Email       string `json:"email,omitempty" gorm:"not null;unique"`
	IsConfirmed bool   `json:"is_confirmed,omitempty"`
	FistName    string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
}
