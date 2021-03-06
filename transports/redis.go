package transports

import (
	"fmt"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/keimoon/gore"
)

type RedisTransport struct {
	c   *gore.Conn
	url *url.URL
}

func NewRedisTransport(u *url.URL) (*RedisTransport, error) {
	c, err := gore.Dial(u.Host)
	if err != nil {
		return nil, err
	}

	if u.User != nil {
		if pass, ok := u.User.Password(); ok {
			err := c.Auth(pass)
			if err != nil {
				return nil, err
			}
		}
	}

	return &RedisTransport{
		c:   c,
		url: u,
	}, nil
}

func (t *RedisTransport) Describe() string {
	return SafeURLString(t.url)
}

func (t *RedisTransport) Subscribe(topic string) (Subscription, error) {
	if !t.c.IsConnected() {
		return nil, fmt.Errorf("not connected")
	}

	log.
		WithField("topic", topic).
		WithField("transport", t.Describe()).
		Debug("Creating new subscriber for transport")

	conn, err := gore.Dial(t.c.GetAddress())
	if err != nil {
		log.
			WithField("address", t.c.GetAddress()).
			WithField("transport", t.Describe()).
			WithError(err).
			Warn("Failed to create connection for subscriber")
		return nil, err
	}

	if t.url.User != nil {
		if pass, ok := t.url.User.Password(); ok {
			err := conn.Auth(pass)
			if err != nil {
				return nil, err
			}
		}
	}

	s := gore.NewSubscriptions(conn)

	err = s.Subscribe(GetFullTopic(t.url, topic))
	if err != nil {
		log.
			WithField("topic", topic).
			WithField("transport", t.Describe()).
			WithError(err).
			Warn("Failed to subscribe to topic for subscriber")

		s.Close()
		conn.Close()
		return nil, err
	}

	c := make(chan []byte)
	go func() {
		for m := range s.Message() {
			log.
				WithField("topic", topic).
				WithField("transport", t.Describe()).
				WithField("entry", m).
				Debug("Received entry for subscription")

			if m == nil {
				break
			}

			c <- m.Message
		}

		log.
			WithField("topic", topic).
			WithField("transport", t.Describe()).
			Debug("Closing subscriber")

		close(c)
	}()

	return &redisSubscription{
		conn: conn,
		sub:  s,
		c:    c,
	}, nil
}

func (t *RedisTransport) Publish(topic string, data []byte) error {
	if !t.c.IsConnected() {
		return fmt.Errorf("not connected")
	}

	log.
		WithField("topic", topic).
		WithField("transport", t.Describe()).
		Debug("Publishing message to transport")

	return gore.Publish(t.c, GetFullTopic(t.url, topic), data)
}

func (t *RedisTransport) Close() error {
	return t.c.Close()
}

type redisSubscription struct {
	conn *gore.Conn
	sub  *gore.Subscriptions
	c    chan []byte
}

func (s *redisSubscription) Channel() <-chan []byte {
	return s.c
}

func (s *redisSubscription) Close() error {
	defer s.sub.Close()
	return s.conn.Close()
}
