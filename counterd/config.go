package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/hcl"
)

// Config is the configuration for the server and snapshot comments
type Config struct {
	// ListenAddress is the HTTP listener address
	ListenAddress string `hcl:"listen_address"`

	// RedisAddress is the address of the redis server
	RedisAddress string `hcl:"redis_address"`

	// PGAddress is the address of the postgresql server
	PGAddress string `hcl:"postgresql_address"`

	// Snapshot has the snapshot specific configuration
	Snapshot *SnapshotConfig
}

// SnapshotConfig has snapshotting configuration
type SnapshotConfig struct {
	// Cron can be configured to have the server invoke snapshots periodically.
	// This is independent from invoking the snapshot command.
	Cron string `hcl:"cron"`

	// UpdateThreshold is how far back we scan for relevant updates.
	// This prevents old counters from being updated. This should be relative to the
	// snapshot rate. For example, if you snapshot hourly, consider a two hour update threshold.
	UpdateThresholdRaw string        `hcl:"update_threshold"`
	UpdateThreshold    time.Duration `hcl:"-"`

	// DeleteThreshold is how far back a counter needs to be for deletion.
	// This should be at least 2x the longest interval that is tracked. For example,
	// if monthly counters are enabled, consider a two month delete threshold.
	DeleteThresholdRaw string        `hcl:"delete_threshold"`
	DeleteThreshold    time.Duration `hcl:"-"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		ListenAddress: "127.0.0.1:8001",
		RedisAddress:  "127.0.0.1:6379",
		PGAddress:     "postgres://postgres@localhost/postgres?sslmode=disable",
		Snapshot: &SnapshotConfig{
			UpdateThreshold: 3 * time.Hour,
			DeleteThreshold: 3 * 30 * 24 * time.Hour,
		},
	}
}

// ParseConfig is used to parse the configuration
func ParseConfig(raw string) (*Config, error) {
	config := DefaultConfig()

	// Attempt to decode the configuration
	if err := hcl.Decode(config, raw); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	if raw := config.Snapshot.UpdateThresholdRaw; raw != "" {
		dur, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse duration: %v", err)
		}
		config.Snapshot.UpdateThreshold = dur
	}
	if raw := config.Snapshot.DeleteThresholdRaw; raw != "" {
		dur, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse duration: %v", err)
		}
		config.Snapshot.DeleteThreshold = dur
	}
	return config, nil
}
