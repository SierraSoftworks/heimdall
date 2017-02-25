package main

import (
	"fmt"
	"sync"
	"time"

	"bytes"

	"github.com/SierraSoftworks/heimdall/models"
	"github.com/SierraSoftworks/scheduler"
	log "github.com/Sirupsen/logrus"
)

type Client struct {
	Checks          []models.Check
	scheduledChecks []*scheduler.ActiveTask
	Runner          *Runner
	Transports      []*ClientTransport

	l                 sync.Mutex
	reports           chan *models.Report
	keepaliveSchedule *scheduler.ActiveTask
}

func NewClient() *Client {
	c := &Client{
		Checks:     []models.Check{},
		Runner:     NewDefaultRunner(),
		Transports: []*ClientTransport{},

		scheduledChecks: []*scheduler.ActiveTask{},
		reports:         make(chan *models.Report),
	}

	c.keepaliveSchedule = scheduler.Do(func(t time.Time) error {
		c.reports <- &models.Report{
			Check: &models.Check{
				Name:     "Keep Alive",
				Command:  "",
				Interval: 30 * time.Second,
				Timeout:  0,
			},
			Client: config.Client,
			Execution: &models.Execution{
				Scheduled: t,
				Executed:  t,
				Duration:  0,
				Status:    models.StatusOkay,
				Output:    c.Describe(),
			},
		}

		return nil
	}).Every(30 * time.Second).Schedule()

	go func() {
		for check := range c.reports {
			for _, tr := range c.Transports {
				err := tr.PublishCheck(check)
				if err != nil {
					log.
						WithField("check", check).
						WithField("transport", tr.Transport.Describe()).
						WithError(err).
						Error("Failed to publish check result")
				}
			}
		}
	}()

	return c
}

func (c *Client) Reschedule() {
	c.l.Lock()
	defer c.l.Unlock()

	for _, t := range c.scheduledChecks {
		t.Cancel()
	}

	c.scheduledChecks = []*scheduler.ActiveTask{}
	for _, check := range c.Checks {
		func(check models.Check) {
			schedule := scheduler.Do(func(t time.Time) error {
				c.reports <- &models.Report{
					Check:     &check,
					Client:    config.Client,
					Execution: c.Runner.Execute(&check),
				}

				return nil
			}).Every(check.Interval).Schedule()

			c.scheduledChecks = append(c.scheduledChecks, schedule)
		}(check)
	}
}

func (c *Client) Shutdown() {
	c.l.Lock()
	defer c.l.Unlock()

	for _, t := range c.scheduledChecks {
		t.Cancel()
	}
	c.scheduledChecks = []*scheduler.ActiveTask{}

	c.keepaliveSchedule.Cancel()
}

func (c *Client) Describe() string {
	b := bytes.NewBuffer([]byte{})
	b.WriteString("Transports:\n")
	for _, t := range c.Transports {
		b.WriteString(fmt.Sprintf("  - %s\n", t.Transport.Describe()))
	}

	b.WriteString("\n")
	b.WriteString("Checks:\n")
	for _, c := range c.Checks {
		b.WriteString(fmt.Sprintf("  - %s (every %s)\n", c.Name, c.Interval))
	}

	return b.String()
}
