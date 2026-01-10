package persistence

import (
	"fmt"
	"strings"

	mysqlp "github.com/brunojet/my-store-go/app/infra/persistence/adapters/mysql"
	sqlitep "github.com/brunojet/my-store-go/app/infra/persistence/adapters/sqlite"
	"github.com/brunojet/my-store-go/app/infra/persistence/contracts"
)

// Select selects a DB connector implementation based on config.
func Select(cfg Config) (contracts.DatabaseConnector, error) {
	driver := strings.TrimSpace(strings.ToLower(cfg.Driver))
	if driver == "" {
		driver = "sqlite"
	}

	dsn := strings.TrimSpace(cfg.DSN)

	switch driver {
	case "sqlite":
		return sqlitep.New(dsn), nil
	case "mysql":
		return mysqlp.New(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER %q", cfg.Driver)
	}
}

// SelectFromEnv selects a DB connector implementation based on environment variables.
//
// Supported:
//   - DB_DRIVER: "sqlite" (default) | "mysql"
//   - DB_DSN: driver-specific DSN
//
// The returned connector is NOT opened yet; call Open() and defer Close() in main.
func SelectFromEnv() (contracts.DatabaseConnector, error) {
	return Select(ConfigFromEnv())
}
