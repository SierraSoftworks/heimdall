package models

// Report represents a message submitted by a check client
// over the message transport to the Heimdall server. It includes
// the details of the client and the check that was executed.
type Report struct {
	Client    *Client    `json:"client"`
	Check     *Check     `json:"check"`
	Execution *Execution `json:"execution"`
}
