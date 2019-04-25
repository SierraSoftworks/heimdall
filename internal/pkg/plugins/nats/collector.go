package nats

import (
	"fmt"
	"path"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
	log "github.com/Sirupsen/logrus"
	nats "github.com/nats-io/go-nats"
)

type NATSCollector struct {
	conf   *driver.Driver
	conn   *nats.Conn
	logger *log.Entry
}

func NewNATSCollector(cfg *driver.Driver) (plugins.Collector, error) {
	logger := log.WithFields(log.Fields{
		"driver":    cfg.Type,
		"collector": cfg.SafeURLString(),
		"url":       cfg.URL,
	})

	logger.Debug("Connecting to NATS server")
	c, err := nats.Connect(cfg.URL.String())
	if err != nil {
		logger.WithError(err).Error("Failed to connect to NATS server")
		return nil, err
	}

	return &NATSCollector{
		conn:   c,
		conf:   cfg,
		logger: logger,
	}, nil
}

func (c *NATSCollector) Driver() *driver.Driver {
	return c.conf
}

func (c *NATSCollector) Describe() string {
	return c.conf.SafeURLString()
}

func (c *NATSCollector) Subscribe(channel string) (plugins.Subscription, error) {
	queueGroup := c.conf.GetOption("queue_group", "heimdall_servers")
	logger := c.logger.WithField("queue_group", queueGroup)

	if c.conn == nil || c.conn.IsClosed() {
		err := fmt.Errorf("not connected")
		logger.WithError(err).Error("Not connected to NATS")
		return nil, err
	}

	ch := make(chan []byte)
	s, err := c.conn.QueueSubscribe(path.Join(c.conf.URL.Path, channel), queueGroup, func(m *nats.Msg) {
		logger := logger.WithField("message", string(m.Data))
		logger.Debug("Received message from subscription")

		ch <- m.Data
	})

	if err != nil {
		logger.WithError(err).Error("Failed to create subscription")
		return nil, err
	}

	return &natsSubscription{
		c:   ch,
		sub: s,
	}, nil
}

func (c *NATSCollector) Close() error {
	c.conn.Close()
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
