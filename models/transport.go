package models

import "net/url"
import "encoding/json"
import "strings"

// Transport defines the configuration used by
// a specific transport instance. These transports are
// used by clients to submit event entries and by servers
// to receive those events.
type Transport struct {
	Driver string   `json:"driver"`
	URL    *url.URL `json:"url"`
}

func (t *Transport) MarshalMap() (interface{}, error) {
	return map[string]string{
		"driver": t.Driver,
		"url":    t.URL.String(),
	}, nil
}

func (t *Transport) UnmarshalMap(from map[string]string) error {
	for k, v := range from {
		if strings.ToLower(k) == "driver" {
			t.Driver = v
		} else if strings.ToLower(k) == "url" {
			u, err := url.Parse(v)
			if err != nil {
				return err
			}
			t.URL = u
		}
	}

	return nil
}

func (t *Transport) MarshalJSON(j []byte) ([]byte, error) {
	obj, err := t.MarshalMap()
	if err != nil {
		return nil, err
	}
	return json.Marshal(obj)
}

func (t *Transport) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	return t.UnmarshalMap(rawStrings)
}

func (t *Transport) MarshalYAML() (interface{}, error) {
	return t.MarshalMap()
}

func (t *Transport) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var rawValues map[string]string
	err := unmarshal(&rawValues)
	if err != nil {
		return err
	}

	return t.UnmarshalMap(rawValues)
}
