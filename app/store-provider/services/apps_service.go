package services

import (
	"context"
	"strings"

	dtos "github.com/brunojet/my-store-go/app/core/dtos"
	models "github.com/brunojet/my-store-go/app/core/models"
	appsrepo "github.com/brunojet/my-store-go/app/core/repositories"
)

type AppsService struct {
	repo appsrepo.Repository
}

func NewService(repo appsrepo.Repository) *AppsService {
	return &AppsService{repo: repo}
}

func (s *AppsService) List(ctx context.Context) ([]models.App, error) {
	return s.repo.List(ctx)
}

func (s *AppsService) Get(ctx context.Context, id uint) (*models.App, error) {
	a, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *AppsService) Create(ctx context.Context, req dtos.CreateAppRequest) (*models.App, error) {
	nome := strings.TrimSpace(req.Nome)
	codigo := strings.TrimSpace(req.CodigoParceiroExterno)
	if nome == "" || codigo == "" {
		return nil, ErrValidation
	}

	a := &models.App{
		Nome:                  nome,
		Descricao:             req.Descricao,
		CodigoParceiroExterno: codigo,
	}
	return s.repo.Create(ctx, a)
}

func (s *AppsService) Patch(ctx context.Context, id uint, req dtos.PatchAppRequest) (*models.App, error) {
	if req.Nome == nil && req.Descricao == nil && req.CodigoParceiroExterno == nil {
		return nil, ErrValidation
	}

	updates := map[string]any{}
	if req.Descricao != nil {
		updates["descricao"] = *req.Descricao
	}

	if req.Nome != nil {
		nome := strings.TrimSpace(*req.Nome)
		if nome == "" {
			return nil, ErrValidation
		}
		updates["nome"] = nome
	}

	if req.CodigoParceiroExterno != nil {
		codigo := strings.TrimSpace(*req.CodigoParceiroExterno)
		if codigo == "" {
			return nil, ErrValidation
		}
		updates["codigo_parceiro_externo"] = codigo
	}

	a, err := s.repo.Patch(ctx, id, updates)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *AppsService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
