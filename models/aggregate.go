package models

import "time"

// Aggregate represents a grouped set of checks on various clients,
// allowing customers to create collections of related checks to
// represent the state of specific services and resources with ease.
type Aggregate struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
}

// AggregateDetails provides detailed insight into the state of an
// aggregate
type AggregateDetails struct {
	*Aggregate
	Entries []AggregateEntry `json:"entries"`
}

// AggregateEntry represents the state of a specific check which is
// a member of this aggregate.
type AggregateEntry struct {
	CheckName  string    `json:"check"`
	ClientName string    `json:"client"`
	Status     Status    `json:"status"`
	Executed   time.Time `json:"executed"`
}
