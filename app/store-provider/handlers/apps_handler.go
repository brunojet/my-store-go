package handlers

import (
	appsdtos "github.com/brunojet/my-store-go/app/core/dtos"
	"github.com/brunojet/my-store-go/app/core/models"
	httpports "github.com/brunojet/my-store-go/app/infra/http/ports"
	services "github.com/brunojet/my-store-go/app/store-provider/services"
)

type AppsHandler struct {
	*httpports.CRUDHandler[models.App, appsdtos.CreateAppRequest, appsdtos.PatchAppRequest, appsdtos.AppResponse]
}

func NewAppsHandler(svc *services.AppsService) *AppsHandler {
	return &AppsHandler{
		CRUDHandler: &httpports.CRUDHandler[models.App, appsdtos.CreateAppRequest, appsdtos.PatchAppRequest, appsdtos.AppResponse]{
			Service:         svc,
			ToResponse:      appsdtos.ToAppResponse,
			NotFoundError:   services.ErrNotFound,
			ValidationError: services.ErrValidation,
		},
	}
}
