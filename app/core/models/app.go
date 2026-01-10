package models

import "gorm.io/gorm"

// App is the GORM model for the Aplicativo entity.
// We embed gorm.Model to get ID, CreatedAt, UpdatedAt and soft-delete fields.
type App struct {
	gorm.Model
	Nome                  string `gorm:"not null"`
	Descricao             string
	CodigoParceiroExterno string `gorm:"not null"`
}
