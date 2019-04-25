package redis

import (
	"fmt"
	"path"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
	log "github.com/Sirupsen/logrus"
	"github.com/keimoon/gore"
)

type RedisPublisher struct {
	conn   *gore.Conn
	conf   *driver.Driver
	logger *log.Entry
}

func NewRedisPublisher(cfg *driver.Driver) (plugins.Publisher, error) {
	logger := log.WithFields(log.Fields{
		"driver":    cfg.Type,
		"publisher": cfg.SafeURLString(),
		"url":       cfg.URL,
	})

	c, err := gore.Dial(cfg.URL.Host)
	if err != nil {
		logger.WithError(err).Error("Failed to connect to Redis")
		return nil, err
	}

	if cfg.URL.User != nil {
		if pass, ok := cfg.URL.User.Password(); ok {
			logger.Debug("Authenticating to Redis")
			err := c.Auth(pass)
			if err != nil {
				logger.WithError(err).Error("Failed to authenticate to Redis")
				return nil, err
			}
		}
	}

	return &RedisPublisher{
		conn:   c,
		conf:   cfg,
		logger: logger,
	}, nil
}

func (c *RedisPublisher) Driver() *driver.Driver {
	return c.conf
}

func (p *RedisPublisher) Describe() string {
	return p.conf.SafeURLString()
}

func (p *RedisPublisher) Publish(channel string, data []byte) error {
	logger := p.logger.WithFields(log.Fields{
		"data": data,
	})

	if !p.conn.IsConnected() {
		err := fmt.Errorf("not connected")
		logger.WithError(err).Error("Not connected to publisher")
		return err
	}

	logger.Debug("Publishing message to transport")
	if err := gore.Publish(p.conn, path.Join(p.conf.URL.Path, channel), data); err != nil {
		logger.WithError(err).Warn("Failed to publish data")
		return err
	}

	return nil
}

func (p *RedisPublisher) Close() error {
	return p.conn.Close()
}
