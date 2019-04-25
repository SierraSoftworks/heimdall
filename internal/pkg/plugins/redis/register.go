package redis

import "github.com/SierraSoftworks/heimdall/pkg/plugins"

func init() {
	plugins.RegisterCollector("redis", NewRedisCollector)
	plugins.RegisterPublisher("redis", NewRedisPublisher)
}
