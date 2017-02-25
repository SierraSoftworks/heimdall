package transports

import (
	"fmt"

	"github.com/nats-io/go-nats"
)

type NATSTransport struct {
	QueueGroup string

	url *TransportURL
	c   *nats.Conn
}

func NewNATSTransport(u *TransportURL) (*NATSTransport, error) {
	c, err := nats.Connect(u.String())
	if err != nil {
		return nil, err
	}

	return &NATSTransport{
		QueueGroup: u.GetOption("queue_group", "heimdall_servers"),

		c:   c,
		url: u,
	}, nil
}

func (t *NATSTransport) Describe() string {
	return t.url.SafeString()
}

func (t *NATSTransport) Subscribe(topic string) (Subscription, error) {
	if t.c == nil || t.c.IsClosed() {
		return nil, fmt.Errorf("not connected")
	}

	c := make(chan []byte)
	s, err := t.c.QueueSubscribe(t.url.GetFullTopic(topic), t.QueueGroup, func(m *nats.Msg) {
		c <- m.Data
	})

	if err != nil {
		return nil, err
	}

	return &natsSubscription{
		c:   c,
		sub: s,
	}, nil
}

func (t *NATSTransport) Publish(topic string, data []byte) error {
	if t.c == nil || t.c.IsClosed() {
		return fmt.Errorf("not connected")
	}

	return t.c.Publish(t.url.GetFullTopic(topic), data)
}

func (t *NATSTransport) Close() error {
	t.c.Close()
	return nil
}

type natsSubscription struct {
	c   chan []byte
	sub *nats.Subscription
}

func (s *natsSubscription) Channel() <-chan []byte {
	return s.c
}

func (s *natsSubscription) Close() error {
	defer close(s.c)
	return s.sub.Unsubscribe()
}
