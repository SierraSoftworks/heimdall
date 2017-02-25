package models

import "time"

// Client represents a Heimdall client which will execute checks
// and report them, over its configured transport, to the Heimdall
// server(s).
// Clients are uniquely identified by their name and may also include
// tags which may be used to perform filtering of check results across
// your infrastructure.
type Client struct {
	Name string            `json:"name"`
	Tags map[string]string `json:"tags"`
}

// ClientDetails represents a client and a rollup of its state to enable
// quick and simple rendering of a client's information on dashboards
// and other interfaces.
type ClientDetails struct {
	*Client
	Status   Status    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}
