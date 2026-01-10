package dtos

import (
	"time"

	"github.com/brunojet/my-store-go/app/core/models"
)

type CreateAppRequest struct {
	Nome                  string `json:"nome"`
	Descricao             string `json:"descricao"`
	CodigoParceiroExterno string `json:"codigo_parceiro_externo"`
}

type PatchAppRequest struct {
	Nome                  *string `json:"nome"`
	Descricao             *string `json:"descricao"`
	CodigoParceiroExterno *string `json:"codigo_parceiro_externo"`
}

type AppResponse struct {
	ID                    uint      `json:"id"`
	Nome                  string    `json:"nome"`
	Descricao             string    `json:"descricao"`
	CodigoParceiroExterno string    `json:"codigo_parceiro_externo"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func ToAppResponse(a models.App) AppResponse {
	return AppResponse{
		ID:                    a.ID,
		Nome:                  a.Nome,
		Descricao:             a.Descricao,
		CodigoParceiroExterno: a.CodigoParceiroExterno,
		CreatedAt:             a.CreatedAt,
		UpdatedAt:             a.UpdatedAt,
	}
}
