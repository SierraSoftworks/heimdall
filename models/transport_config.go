package models

// TransportConfig defines the configuration used by
// a specific transport instance. These transports are
// used by clients to submit event entries and by servers
// to receive those events.
type TransportConfig struct {
	Driver string `json:"driver"`
	URL    string `json:"url"`
}
