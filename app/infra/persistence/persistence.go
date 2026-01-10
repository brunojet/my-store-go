package persistence

import (
	"fmt"
	"os"
	"strings"

	mysqlp "github.com/brunojet/my-store-go/app/infra/persistence/adapters/mysql"
	sqlitep "github.com/brunojet/my-store-go/app/infra/persistence/adapters/sqlite"
	"github.com/brunojet/my-store-go/app/infra/persistence/contracts"
)

// SelectFromEnv selects a DB connector implementation based on environment variables.
//
// Supported:
//   - DB_DRIVER: "sqlite" (default) | "mysql"
//   - DB_DSN: driver-specific DSN
//
// The returned connector is NOT opened yet; call Open() and defer Close() in main.
func SelectFromEnv() (contracts.DatabaseConnector, error) {
	driver := strings.TrimSpace(os.Getenv("DB_DRIVER"))
	if driver == "" {
		driver = "sqlite"
	}
	driver = strings.ToLower(driver)

	dsn := strings.TrimSpace(os.Getenv("DB_DSN"))

	switch driver {
	case "sqlite":
		return sqlitep.New(dsn), nil
	case "mysql":
		return mysqlp.New(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER %q", driver)
	}
}
