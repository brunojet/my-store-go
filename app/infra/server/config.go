package server

import (
	"time"

	"github.com/brunojet/my-store-go/app/infra/config/adapters/env"
	"github.com/brunojet/my-store-go/app/infra/config/ports"
)

// Config holds HTTP server configuration.
//
// Env:
//   - PORT (default: 8080)
//   - READ_HEADER_TIMEOUT (default: 5s)
type Config struct {
	Addr              string
	ReadHeaderTimeout time.Duration
}

func LoadConfig(src ports.Source) Config {
	port := ports.Trimmed(src, "PORT")
	if port == "" {
		port = "8080"
	}
	return Config{
		Addr:              ":" + port,
		ReadHeaderTimeout: ports.Duration(src, "READ_HEADER_TIMEOUT", 5*time.Second),
	}
}

func ConfigFromEnv() Config {
	return LoadConfig(env.New())
}
