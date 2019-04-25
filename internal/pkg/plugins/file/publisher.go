package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/SierraSoftworks/heimdall/pkg/models"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
	log "github.com/Sirupsen/logrus"
)

type FilePublisher struct {
	f      *os.File
	conf   *driver.Driver
	logger *log.Entry
}

func NewFilePublisher(cfg *driver.Driver) (plugins.Publisher, error) {
	p := &FilePublisher{
		conf: cfg,
		logger: log.WithFields(log.Fields{
			"driver":    cfg.Type,
			"publisher": cfg.SafeURLString(),
			"url":       cfg.URL,
		}),
	}

	fileName := p.fileName()
	p.logger = p.logger.WithField("file", fileName)

	p.logger.Debug("Opening file for publisher")
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0664)
	if err != nil {
		p.logger.WithError(err).Error("Failed to open file for publisher")
		return nil, err
	}

	p.f = f

	return p, nil
}

func (c *FilePublisher) Driver() *driver.Driver {
	return c.conf
}

func (p *FilePublisher) Describe() string {
	return p.conf.SafeURLString()
}

func (p *FilePublisher) Publish(channel string, data []byte) error {
	logger := p.logger.WithField("data", string(data))

	if p.f == nil {
		err := fmt.Errorf("not connected")
		logger.WithError(err).Error("File publisher has been clsoed")
		return err
	}

	logger.Debug("Publishing report")
	defer p.f.Sync()

	entry := FileFormat{
		Channel: channel,
		Data:    data,
	}

	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(&entry); err != nil {
		logger.WithError(err).Error("Failed to encode entry for file")
		return err
	}

	if _, err := p.f.Write(buf.Bytes()); err != nil {
		logger.WithError(err).Error("Failed to publish report")
		return err
	}

	return nil
}

func (p *FilePublisher) Close() error {
	if p.f == nil {
		return nil
	}

	p.logger.Debug("Closing publisher")

	f := p.f
	p.f = nil
	return f.Close()
}

func (p *FilePublisher) fileName() string {
	return filepath.Join(append([]string{p.conf.URL.Host}, strings.Split(p.conf.URL.Path, "/")...)...)
}

type fileEntry struct {
	Topic string         `json:"topic"`
	Data  *models.Report `json:"data"`
}
