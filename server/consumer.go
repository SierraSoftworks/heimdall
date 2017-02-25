package main

import (
	"sync"

	"github.com/SierraSoftworks/heimdall/models"
	"github.com/SierraSoftworks/heimdall/stores"
	log "github.com/Sirupsen/logrus"
)

type Consumer struct {
	transports   []*ServerTransport
	c            chan *models.Report
	wg           sync.WaitGroup
	noTransports sync.Mutex
}

func NewConsumer() *Consumer {
	c := &Consumer{
		transports: []*ServerTransport{},
		c:          make(chan *models.Report),
	}
	c.noTransports.Lock()

	go func() {
		c.noTransports.Lock()
		c.wg.Wait()
		close(c.c)
	}()

	return c
}

func (c *Consumer) AddTransport(t *ServerTransport) {
	for _, x := range c.transports {
		if x == t {
			return
		}
	}

	c.transports = append(c.transports, t)
	go func() {
		c.wg.Add(1)
		for r := range t.Reports() {
			c.c <- &r
		}
		c.wg.Done()
	}()

	c.noTransports.Unlock()
}

func (c *Consumer) Run() {
	for r := range c.c {
		err := stores.GetStore().AddReport(r)
		if err != nil {
			log.
				WithError(err).
				WithField("report", r).
				Error("Failed to persist report to store")
		}
	}
}

func (c *Consumer) Wait() {
	c.wg.Wait()
}
