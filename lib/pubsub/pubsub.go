package pubsub

import (
	"context"

	"github.com/google/uuid"
)

type Event interface {
	// ID is the ID of the object that was created, updated or deleted.
	ID() uuid.UUID
	// Action is the type of the event.
	Action() ActionType
}

type Publisher interface {
	// Publish publishes a message.
	Publish(message Event) error
}

type Subscriber interface {
	// Subscribe returns a channel that will receive messages.
	// This function blocks until the subscription is ready (unbuffered channel).
	Subscribe() (<-chan SubscribeDater, error)
}

type SubscribeDater interface {
	// Ack acknowledges the message.
	Ack()
	// Data returns the message data.
	Data() []byte
}

type PublishSubscriber interface {
	Publisher
	Subscriber
}

type Client interface {
	PublishSubscriber
	// WithTopic creates a new Client with the given topic.
	SetTopic(topic string) Client
	// Close closes the connection.
	Close(context.Context) error
}

type ActionType string

const (
	ActionTypeCreate ActionType = "create"
	ActionTypeUpdate ActionType = "update"
	ActionTypeDelete ActionType = "delete"
)

type DefaultEvent struct {
	// ResourceID is the ResourceID of the object that was created, updated or deleted.
	ResourceID uuid.UUID
	// ActionType is the type of the event.
	ActionType ActionType
}

// ID is the ID of the object that was created, updated or deleted.
func (d *DefaultEvent) ID() uuid.UUID {
	return d.ResourceID
}

// Action is the type of the event.
func (d *DefaultEvent) Action() ActionType {
	return d.ActionType
}
