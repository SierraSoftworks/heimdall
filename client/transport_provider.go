package main

import (
	"bytes"
	"encoding/json"

	"github.com/SierraSoftworks/heimdall/models"
	"github.com/SierraSoftworks/heimdall/transports"
	"github.com/rubyist/circuitbreaker"
	log "github.com/Sirupsen/logrus"
)

type TransportProvider struct {
	transports map[string]*clientTransport
}

func NewTransportProvider() *TransportProvider {
	return &TransportProvider{
		transports: map[string]*clientTransport{},
	}
}

func (p *TransportProvider) PublishReport(r *models.Report) {
	logger := log.
		WithField("report", r)

	success := false

	for name, tr := range p.transports {
		logger := log.
			WithField("transport", name).
			WithField("url", tr.URL)

		err := tr.PublishReport(r)
		if err != nil {
			logger.WithError(err).Error("failed to publish report to transport")
			continue
		}

		logger.Debug("published report to transport")
		success = true
	}

	if success {
		logger.Info("report published")
	} else {
		logger.Warning("failed to publish report")
	}
}

func (p *TransportProvider) Configure(conf *Config) error {
	active := map[string]struct{}{}

	for _, tc := range conf.Transports {
		active[tc.Name] = struct{}{}
		logger := log.
			WithField("transport", tc.Name).
			WithField("url", tc.URL)

		if tr, ok := p.transports[tc.Name]; !ok {
			logger.Debug("connecting to transport")
			ct, err := newClientTransport(tc.URL)
			if err != nil {
				logger.WithError(err).Error("failed to connect to transport")
				continue
			}

			logger.Info("connected to transport")
			p.transports[tc.Name] = ct
		} else if tr.URL != tc.URL {
			logger.Debug("reconfiguring transport")
			oldTransport := tr.Transport

			logger.Debug("connecting to new transport URL")
			ntr, err := transports.GetTransport(tc.URL)
			if err != nil {
				logger.WithError(err).Error("failed to connect to transport")
				continue
			}
			tr.Transport = ntr

			logger.Debug("closing old transport connection")
			if err := oldTransport.Close(); err != nil {
				logger.WithError(err).Error("failed to close old transport")
			}
		} else {
			logger.Debug("no changes to transport")
		}
	}

	log.Debug("cleaning out removed transports")
	for name, tr := range p.transports {
		logger := log.
			WithField("transport", name).
			WithField("url", tr.URL)

		if _, ok := active[name]; !ok {
			logger.Debug("removing transport")
			delete(p.transports, name)
			err := tr.Transport.Close()
			if err != nil {
				logger.WithError(err).Error("failed to close old transport")
				continue
			}
		}
	}

	return nil
}

func (p *TransportProvider) Shutdown() {
	for _, tr := range p.transports {
		tr.Transport.Close()
	}
}

type clientTransport struct {
	URL       string
	Transport transports.Transport
	Breaker *circuit.Breaker
}

func newClientTransport(url string) (*clientTransport, error) {
	tr, err := transports.GetTransport(url)
	if err != nil {
		return nil, err
	}

	return &clientTransport{
		URL:       url,
		Transport: tr,
		Breaker: circuit.NewRateBreaker(0.95, 100),
	}, nil
}

func (t *clientTransport) PublishReport(c *models.Report) error {
	return t.Breaker.Call(func() error {
		b := bytes.NewBuffer([]byte{})
		err := json.NewEncoder(b).Encode(c)
		if err != nil {
			return err
		}

		return t.Transport.Publish(transports.CompletedCheckTopic, b.Bytes())
	}, 0)
}
