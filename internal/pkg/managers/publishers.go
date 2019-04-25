package managers

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/SierraSoftworks/heimdall/internal/pkg/config"
	"github.com/SierraSoftworks/heimdall/internal/pkg/hub"
	"github.com/SierraSoftworks/heimdall/internal/pkg/lamport"
	"github.com/SierraSoftworks/heimdall/pkg/models"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"

	log "github.com/Sirupsen/logrus"
)

func init() {
	RegisterManager(&PublishersManager{
		publishers: []plugins.Publisher{},
		clock:      lamport.NewClock(0),
	})
}

type PublishersManager struct {
	hub hub.Hub

	publishers []plugins.Publisher
	clock      *lamport.Clock
	source     *models.Source

	l sync.RWMutex
}

func (m *PublishersManager) Setup(h hub.Hub) error {
	m.hub = h

	m.hub.Subscribe(NewConfigSubscriber(m.ConfigUpdated))
	m.hub.Subscribe(NewReportSubscriber(m.ReportReceived))

	return nil
}

func (m *PublishersManager) Shutdown() {
	m.l.Lock()
	defer m.l.Unlock()

	for _, pub := range m.publishers {
		logger := log.WithField("publisher", pub.Describe())
		logger.Debug("Shutting down publisher")
		if err := pub.Close(); err != nil {
			logger.WithError(err).Warning("Failed to shutdown publisher cleanly")
		}
	}

	m.publishers = []plugins.Publisher{}
}

func (m *PublishersManager) ConfigUpdated(conf *config.Config) {
	m.l.Lock()
	defer m.l.Unlock()

	skip := map[int]bool{}
	remove := map[int]bool{}

	for i, pub := range m.publishers {
		remove[i] = true

	inner:
		for j, newPub := range conf.Publishers {
			if pub.Driver().Equals(&newPub) {
				skip[j] = true
				remove[i] = false
				break inner
			}
		}
	}

	newPubs := []plugins.Publisher{}
	for i, pub := range m.publishers {
		logger := log.WithField("publisher", pub.Describe())

		if remove[i] {
			logger.Info("Removing existing publisher")
			if err := pub.Close(); err != nil {
				logger.WithError(err).Warning("Failed to shutdown publisher cleanly")
			}
			continue
		}

		logger.Debug("Keeping existing publisher")
		newPubs = append(newPubs, pub)
	}

	for j, newPub := range conf.Publishers {
		logger := log.WithField("publisher", newPub.Describe())

		if !skip[j] {
			logger.Debug("Adding new publisher")
			pub, err := plugins.GetPublisher(&newPub)
			if err != nil {
				logger.WithError(err).Error("Failed to add publisher")
				continue
			} else {
				logger.Info("Added new publisher")
				newPubs = append(newPubs, pub)
			}
		}
	}

	m.publishers = newPubs
	m.source = conf.Source
}

func (m *PublishersManager) ReportReceived(report *models.Report) {
	if report.Source != nil && report.Source.ID == m.source.ID {
		// Don't re-send reports that we generated
		return
	}

	m.l.RLock()
	defer m.l.RUnlock()

	// Mark the report as originating from our current agent
	report.Source = m.source
	// Update the report's sequence number (lamport clock)
	report.SequenceNumber = m.clock.Update(report.SequenceNumber)

	logger := log.WithField("report", report)

	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(report); err != nil {
		logger.
			WithError(err).
			Warning("Failed to serialize report object for publication")
	}

	rm := report.ToMap()
	for _, p := range m.publishers {
		if p.Driver().Matches(rm) {
			logger := logger.WithField("publisher", p.Describe())

			logger.Debug("Publishing report")
			if err := p.Publish(plugins.ReportsChannel, buf.Bytes()); err != nil {
				logger.WithError(err).Warning("Failed to publish report")
			}
		}
	}
}
