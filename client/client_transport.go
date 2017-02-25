package main

import (
	"bytes"
	"encoding/json"

	"github.com/SierraSoftworks/heimdall/models"
	"github.com/SierraSoftworks/heimdall/transports"
)

type ClientTransport struct {
	Transport transports.Transport
}

func NewClientTransport(t transports.Transport) *ClientTransport {
	return &ClientTransport{
		Transport: t,
	}
}

func (t *ClientTransport) PublishCheck(c *models.Report) error {
	b := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(b).Encode(c)
	if err != nil {
		return err
	}

	return t.Transport.Publish(transports.CompletedCheckTopic, b.Bytes())
}
