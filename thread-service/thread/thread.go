package thread

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Thread struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Title string `gorm:"not null"`
	Body  string `gorm:"not null"`

	// The ID of the author of this thread.
	AuthorID uuid.UUID `gorm:"type:uuid;not null"`
	// A list of keywords that are associated with this thread.
	KeywordIDs pq.StringArray `json:"keyword_ids,omitempty" gorm:"type:uuid[]"`
}
