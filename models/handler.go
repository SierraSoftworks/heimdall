package models

// Handler describes a task which is executed to handle the
// results of check executions for the purpose of integrating
// with external systems.
// Specifically, handlers enable push mode data flows rather
// instead of relying on polling requests. They are useful for
// propagating state changes to external systems.
type Handler struct {
	Name   string                 `json:"name"`
	Filter map[string]interface{} `json:"filter"`

	Driver string `json:"driver"`
	URL    string `json:"url"`
}
