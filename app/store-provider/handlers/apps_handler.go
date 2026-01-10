package handlers

import (
	appsdtos "github.com/brunojet/my-store-go/app/core/dtos"
	"github.com/brunojet/my-store-go/app/core/models"
	httpbase "github.com/brunojet/my-store-go/app/infra/http/base"
	services "github.com/brunojet/my-store-go/app/store-provider/services"
)

type AppsHandler struct {
	*httpbase.CRUDHandler[models.App, appsdtos.CreateAppRequest, appsdtos.PatchAppRequest, appsdtos.AppResponse]
}

func NewAppsHandler(svc *services.AppsService) *AppsHandler {
	return &AppsHandler{
		CRUDHandler: &httpbase.CRUDHandler[models.App, appsdtos.CreateAppRequest, appsdtos.PatchAppRequest, appsdtos.AppResponse]{
			Service:         svc,
			ToResponse:      appsdtos.ToAppResponse,
			NotFoundError:   services.ErrNotFound,
			ValidationError: services.ErrValidation,
		},
	}
}
