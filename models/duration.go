package models

import (
	"fmt"
	"time"
)

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Duration(d).String())), nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) < 2 {
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

	pd, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = Duration(pd)

	return nil
}
