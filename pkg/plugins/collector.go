package plugins

import (
	"fmt"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
)

// Collector is the common interface used to enable
// different collectors to be used by Heimdall's
// pipeline.
type Collector interface {
	Driver() *driver.Driver
	Describe() string
	Subscribe(channel string) (Subscription, error)
	Close() error
}

// Subscription represents an active collector subscription
// which emits check status information on a channel.
type Subscription interface {
	Channel() <-chan []byte
	Close() error
}

var collectors = map[string]func(*driver.Driver) (Collector, error){}

func RegisterCollector(name string, constructor func(*driver.Driver) (Collector, error)) {
	collectors[name] = constructor
}

func GetCollector(driver *driver.Driver) (Collector, error) {
	constructor, ok := collectors[driver.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown collector driver '%s'", driver.Type)
	}

	return constructor(driver)
}
