package file

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
	log "github.com/Sirupsen/logrus"
)

type FileCollector struct {
	conf   *driver.Driver
	logger *log.Entry
}

func NewFileCollector(cfg *driver.Driver) (plugins.Collector, error) {
	logger := log.WithFields(log.Fields{
		"driver":    cfg.Type,
		"collector": cfg.SafeURLString(),
		"url":       cfg.URL,
	})

	return &FileCollector{
		conf:   cfg,
		logger: logger,
	}, nil
}

func (c *FileCollector) Driver() *driver.Driver {
	return c.conf
}

func (c *FileCollector) Describe() string {
	return fmt.Sprintf("file://%s", c.fileName())
}

func (c *FileCollector) Subscribe(channel string) (plugins.Subscription, error) {
	logger := c.logger.WithFields(log.Fields{
		"file": c.fileName(),
	})

	logger.Debug("Creating new subscriber")

	logger.Debug("Opening file for subscription")
	f, err := os.OpenFile(c.fileName(), os.O_CREATE|os.O_RDONLY|os.O_SYNC, 0664)
	if err != nil {
		logger.WithError(err).Error("Failed to open file for subscription")
		return nil, err
	}

	ch := make(chan []byte)
	go func() {
		r := bufio.NewScanner(NewTailReader(f, 10*time.Millisecond))

		for r.Scan() {
			logger := logger.WithField("entry", r.Text())
			logger.Debug("Read entry from file")

			entry := FileFormat{}
			buf := bytes.NewBuffer(r.Bytes())
			if err := json.NewDecoder(buf).Decode(&entry); err != nil {
				logger.WithError(err).Error("Failed to decode file entry")
				continue
			}

			if entry.Channel != channel {
				continue
			}

			ch <- entry.Data
		}

		if err := r.Err(); err != nil {
			logger.WithError(err).Warn("File subscription exited with an error")
		}

		logger.Debug("Closing subscriber")

		close(ch)
	}()

	return &fileSubscription{
		f: f,
		c: ch,
	}, nil
}

func (c *FileCollector) Close() error {
	return nil
}

func (c *FileCollector) fileName() string {
	return filepath.Join(append([]string{c.conf.URL.Host}, strings.Split(c.conf.URL.Path, "/")...)...)
}

type fileSubscription struct {
	f *os.File
	c chan []byte
}

func (s *fileSubscription) Channel() <-chan []byte {
	return s.c
}

func (s *fileSubscription) Close() error {
	return s.f.Close()
}
