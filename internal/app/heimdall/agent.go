package heimdall

import (
	"github.com/SierraSoftworks/heimdall/internal/pkg/api"
	"github.com/SierraSoftworks/heimdall/internal/pkg/config"
	"github.com/SierraSoftworks/heimdall/internal/pkg/hub"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
	log "github.com/Sirupsen/logrus"
)

type Agent struct {
	Config *config.Config

	hub hub.Hub

	store plugins.Store
	api   *api.API

	colls []plugins.Collector
	pubs  []plugins.Publisher

	logger *log.Entry
}

func NewAgent(c *config.Config) (*Agent, error) {
	store, err := plugins.GetStore(c.Store)
	if err != nil {
		return nil, err
	}

	api, err := api.NewAPI(store)
	if err != nil {
		return nil, err
	}

	h := hub.NewMemoryHub()

	return &Agent{
		Config: c,

		hub: h,

		store: store,
		api:   api,

		colls:  []plugins.Collector{},
		pubs:   []plugins.Publisher{},
		logger: log.WithFields(log.Fields{}),
	}, nil
}

func (a *Agent) Run() error {
	a.logger.Debug("Loading collectors")
	if err := a.loadCollectors(); err != nil {
		a.logger.WithError(err).Error("Failed to load collectors")
		return err
	}

	a.logger.Debug("Loading publishers")
	if err := a.loadPublishers(); err != nil {
		a.logger.WithError(err).Error("Failed to load publishers")
		return err
	}

	return a.api.ListenAndServe(a.Config.API.Listen)
}

func (a *Agent) loadCollectors() error {
	for _, driver := range a.Config.Collectors {
		logger := a.logger.WithField("driver", driver.Type)
		coll, err := plugins.GetCollector(&driver)
		if err != nil {
			logger.WithError(err).Error("Failed to create collector")
			return err
		}

		a.colls = append(a.colls, coll)
	}

	return nil
}

func (a *Agent) loadPublishers() error {
	for _, driver := range a.Config.Publishers {
		logger := a.logger.WithField("driver", driver.Type)
		pub, err := plugins.GetPublisher(&driver)
		if err != nil {
			logger.WithError(err).Error("Failed to create publisher")
			return err
		}

		a.pubs = append(a.pubs, pub)
	}

	return nil
}
