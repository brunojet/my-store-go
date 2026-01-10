package interfaces

import "gorm.io/gorm"

type ConnectorInterface interface {
	Open() (*gorm.DB, error)
	Close() error
}
