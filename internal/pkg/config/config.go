package config

import (
	"os"
	"time"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/SierraSoftworks/heimdall/pkg/models"
)

// Config describes the Heimdall config file
// structure.
type Config struct {
	API *APIConfig `json:"api"`

	Source    *models.Source   `json:"source"`
	Keepalive *KeepaliveConfig `json:"keepalive"`

	Checks []models.Check `json:"checks"`

	Stores     []driver.Driver `json:"stores"`
	Collectors []driver.Driver `json:"collectors"`
	Publishers []driver.Driver `json:"publishers"`
}

// NewConfig creates a default configuration
// file which should be populated using the
// config.Update function.
func NewConfig() (*Config, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &Config{
		API: &APIConfig{
			Listen: ":80",
		},

		Keepalive: &KeepaliveConfig{
			Interval: time.Minute,
		},

		Source: &models.Source{
			Name: hostname,
			Tags: map[string]string{},
		},

		Checks: []models.Check{},

		Stores:     []driver.Driver{},
		Collectors: []driver.Driver{},
		Publishers: []driver.Driver{},
	}, nil
}

// Update is used to apply a differential config
// to this configuration object, updating any
// new fields.
func (c *Config) Update(dc *Config) {
	if dc.API != nil {
		c.API.Update(dc.API)
	}

	if dc.Keepalive != nil {
		c.Keepalive.Update(dc.Keepalive)
	}

	if dc.Source != nil {
		if dc.Source.Name != "" {
			c.Source.Name = dc.Source.Name
		}

		if dc.Source.Tags != nil {
			for k, v := range dc.Source.Tags {
				c.Source.Tags[k] = v
			}
		}
	}

	c.Checks = append(c.Checks, dc.Checks...)

	c.Stores = append(c.Stores, dc.Stores...)
	c.Collectors = append(c.Collectors, dc.Collectors...)
	c.Publishers = append(c.Publishers, dc.Publishers...)
}
