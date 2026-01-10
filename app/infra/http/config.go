package http

import (
	"github.com/brunojet/my-store-go/app/infra/config/adapters/env"
	"github.com/brunojet/my-store-go/app/infra/config/ports"
)

// Config holds HTTP runtime selection config.
//
// Env:
//   - HTTP_DRIVER (default: "gin")
type Config struct {
	Driver string
}

func LoadConfig(src ports.Source) Config {
	driver := ports.LowerTrimmed(src, "HTTP_DRIVER")
	if driver == "" {
		driver = "gin"
	}
	return Config{Driver: driver}
}

func ConfigFromEnv() Config {
	return LoadConfig(env.New())
}
