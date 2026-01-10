package ports

import "time"

// Audit existe também em DTO porque o contrato da API pode divergir do modelo de dados.

type AuditDTO struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
