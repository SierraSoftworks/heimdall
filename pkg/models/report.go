package models

import (
	"bytes"
	"encoding/json"
)

// Report represents a message submitted by a check client
// over the message transport to the Heimdall server. It includes
// the details of the client and the check that was executed.
type Report struct {
	SequenceNumber int64 `json:"sno"`

	Source    *Source    `json:"source,omitempty"`
	Execution *Execution `json:"execution"`
}

func (r *Report) ToMap() map[string]interface{} {
	buf := bytes.NewBuffer([]byte{})
	json.NewEncoder(buf).Encode(r)

	var m map[string]interface{}
	json.NewDecoder(buf).Decode(&m)

	return m
}
