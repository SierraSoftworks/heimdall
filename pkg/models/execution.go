package models

import (
	"time"

	"github.com/SierraSoftworks/heimdall/pkg/duration"
)

// Execution represents a check that has been executed by
// a Heimdall client. It includes information about the time
// taken to execute the check, the output and exit code.
type Execution struct {
	Check *Check  `json:"check"`
	Host  *Source `json:"host"`

	Scheduled time.Time         `json:"scheduled"`
	Executed  time.Time         `json:"executed"`
	Duration  duration.Duration `json:"duration"`
	Status    Status            `json:"status"`
	Output    string            `json:"output"`
}
