package persistence

import (
	"github.com/brunojet/my-store-go/app/infra/config/adapters/env"
	"github.com/brunojet/my-store-go/app/infra/config/ports"
)

// Config holds persistence selection config.
//
// Env:
//   - DB_DRIVER (default: "sqlite")
//   - DB_DSN (default: ""; sqlite adapter will fall back to in-memory)
type Config struct {
	Driver string
	DSN    string
}

func LoadConfig(src ports.Source) Config {
	driver := ports.LowerTrimmed(src, "DB_DRIVER")
	if driver == "" {
		driver = "sqlite"
	}
	return Config{Driver: driver, DSN: ports.Trimmed(src, "DB_DSN")}
}

func ConfigFromEnv() Config {
	return LoadConfig(env.New())
}
