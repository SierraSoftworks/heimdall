package models

import (
	"fmt"
	"strings"
)

// Status represents the status of a check in terms of
// an exit code.
type Status int

// StatusUnkn indicates that a check is in an unknown
// state. This usually indicates that the check script
// is either operating in a non-compliant manner or that
// the state of the underlying resource or service could
// not be determined by the check.
const StatusUnkn = Status(3)

// StatusCrit indicates that a check is failing with a
// critical error. This usually indicates that a service or
// resource is entirely unavailable.
const StatusCrit = Status(2)

// StatusWarn indicates that a check is failing with a
// non-critical error. This generally indicates that a service
// or resource is operating in a degraded state.
const StatusWarn = Status(1)

// StatusOkay indicates that a check is passing and that,
// to the best of its knowledge, the underlying system or resource
// is fully available and working correctly.
const StatusOkay = Status(0)

func (s Status) String() string {
	switch s {
	case StatusCrit:
		return "CRIT"
	case StatusWarn:
		return "WARN"
	case StatusOkay:
		return "OK"
	default:
		return "UNKN"
	}
}

// ParseStatus attempts to parse a status string into a Status
// object.
func ParseStatus(status string) Status {
	switch strings.ToUpper(status) {
	case "CRIT":
		return StatusCrit
	case "WARN":
		return StatusWarn
	case "OK":
		return StatusOkay
	case "UNKN":
		return StatusUnkn
	default:
		return StatusUnkn
	}
}

// IsWorseThan tells you whether a specific status is worse than another
// given the following order of precidence  OK < UNKN < WARN < CRIT
func (s Status) IsWorseThan(s2 Status) bool {
	if s >= StatusUnkn && s2 == StatusOkay {
		return true
	}

	if s >= StatusUnkn {
		return false
	}

	return s > s2
}

// MarshalJSON will convert a CheckStatus entry into a format
// compatible with the JSON serializer.
func (s Status) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", s.String())), nil
}

// UnmarshalJSON will convert a JSON serialized CheckStatus
// entry back into its CheckStatus form.
func (s *Status) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != byte('"') || b[len(b)-1] != byte('"') {
		return fmt.Errorf("Expected a string status type")
	}

	*s = ParseStatus(string(b[1 : len(b)-1]))

	return nil
}
