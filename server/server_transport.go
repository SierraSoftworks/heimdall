package main

import (
	"bytes"
	"encoding/json"

	"github.com/SierraSoftworks/heimdall/models"
	"github.com/SierraSoftworks/heimdall/transports"
	log "github.com/Sirupsen/logrus"
)

type ServerTransport struct {
	Transport transports.Transport

	sub transports.Subscription
	c   chan models.Report
}

func NewServerTransport(t transports.Transport) (*ServerTransport, error) {
	st := &ServerTransport{
		Transport: t,
		c:         make(chan models.Report),
	}

	sub, err := st.Transport.Subscribe(transports.CompletedCheckTopic)
	if err != nil {
		return nil, err
	}

	st.sub = sub

	go func() {
		for d := range sub.Channel() {
			var r models.Report
			if err := json.NewDecoder(bytes.NewBuffer(d)).Decode(&r); err != nil {
				log.
					WithError(err).
					WithField("report", string(d)).
					WithField("transport", t.Describe()).
					Error("Failed to parse report")
			} else {
				st.c <- r
			}
		}

		close(st.c)
	}()

	return st, nil
}

func (t *ServerTransport) Reports() <-chan models.Report {
	return t.c
}
