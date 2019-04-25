package models

import (
	"reflect"

	"github.com/SierraSoftworks/heimdall/pkg/duration"
)

// Check describes a check that should be executed on a client
// to determine the state of a service or resource within that
// client's domain.
// Checks are uniquely identified by their name however they may
// also be placed in aggregates for the purpose of analysis.
type Check struct {
	ID          string            `json:"id"`
	Command     string            `json:"command"`
	Interval    duration.Duration `json:"interval"`
	Timeout     duration.Duration `json:"timeout"`
	Collections []string          `json:"collections"`
}

func (c *Check) Equals(o *Check) bool {
	if c.ID != "" || o.ID != "" {
		return c.ID == o.ID
	}

	return c.Command == o.Command &&
		c.Interval.Equals(o.Interval) &&
		c.Timeout.Equals(o.Timeout) &&
		reflect.DeepEqual(c.Collections, o.Collections)
}
