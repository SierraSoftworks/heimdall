package models

// Transport defines the configuration used by
// a specific transport instance. These transports are
// used by clients to submit event entries and by servers
// to receive those events.
type Transport struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
