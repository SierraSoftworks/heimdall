package main

import (
	"os"

	"github.com/SierraSoftworks/heimdall/models"
	"github.com/SierraSoftworks/heimdall/utils"
	log "github.com/Sirupsen/logrus"
)

type Config struct {
	Client     *models.Client           `json:"client"`
	Checks     []models.Check           `json:"checks"`
	Transports []models.TransportConfig `json:"transports"`
}

func ReadConfig(path string) (*Config, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	c := &Config{
		Client: &models.Client{
			Name: hostname,
			Tags: map[string]string{},
		},
		Transports: []models.TransportConfig{},
		Checks:     []models.Check{},
	}

	cfiles, err := utils.FindConfig(path)
	if err != nil {
		log.
			WithField("path", path).
			WithError(err).
			Warn("Failed to find config files")
		return c, err
	}

	for _, cfile := range cfiles {
		dc := Config{}
		err := utils.LoadConfig(cfile, &dc)
		if err != nil {
			log.
				WithField("file", cfile).
				WithError(err).
				Warn("Failed to load config file")
			continue
		}

		c.Transports = append(c.Transports, dc.Transports...)
		c.Checks = append(c.Checks, dc.Checks...)
	}

	return c, nil
}
