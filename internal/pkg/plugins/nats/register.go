package nats

import "github.com/SierraSoftworks/heimdall/pkg/plugins"

func init() {
	plugins.RegisterCollector("nats", NewNATSCollector)
	plugins.RegisterPublisher("nats", NewNATSPublisher)
}
