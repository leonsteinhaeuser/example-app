package pubsub

import "context"

type Publisher interface {
	// Publish publishes a message.
	Publish(message any) error
}

type Subscriber interface {
	// Subscribe returns a channel that will receive messages.
	// This function blocks until the subscription is ready (unbuffered channel).
	Subscribe() (<-chan SubscribeDataer, error)
}

type SubscribeDataer interface {
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

type Event struct {
	// ID is the ID of the object that was created, updated or deleted.
	ID string
	// Action is the type of the event.
	Action ActionType
}
