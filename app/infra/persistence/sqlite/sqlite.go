package sqlite

import (
	"database/sql"
	"sync"

	"github.com/brunojet/my-store-go/app/infra/persistence/interfaces"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const defaultInMemoryDSN = "file::memory:?cache=shared"

type Connector struct {
	DSN string

	mu    sync.Mutex
	gorm  *gorm.DB
	sqlDB *sql.DB
}

var _ interfaces.ConnectorInterface = (*Connector)(nil)

func New(dsn string) *Connector {
	return &Connector{DSN: dsn}
}

func (c *Connector) Open() (*gorm.DB, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.gorm != nil {
		return c.gorm, nil
	}

	dsn := c.DSN
	if dsn == "" {
		dsn = defaultInMemoryDSN
	}

	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	// SQLite works best with a single open connection.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	c.gorm = gormDB
	c.sqlDB = sqlDB
	return gormDB, nil
}

func (c *Connector) DB() *gorm.DB {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.gorm
}

func (c *Connector) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sqlDB == nil {
		return nil
	}
	err := c.sqlDB.Close()
	c.sqlDB = nil
	c.gorm = nil
	return err
}
