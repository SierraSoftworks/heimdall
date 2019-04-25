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
)

func init() {
	RegisterManager(&StoreManager{
		stores: []plugins.Store{},
	})
}

type StoreManager struct {
	hub    hub.Hub
	stores []plugins.Store

	l sync.RWMutex
}

func (m *StoreManager) Setup(h hub.Hub) error {
	m.hub = h
	m.hub.Subscribe(NewConfigSubscriber(m.ConfigUpdated))
	m.hub.Subscribe(NewReportSubscriber(m.ReportReceived))

	return nil
}

func (m *StoreManager) Shutdown() {
	m.l.Lock()
	defer m.l.Unlock()

	for _, store := range m.stores {
		logger := log.WithField("store", store.Driver().Describe())
		logger.Debug("Shutting down store")
		if err := store.Close(); err != nil {
			logger.WithError(err).Warning("Failed to shutdown store cleanly")
		}
	}
}

func (m *StoreManager) ConfigUpdated(conf *config.Config) {
	m.l.Lock()
	defer m.l.Unlock()

	skip := map[int]bool{}
	remove := map[int]bool{}

	for i, store := range m.stores {
		remove[i] = true

	inner:
		for j, newStore := range conf.Stores {
			if store.Driver().Equals(&newStore) {
				skip[j] = true
				remove[i] = false
				break inner
			}
		}
	}

	newStores := []plugins.Store{}
	for i, store := range m.stores {
		logger := log.WithField("store", store.Driver().Describe())

		if remove[i] {
			logger.Info("Removing existing store")
			if err := store.Close(); err != nil {
				logger.WithError(err).Warning("Failed to shutdown store cleanly")
			}
			continue
		}

		logger.Debug("Keeping existing store")
		newStores = append(newStores, store)
	}

	for j, newStore := range conf.Stores {
		logger := log.WithField("store", newStore.Describe())

		if !skip[j] {
			logger.Debug("Adding new store")
			store, err := plugins.GetStore(&newStore)
			if err != nil {
				logger.WithError(err).Error("Failed to add store")
			} else {
				logger.Info("Added new store")
				newStores = append(newStores, store)
			}
		}
	}

	m.stores = newStores
}

func (m *StoreManager) ReportReceived(report *models.Report) {
	if report.Source != nil && report.Source.ID == m.source.ID {
		// Don't process reports that we sent already
		return
	}

	m.l.RLock()
	defer m.l.RUnlock()

	logger := log.WithField("report", report)

	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(report); err != nil {
		logger.
			WithError(err).
			Warning("Failed to serialize report object for publication")
	}

	rm := report.ToMap()
	for _, s := range m.stores {
		if s.Driver().Matches(rm) {
			logger := logger.WithField("store", s.Driver().Describe())

			logger.Debug("Storing report")
			if err := s.Store(report); err != nil {
				logger.WithError(err).Warning("Failed to store report")
			}
		}
	}
}
