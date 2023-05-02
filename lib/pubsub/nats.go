package pubsub

import (
	"context"

	"github.com/leonsteinhaeuser/example-app/lib/log"
	"github.com/nats-io/nats.go"
)

type natsResponse struct {
	*nats.Msg
}

func (n *natsResponse) Ack() {
	n.Ack()
}

func (n *natsResponse) Data() any {
	return n.Data
}

type natsClient struct {
	log log.Logger

	topic string
	conn  *nats.EncodedConn
}

func NewNatsClient(log log.Logger, connectStr string, defaultTopic string) (Client, error) {
	conn, err := nats.Connect(connectStr)
	if err != nil {
		return nil, err
	}
	ec, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	return &natsClient{
		log:   log,
		topic: defaultTopic,
		conn:  ec,
	}, nil
}

func (n *natsClient) SetTopic(topic string) Client {
	natsClients2 := &natsClient{}
	*natsClients2 = *n
	natsClients2.topic = topic
	return natsClients2
}

func (n *natsClient) Close(context.Context) error {
	n.conn.Close()
	return nil
}

func (n *natsClient) Publish(message any) error {
	return n.conn.Publish(n.topic, message)
}

func (n *natsClient) Subscribe() (<-chan SubscribeDataer, error) {
	streamCh := make(chan SubscribeDataer)
	_, err := n.conn.Subscribe(n.topic, func(msg *nats.Msg) {
		streamCh <- &natsResponse{msg}
	})
	if err != nil {
		return nil, err
	}
	return streamCh, nil
}
