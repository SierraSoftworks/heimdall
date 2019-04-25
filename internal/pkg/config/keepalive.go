package config

import (
	"time"
)

type KeepaliveConfig struct {
	Interval    time.Duration `json:"interval"`
	Collections []string      `json:"collections"`
}

func (c *KeepaliveConfig) Update(dc *KeepaliveConfig) {
	if dc.Interval != 0 {
		c.Interval = dc.Interval
	}

	if dc.Collections != nil {
		c.Collections = append(c.Collections, dc.Collections...)
	}
}
