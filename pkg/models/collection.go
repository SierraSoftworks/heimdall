package models

// Collection represents a grouped set of checks on various clients,
// allowing customers to create collections of related checks to
// represent the state of specific services and resources with ease.
type Collection struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
}

// CollectionDetails provides detailed insight into the state of an
// aggregate
type CollectionDetails struct {
	*Collection
	Entries []*Execution `json:"entries"`
}
