package main

import (
	"github.com/SierraSoftworks/heimdall/utils"
)

type Client struct {
	Config *ConfigProvider

	checker   *Checker
	transport *TransportProvider

	cleanup *utils.ShutdownManager
}

func NewClient(configPath string) (*Client, error) {
	conf, err := NewConfigProvider(configPath)
	if err != nil {
		return nil, err
	}

	cl := &Client{
		Config: conf,

		checker:   NewChecker(),
		transport: NewTransportProvider(),
		cleanup:   utils.NewShutdownManager(),
	}

	conf.AddTarget(cl.checker)
	cl.cleanup.AddTarget(cl.checker)

	conf.AddTarget(cl.transport)
	cl.cleanup.AddTarget(cl.transport)

	go func() {
		for report := range cl.checker.Reports() {
			cl.transport.PublishReport(report)
		}
	}()

	if err := conf.Reload(); err != nil {
		cl.Shutdown()
		return nil, err
	}

	return cl, nil
}

func (c *Client) Shutdown() {
	c.cleanup.Shutdown()
}
