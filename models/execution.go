package models

import "time"

// Execution represents a check that has been executed by
// a Heimdall client. It includes information about the time
// taken to execute the check, the output and exit code.
type Execution struct {
	Scheduled time.Time     `json:"time"`
	Executed  time.Time     `json:"time"`
	Duration  time.Duration `json:"duration"`
	Status    Status        `json:"status"`
	Output    string        `json:"output"`
}
