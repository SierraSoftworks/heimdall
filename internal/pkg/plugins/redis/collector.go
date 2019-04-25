package redis

import (
	"path"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"

	log "github.com/Sirupsen/logrus"
	"github.com/keimoon/gore"
)

type RedisCollector struct {
	conf   *driver.Driver
	logger *log.Entry
}

func NewRedisCollector(cfg *driver.Driver) (plugins.Collector, error) {
	logger := log.WithFields(log.Fields{
		"driver":    cfg.Type,
		"collector": cfg.SafeURLString(),
		"url":       cfg.URL,
	})

	return &RedisCollector{
		conf:   cfg,
		logger: logger,
	}, nil
}

func (c *RedisCollector) Driver() *driver.Driver {
	return c.conf
}

func (c *RedisCollector) Describe() string {
	return c.conf.SafeURLString()
}

func (c *RedisCollector) Subscribe(channel string) (plugins.Subscription, error) {
	logger := c.logger.WithField("channel", channel)

	logger.Debug("Connecting to Redis")
	conn, err := gore.Dial(c.conf.URL.Host)
	if err != nil {
		logger.
			WithError(err).
			Error("Failed to connect to Redis")
		return nil, err
	}

	logger.Debug("Connected to Redis")

	if c.conf.URL.User != nil {
		if pass, ok := c.conf.URL.User.Password(); ok {
			logger.Debug("Authenticating to Redis")
			err := conn.Auth(pass)
			if err != nil {
				logger.
					WithError(err).
					Error("Failed to authenticate to Redis")
				return nil, err
			}
		}
	}

	logger = logger.WithField("path", path.Join(c.conf.URL.Path, channel))

	logger.Debug("Subscribing to path")
	s := gore.NewSubscriptions(conn)

	if err := s.Subscribe(path.Join(c.conf.URL.Path, channel)); err != nil {
		logger.WithError(err).Error("Failed to subscribe to topic")

		s.Close()
		conn.Close()
		return nil, err
	}

	logger.Debug("Subscribed to path")
	logger.Debug("Starting subscription Goroutine")

	ch := make(chan []byte)
	go func() {
		for m := range s.Message() {
			if m == nil {
				break
			}

			logger := logger.WithField("message", string(m.Message))
			logger.Debug("Received message for subscription")

			ch <- m.Message
		}

		logger.Debug("Closing subscription")
		close(ch)
		s.Close()
		conn.Close()
	}()

	logger.Debug("Subscription Goroutine running")
	return &redisSubscription{
		conn: conn,
		sub:  s,
		c:    ch,
	}, nil
}

func (c *RedisCollector) Close() error {
	return nil
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
