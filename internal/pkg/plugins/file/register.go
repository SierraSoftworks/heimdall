package file

import "github.com/SierraSoftworks/heimdall/pkg/plugins"

func init() {
	plugins.RegisterCollector("file", NewFileCollector)
	plugins.RegisterPublisher("file", NewFilePublisher)
}
