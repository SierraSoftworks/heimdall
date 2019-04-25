package config

// APIConfig describes configuration options for the
// Heimdall API.
type APIConfig struct {
	Listen string `json:"listen"`
}

// Update is used to apply a differential config
// to this config entry, updating any new fields.
func (c *APIConfig) Update(dc *APIConfig) {
	if dc.Listen != "" {
		c.Listen = dc.Listen
	}
}
