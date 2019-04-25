// The duration package provides a JSON serializable time.Duration
// type which allows you to specify durations in human-readable formats.

package duration

import (
	"fmt"
	"strconv"
	"time"
)

// Duration is a time.Duration object representing a number of nanoseconds.
// It may be serialized to and deserialized from JSON in a human readable format
// like '10s' or '1h30m'.
type Duration time.Duration

// Equals determines whether two duration objects are identical or not
func (d Duration) Equals(o Duration) bool {
	return int64(d) == int64(o)
}

// MarshalJSON is responsible for converting the duration object into a human readable
// string.
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Duration(d).String())), nil
}

// UnmarshalJSON is responsible for converting a human readable duration string into
// a Duration object.
func (d *Duration) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) == 0 {
		return fmt.Errorf("not a valid duration")
	}

	if s[0] == '"' && s[len(s)-1] == '"' {
		pd, err := time.ParseDuration(s[1 : len(s)-1])
		if err != nil {
			return err
		}

		*d = Duration(pd)
		return nil
	}

	pd, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*d = Duration(time.Duration(pd) * time.Second)

	return nil
}

func (d Duration) MarshalYAML() (interface{}, error) {
	return time.Duration(d).String(), nil
}

func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	pd, err := time.ParseDuration(str)
	if err != nil {
		return err
	}

	*d = Duration(pd)

	return nil
}
