package thread

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Thread is the API model for a thread.
type Thread struct {
	ID        uuid.UUID `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`

	Title string `json:"title"`
	Body  string `json:"body"`

	AuthorID   uuid.UUID   `json:"author_id"`
	KeywordIDs []uuid.UUID `json:"keyword_ids,omitempty"`
}

// ThreadModel is the database model for a thread.
type ThreadModel struct {
	ID        uuid.UUID      `json:"id,omitempty" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Title string `json:"title" gorm:"not null"`
	Body  string `json:"body" gorm:"not null"`

	// The ID of the author of this thread.
	AuthorID uuid.UUID `json:"author_id" gorm:"type:uuid;not null"`
	// A list of keywords that are associated with this thread.
	KeywordIDs pq.StringArray `json:"keyword_ids,omitempty" gorm:"type:uuid[]"`
}

func (ThreadModel) TableName() string {
	return "threads"
}

// ThreadFromThreadModel converts a ThreadModel to a Thread.
func ThreadFromThreadModel(m *ThreadModel) *Thread {
	return &Thread{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt.Time,
		Title:     m.Title,
		Body:      m.Body,
		AuthorID:  m.AuthorID,
		KeywordIDs: func() []uuid.UUID {
			uuids := make([]uuid.UUID, len(m.KeywordIDs))
			for i, id := range m.KeywordIDs {
				uuids[i] = uuid.MustParse(id)
			}
			return uuids
		}(),
	}
}

// ThreadModelFromThread converts a Thread to a ThreadModel.
func ThreadModelFromThread(t *Thread) *ThreadModel {
	return &ThreadModel{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		DeletedAt: gorm.DeletedAt{Time: t.DeletedAt},
		Title:     t.Title,
		Body:      t.Body,
		AuthorID:  t.AuthorID,
		KeywordIDs: func() pq.StringArray {
			ids := make([]string, len(t.KeywordIDs))
			for i, id := range t.KeywordIDs {
				ids[i] = id.String()
			}
			return ids
		}(),
	}
}

// ThreadsFromThreadModels converts a slice of ThreadModels to a slice of Threads.
func ThreadsFromThreadModels(models []*ThreadModel) []*Thread {
	threads := make([]*Thread, len(models))
	for i, m := range models {
		threads[i] = ThreadFromThreadModel(m)
	}
	return threads
}

// ThreadModelsFromThreads converts a slice of Threads to a slice of ThreadModels.
func ThreadModelsFromThreads(threads []*Thread) []*ThreadModel {
	models := make([]*ThreadModel, len(threads))
	for i, t := range threads {
		models[i] = ThreadModelFromThread(t)
	}
	return models
}
