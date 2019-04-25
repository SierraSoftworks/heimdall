package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/SierraSoftworks/connor"
	log "github.com/Sirupsen/logrus"
)

// Driver is used to represent a configurable driver
// which offers functionality used by Heimdall.
// Drivers are used to provide functionality for
// collecting status information, publishing it
// and storing it. This Driver construct is used
// in the configuration files to provide a common
// configuration interface for all of the above.
type Driver struct {
	ID     string                 `json:"id"`
	Type   string                 `json:"type"`
	URL    *url.URL               `json:"url"`
	Filter map[string]interface{} `json:"filter"`
}

// Describe will create a string representation of a driver
// object for inclusion in logs and other output elements.
func (t *Driver) Describe() string {
	if t.ID != "" {
		return t.ID
	}

	out := t.Type

	if t.URL != nil {
		out = out + "(" + t.URL.String() + ")"
	}

	if t.Filter != nil {
		buf := bytes.NewBuffer([]byte{})
		if err := json.NewEncoder(buf).Encode(t.Filter); err == nil {
			out = out + " " + strings.TrimSpace(buf.String())
		}
	}

	return out
}

// Matches determines whether this driver's filter matches a
// provided context object.
func (t *Driver) Matches(context map[string]interface{}) bool {
	if t.Filter == nil {
		return true
	}

	match, err := connor.Match(t.Filter, context)
	if err != nil {
		log.
			WithField("filter", t.Filter).
			WithField("context", context).
			WithError(err).
			Debug("Failed to run filter against context")

		return false
	}

	return match
}

// Equals determines whether two driver objects are equivalent or not
// ID equivalency takes precidence where it is provided, otherwise value
// equality will be used.
func (t *Driver) Equals(o *Driver) bool {
	if t.ID != "" {
		return t.ID == o.ID
	}

	return t.Type == o.Type &&
		t.URL.String() == o.URL.String() &&
		reflect.DeepEqual(t.Filter, o.Filter)
}

func (t *Driver) SafeURLString() string {
	return strings.Trim(fmt.Sprintf("%s/%s", t.SafeURLHost(), strings.TrimLeft(t.URL.Path, "/")), "/:")
}

func (t *Driver) SafeURLHost() string {
	return strings.Trim(fmt.Sprintf("%s://%s", t.URL.Scheme, t.URL.Host), "/:")
}

func (t *Driver) GetPath(suffix string) string {

	return strings.Join([]string{
		strings.TrimRight(t.URL.Path, "/"),
		strings.TrimLeft(suffix, "/"),
	}, "/")
}

func (t *Driver) GetOption(name, defaultValue string) string {
	v := t.URL.Query().Get(name)
	if v == "" {
		return defaultValue
	}

	return v
}

func (t *Driver) MarshalMap() (interface{}, error) {
	m := map[string]interface{}{
		"id":   t.ID,
		"type": t.Type,
		"url":  t.URL.String(),
	}

	if t.Filter != nil && len(t.Filter) > 0 {
		m["filter"] = t.Filter
	}

	return m, nil
}

func (t *Driver) UnmarshalMap(from map[string]interface{}) error {
	for k, v := range from {
		switch strings.ToLower(k) {
		case "id":
			if vstring, ok := v.(string); ok {
				t.ID = vstring
			} else {
				return fmt.Errorf("Driver ID must be a string, got %#v", v)
			}
		case "type":
			if vstring, ok := v.(string); ok {
				t.Type = vstring
			} else {
				return fmt.Errorf("Driver type must be a string, got %#v", v)
			}
		case "url":
			if vstring, ok := v.(string); ok {
				u, err := url.Parse(vstring)
				if err != nil {
					return err
				}
				t.URL = u
			} else {
				return fmt.Errorf("Driver url must be a string, got %#v", v)
			}
		case "filter":
			if vmap, ok := v.(map[string]interface{}); ok {
				t.Filter = vmap
			} else {
				return fmt.Errorf("Driver filter must be a dictionary")
			}
		default:
		}
	}

	return nil
}

func (t *Driver) MarshalJSON() ([]byte, error) {
	obj, err := t.MarshalMap()
	if err != nil {
		return nil, err
	}
	return json.Marshal(obj)
}

func (t *Driver) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	return t.UnmarshalMap(rawStrings)
}

func (t *Driver) MarshalYAML() (interface{}, error) {
	return t.MarshalMap()
}

func (t *Driver) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var rawValues map[string]interface{}
	err := unmarshal(&rawValues)
	if err != nil {
		return err
	}

	return t.UnmarshalMap(rawValues)
}
