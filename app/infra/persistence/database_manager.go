package persistence

import (
	"errors"
	"fmt"
	"strings"

	mysqlp "github.com/brunojet/my-store-go/app/infra/persistence/adapters/mysql"
	sqlitep "github.com/brunojet/my-store-go/app/infra/persistence/adapters/sqlite"
	"github.com/brunojet/my-store-go/app/infra/persistence/contracts"
	"gorm.io/gorm"
)

var ErrMissingMigrator = errors.New("migrate enabled but no migrator provided")

type DatabaseConfig struct {
	Driver string
	DSN    string
}

// SelectFromParams selects a DB connector implementation based on DatabaseParams.
//
// It reuses the existing Select(DatabaseConfig) and therefore supports the same drivers.
func SelectFromParams(params DatabaseParams) (contracts.DatabaseConnector, error) {
	dsn, err := params.BuildDSN()
	if err != nil {
		return nil, err
	}
	return Select(DatabaseConfig{Driver: params.NormalizedDriver(), DSN: dsn})
}

func Select(cfg DatabaseConfig) (contracts.DatabaseConnector, error) {
	driver := strings.TrimSpace(strings.ToLower(cfg.Driver))
	if driver == "" {
		driver = "mysql"
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

// DatabaseManager owns a DatabaseConnector lifecycle (open/close) and can
// optionally run migrations via a provided callback when Params.Migrate is true.
//
// Typical usage:
//
//	mgr, _ := persistence.NewDatabaseManager(persistence.DatabaseParams{Driver: "sqlite", Migrate: true})
//	db, _ := mgr.OpenAndMigrate(core.Register)
//	defer mgr.Close()
type DatabaseManager struct {
	Params    DatabaseParams
	Connector contracts.DatabaseConnector

	db *gorm.DB
}

func NewDatabaseManager(params DatabaseParams) (*DatabaseManager, error) {
	conn, err := SelectFromParams(params)
	if err != nil {
		return nil, err
	}
	return &DatabaseManager{Params: params, Connector: conn}, nil
}

func (m *DatabaseManager) Open() (*gorm.DB, error) {
	if m.Connector == nil {
		return nil, fmt.Errorf("nil Connector")
	}
	if m.db != nil {
		return m.db, nil
	}
	db, err := m.Connector.Open()
	if err != nil {
		return nil, err
	}
	m.db = db
	return db, nil
}

func (m *DatabaseManager) OpenAndMigrate(migrator func(*gorm.DB) error) (*gorm.DB, error) {
	db, err := m.Open()
	if err != nil {
		return nil, err
	}
	if !m.Params.Migrate {
		return db, nil
	}
	if migrator == nil {
		return nil, ErrMissingMigrator
	}
	if err := migrator(db); err != nil {
		return nil, err
	}
	return db, nil
}

func (m *DatabaseManager) Close() error {
	m.db = nil
	if m.Connector == nil {
		return nil
	}
	return m.Connector.Close()
}
