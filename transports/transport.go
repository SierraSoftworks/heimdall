package transports

import (
	"fmt"
	"net/url"
)

type Transport interface {
	Describe() string
	Subscribe(topic string) (Subscription, error)
	Publish(topic string, data []byte) error
	Close() error
}

type Subscription interface {
	Channel() <-chan []byte
	Close() error
}

func GetTransport(driver string, u *url.URL) (Transport, error) {
	switch driver {
	case "nats":
		return NewNATSTransport(u)
	case "redis":
		return NewRedisTransport(u)
	case "file":
		return NewFileTransport(u)
	default:
		return nil, fmt.Errorf("unknown transport driver")
	}
}
