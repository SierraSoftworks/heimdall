package managers

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/SierraSoftworks/heimdall/internal/pkg/config"
	"github.com/SierraSoftworks/heimdall/internal/pkg/hub"
	"github.com/SierraSoftworks/heimdall/pkg/models"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

func init() {
	RegisterManager(&CollectorsManager{})
}

type CollectorsManager struct {
	hub hub.Hub

	collectors []plugins.Collector

	l sync.Mutex
}

func (m *CollectorsManager) Setup(h hub.Hub) error {
	m.hub = h

	m.hub.Subscribe(NewConfigSubscriber(m.ConfigUpdated))

	return nil
}

func (m *CollectorsManager) Shutdown() {
	m.l.Lock()
	defer m.l.Unlock()

	for _, col := range m.collectors {
		logger := log.WithField("collector", col.Describe())
		logger.Debug("Shutting down collector")
		if err := col.Close(); err != nil {
			logger.WithError(err).Warning("Failed to shutdown collector cleanly")
		}
	}
}

func (m *CollectorsManager) ConfigUpdated(conf *config.Config) {
	m.l.Lock()
	defer m.l.Unlock()

	skip := map[int]bool{}
	remove := map[int]bool{}

	for i, col := range m.collectors {
		remove[i] = true

	inner:
		for j, newCol := range conf.Collectors {
			if col.Driver().Equals(&newCol) {
				skip[j] = true
				remove[i] = false
				break inner
			}
		}
	}

	newCols := []plugins.Collector{}
	for i, col := range m.collectors {
		logger := log.WithField("collector", col.Describe())

		if remove[i] {
			logger.Info("Removing existing collector")
			if err := col.Close(); err != nil {
				logger.WithError(err).Warning("Failed to shutdown collector cleanly")
			}
			continue
		}

		logger.Debug("Keeping existing collector")
		newCols = append(newCols, col)
	}

	for j, newCol := range conf.Collectors {
		logger := log.WithField("collector", newCol.Describe())

		if !skip[j] {
			logger.Debug("Adding new collector")
			col, err := plugins.GetCollector(&newCol)
			if err != nil {
				logger.WithError(err).Error("Failed to add collector")
			} else if err := m.setupCollector(col); err != nil {
				logger.WithError(err).Error("Failed to configure collector")
				col.Close()
			} else {
				logger.Info("Added new collector")
				newCols = append(newCols, col)
			}
		}
	}

	m.collectors = newCols
}

func (m *CollectorsManager) setupCollector(col plugins.Collector) error {
	logger := log.WithField("collector", col.Describe())

	logger.WithField("channel", plugins.ReportsChannel).Debug("Subscribing to reports channel")
	reports, err := col.Subscribe(plugins.ReportsChannel)
	if err != nil {
		return errors.Wrap(err, "collectors manager: failed to subscribe to reports channel")
	}

	go func() {
		for reportData := range reports.Channel() {
			buf := bytes.NewBuffer(reportData)
			report := models.Report{}
			if err := json.NewDecoder(buf).Decode(&report); err != nil {
				logger.WithError(err).Warning("Failed to deserialize report")
				continue
			}

			m.hub.Notify(&report)
		}
	}()

	return nil
}
