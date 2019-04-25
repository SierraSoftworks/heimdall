package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"fmt"

	log "github.com/Sirupsen/logrus"
)

func FindConfig(path string) ([]string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		log.WithField("path", path).WithError(err).Error("Failed to stat config")
		return nil, err
	}

	if !fi.IsDir() {
		return []string{path}, nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.WithField("path", path).WithError(err).Error("Failed to enumerate config directory")
		return nil, err
	}

	results := []string{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		ext := filepath.Ext(f.Name())

		if ext == ".yaml" || ext == ".yml" {
			results = append(results, filepath.Join(path, f.Name()))
		}
	}

	return results, nil
}

func LoadConfig(file string, config interface{}) error {
	log.WithField("file", file).Debug("Loading configuration")

	fi, err := os.Stat(file)
	if err != nil {
		log.WithField("file", file).WithError(err).Warn("Failed to stat config file/dir")
		return err
	}

	if fi.IsDir() {
		log.WithField("file", file).WithError(err).Error("Cannot load a directory as a configuration file")
		return fmt.Errorf("directories not supported")
	}

	log.WithField("file", file).Debug("Reading config file")

	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.WithField("file", file).WithError(err).Warn("Failed to read config file")
		return err
	}

	log.WithField("file", file).Debug("Parsing config file")
	err = yaml.Unmarshal(f, config)
	if err != nil {
		log.WithField("file", file).WithError(err).Warn("Failed to parse config file")
		return err
	}

	log.WithField("config", config).Debug("Config ready")

	return nil
}
