package config

import (
	log "github.com/Sirupsen/logrus"
)

// ReadConfig reads Heimdall configuration from a file
// or directory of files.
func ReadConfig(path string) (*Config, error) {
	c, err := NewConfig()
	if err != nil {
		log.WithField("path", path).
			WithError(err).
			Error("Failed to generate default config")
		return nil, err
	}

	cfiles, err := FindConfig(path)
	if err != nil {
		log.
			WithField("path", path).
			WithError(err).
			Warn("Failed to find config files")
		return c, err
	}

	for _, cfile := range cfiles {
		dc := Config{}
		err := LoadConfig(cfile, &dc)
		if err != nil {
			log.
				WithField("file", cfile).
				WithError(err).
				Warn("Failed to load config file")
			continue
		}

		c.Update(&dc)
	}

	return c, nil
}
