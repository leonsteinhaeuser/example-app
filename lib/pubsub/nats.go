package pubsub

import (
	"context"

	"github.com/leonsteinhaeuser/example-app/lib/log"
	"github.com/nats-io/nats.go"
)

type natsResponse struct {
	msg *nats.Msg
}

func (n *natsResponse) Ack() {
	n.msg.Ack()
}

func (n *natsResponse) Data() []byte {
	return n.msg.Data
}

type natsClient struct {
	log log.Logger

	topic    string
	origConn *nats.Conn
	conn     *nats.EncodedConn
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
		log:      log,
		topic:    defaultTopic,
		origConn: conn,
		conn:     ec,
	}, nil
}

func (n *natsClient) SetTopic(topic string) Client {
	n.log.Debug().Field("topic", topic).Log("setting topic")
	natsClients2 := &natsClient{}
	*natsClients2 = *n
	natsClients2.topic = topic
	return natsClients2
}

func (n *natsClient) Close(context.Context) error {
	n.conn.Close()
	return nil
}

func (n *natsClient) Publish(message Event) error {
	return n.conn.Publish(n.topic, message)
}

func (n *natsClient) Subscribe() (<-chan SubscribeDater, error) {
	streamCh := make(chan SubscribeDater)
	_, err := n.conn.Subscribe(n.topic, func(msg *nats.Msg) {
		// add message to channel
		streamCh <- &natsResponse{msg}
	})
	if err != nil {
		return nil, err
	}
	return streamCh, nil
}
