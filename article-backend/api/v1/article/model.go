package article

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID        uuid.UUID  `json:"id,omitempty" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Title is the title of the article.
	Title string `json:"title,omitempty"`
	// Description is the description of the article.
	Description string `json:"description,omitempty"`
	// Content is the content of the article.
	Content string `json:"content,omitempty"`
	// Published is a flag indicating whether the article is published or not.
	Published bool `json:"published"`
	// PublishedAt is the time the article was published.
	PublishedAt *time.Time `json:"published_at,omitempty"`
	// PublishedBy is the ID of the user who published the article.
	PublishedBy *uuid.UUID `json:"published_by,omitempty"`

	// Tags is a list of tags of the article.
	Tags []string `json:"tags,omitempty" gorm:"serializer:json"`

	// AuthorID is the ID of the head author of the article.
	AuthorID uuid.UUID `json:"author_id,omitempty"`
	// CoAuthorIDs is a list of IDs of co-authors of the article.
	CoAuthorIDs []uuid.UUID `json:"co_author_ids,omitempty" gorm:"serializer:json"`
}
