package models

import (
	"time"
)

// Check describes a check that should be executed on a client
// to determine the state of a service or resource within that
// client's domain.
// Checks are uniquely identified by their name however they may
// also be placed in aggregates for the purpose of analysis.
type Check struct {
	Name       string        `json:"name"`
	Command    string        `json:"command"`
	Interval   time.Duration `json:"interval"`
	Timeout    time.Duration `json:"timeout"`
	Aggregates []string      `json:"aggregates"`
}

// CheckDetails represents a check and a quick summary of its current
// state to enable easy rendering on dashboards and other interfaces.
type CheckDetails struct {
	*Check
	Status   Status    `json:"status"`
	Executed time.Time `json:"executed"`
}
