package main

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/SierraSoftworks/heimdall/models"
	"github.com/SierraSoftworks/scheduler"
)

type Checker struct {
	Config *Config

	runner            *Runner
	scheduledChecks   []*scheduler.ActiveTask
	keepaliveSchedule *scheduler.ActiveTask

	l       sync.Mutex
	reports chan *models.Report
}

func NewChecker() *Checker {
	c := &Checker{
		runner:          NewDefaultRunner(),
		scheduledChecks: []*scheduler.ActiveTask{},
		reports:         make(chan *models.Report),
	}

	return c
}

func (c *Checker) Reports() <-chan *models.Report {
	return c.reports
}

func (c *Checker) Configure(conf *Config) error {
	c.l.Lock()
	defer c.l.Unlock()

	c.Config = conf

	if c.keepaliveSchedule == nil {
		c.keepaliveSchedule = scheduler.Do(func(t time.Time) error {
			c.reports <- &models.Report{
				Check: &models.Check{
					Name:     "Keep Alive",
					Command:  "",
					Interval: 30 * time.Second,
					Timeout:  0,
				},
				Client: c.Config.Client,
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
	}

	for _, t := range c.scheduledChecks {
		t.Cancel()
	}

	c.scheduledChecks = []*scheduler.ActiveTask{}
	for _, check := range c.Config.Checks {
		func(check models.Check) {
			schedule := scheduler.Do(func(t time.Time) error {
				c.reports <- &models.Report{
					Check:     &check,
					Client:    c.Config.Client,
					Execution: c.runner.ExecuteCheck(&check),
				}

				return nil
			}).Every(check.Interval).Schedule()

			c.scheduledChecks = append(c.scheduledChecks, schedule)
		}(check)
	}

	return nil
}

func (c *Checker) Shutdown() {
	c.l.Lock()
	defer c.l.Unlock()

	for _, t := range c.scheduledChecks {
		t.Cancel()
	}
	c.scheduledChecks = []*scheduler.ActiveTask{}

	c.keepaliveSchedule.Cancel()

	close(c.reports)
}

func (c *Checker) Describe() string {
	b := bytes.NewBuffer([]byte{})
	b.WriteString("\n")
	b.WriteString("Checks:\n")
	for _, c := range c.Config.Checks {
		b.WriteString(fmt.Sprintf("  - %s (every %s)\n", c.Name, c.Interval))
	}

	return b.String()
}
