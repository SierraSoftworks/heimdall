package nats

import (
	"fmt"
	"path"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
	log "github.com/Sirupsen/logrus"
	nats "github.com/nats-io/go-nats"
)

type NATSPublisher struct {
	conn   *nats.Conn
	conf   *driver.Driver
	logger *log.Entry
}

func NewNATSPublisher(cfg *driver.Driver) (plugins.Publisher, error) {
	logger := log.WithFields(log.Fields{
		"driver":    cfg.Type,
		"collector": cfg.SafeURLString(),
		"url":       cfg.URL,
	})

	c, err := nats.Connect(cfg.URL.String())
	if err != nil {
		logger.WithError(err).Error("Failed to connect to NATS server")
		return nil, err
	}

	return &NATSPublisher{
		conn:   c,
		conf:   cfg,
		logger: logger,
	}, nil
}

func (c *NATSPublisher) Driver() *driver.Driver {
	return c.conf
}

func (p *NATSPublisher) Describe() string {
	return p.conf.SafeURLString()
}

func (p *NATSPublisher) Publish(channel string, data []byte) error {
	logger := p.logger.WithField("data", data)

	if p.conn == nil || p.conn.IsClosed() {
		err := fmt.Errorf("not connected")
		logger.WithError(err).Error("Not connected to the NATS server")
		return err
	}

	logger.Debug("Publishing data")
	if err := p.conn.Publish(path.Join(p.conf.URL.Path, channel), data); err != nil {
		logger.WithError(err).Warn("Failed to publish data")
		return err
	}

	return nil
}

func (p *NATSPublisher) Close() error {
	p.conn.Close()
	return nil
}
