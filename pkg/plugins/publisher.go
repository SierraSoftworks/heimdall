package plugins

import (
	"fmt"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
)

// Publisher is the common interface used to enable different
// publishers to integrate with Heimdall's pipeline.
type Publisher interface {
	Driver() *driver.Driver
	Describe() string
	Publish(channel string, data []byte) error
	Close() error
}

var publishers = map[string]func(*driver.Driver) (Publisher, error){}

func RegisterPublisher(name string, constructor func(*driver.Driver) (Publisher, error)) {
	publishers[name] = constructor
}

func GetPublisher(driver *driver.Driver) (Publisher, error) {
	constructor, ok := publishers[driver.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown publisher driver '%s'", driver.Type)
	}

	return constructor(driver)
}
