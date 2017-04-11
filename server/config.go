package main

import (
	"github.com/SierraSoftworks/heimdall/models"
	"github.com/SierraSoftworks/heimdall/utils"
	log "github.com/Sirupsen/logrus"
)

type Config struct {
	Listen     string             `json:"listen"`
	Transports []models.Transport `json:"transports"`
	Handlers   []models.Handler   `json:"handlers"`
}

func ReadConfig(path string) (*Config, error) {
	c := &Config{
		Listen:     ":80",
		Transports: []models.TransportConfig{},
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

		if dc.Listen != "" {
			c.Listen = dc.Listen
		}

		c.Transports = append(c.Transports, dc.Transports...)
	}

	return c, nil
}
