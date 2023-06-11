package thread

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/leonsteinhaeuser/example-app/internal/db"
	"github.com/leonsteinhaeuser/example-app/internal/pubsub"
)

type Store struct {
	DB db.Repository
	PS pubsub.Publisher
}

func (s *Store) Create(ctx context.Context, t *Thread) error {
	err := s.DB.Create(ctx, ThreadModelFromThread(t))
	if err != nil {
		return fmt.Errorf("unable to create thread: %w", err)
	}
	s.PS.Publish(&pubsub.DefaultEvent{
		ResourceID:     t.ID,
		ActionType:     pubsub.ActionTypeCreate,
		AdditionalData: map[string]interface{}{},
	})
	return nil
}

func (s *Store) GetByID(ctx context.Context, id string) (*Thread, error) {
	t := &ThreadModel{}
	err := s.DB.Find(t).Where("id = ?", id).Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to find thread by ID %q: %w", id, err)
	}
	return ThreadFromThreadModel(t), nil
}

func (s *Store) List(ctx context.Context) ([]*Thread, error) {
	var threads []*ThreadModel
	err := s.DB.Find(&threads).Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list threads: %w", err)
	}
	return ThreadsFromThreadModels(threads), nil
}

func (s *Store) UpdateByID(ctx context.Context, t *Thread) error {
	err := s.DB.Update(ThreadModelFromThread(t)).Where("id = ?", t.ID).Commit(ctx)
	if err != nil {
		return fmt.Errorf("unable to update thread: %w", err)
	}
	s.PS.Publish(&pubsub.DefaultEvent{
		ResourceID:     t.ID,
		ActionType:     pubsub.ActionTypeUpdate,
		AdditionalData: map[string]interface{}{},
	})
	return nil
}

func (s *Store) DeleteByID(ctx context.Context, id string) error {
	err := s.DB.Delete(&ThreadModel{}).Where("id = ?", id).Commit(ctx)
	if err != nil {
		return fmt.Errorf("unable to delete thread: %w", err)
	}
	s.PS.Publish(&pubsub.DefaultEvent{
		ResourceID:     uuid.MustParse(id),
		ActionType:     pubsub.ActionTypeDelete,
		AdditionalData: map[string]interface{}{},
	})
	return nil
}
