package memory

import "github.com/SierraSoftworks/heimdall/pkg/plugins"

func init() {
	plugins.RegisterStore("memory", NewMemoryStore)
}
