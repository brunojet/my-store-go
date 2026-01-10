package observability

import (
	"github.com/brunojet/my-store-go/app/infra/config/adapters/env"
	"github.com/brunojet/my-store-go/app/infra/config/ports"
)

func LoadConfig(src ports.Source) Config {
	cfg := Config{RequestID: true, Telemetry: true, Recovery: true}

	cfg.RequestID = ports.Bool(src, "OBS_REQUEST_ID", cfg.RequestID)

	// Backward-compatible: OBS_ACCESS_LOG previously controlled request logs.
	// If OBS_TELEMETRY is explicitly set, it wins.
	if _, ok := src.Lookup("OBS_TELEMETRY"); ok {
		cfg.Telemetry = ports.Bool(src, "OBS_TELEMETRY", cfg.Telemetry)
	} else {
		cfg.Telemetry = ports.Bool(src, "OBS_ACCESS_LOG", cfg.Telemetry)
	}

	cfg.Recovery = ports.Bool(src, "OBS_RECOVERY", cfg.Recovery)
	return cfg
}

// ConfigFromEnv builds an observability Config from env.
func ConfigFromEnv() Config {
	return LoadConfig(env.New())
}
