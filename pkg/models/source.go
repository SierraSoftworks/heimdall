package models

// Source represents a Heimdall client which will execute checks
// and report them, over its configured transport, to the Heimdall
// server(s).
// Clients are uniquely identified by their name and may also include
// tags which may be used to perform filtering of check results across
// your infrastructure.
type Source struct {
	ID   string            `json:"id"`
	Tags map[string]string `json:"tags"`
}
