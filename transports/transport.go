package transports

import "fmt"

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

func GetTransport(url string) (Transport, error) {
	u, err := ParseURL(url)
	if err != nil {
		return nil, err
	}

	switch u.Driver {
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
