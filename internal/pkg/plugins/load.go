package plugins

import (
	// Load the file plugin
	_ "github.com/SierraSoftworks/heimdall/internal/pkg/plugins/file"

	// Load the memory plugin
	_ "github.com/SierraSoftworks/heimdall/internal/pkg/plugins/memory"

	// Load the NATS plugin
	_ "github.com/SierraSoftworks/heimdall/internal/pkg/plugins/nats"

	// Load the Redis plugin
	_ "github.com/SierraSoftworks/heimdall/internal/pkg/plugins/redis"
)
