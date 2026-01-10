package handlers

import (
	appsdtos "github.com/brunojet/my-store-go/app/core/dtos"
	"github.com/brunojet/my-store-go/app/core/models"
	ginbase "github.com/brunojet/my-store-go/app/infra/http/gin"
	appsservice "github.com/brunojet/my-store-go/app/store-provider/services/apps"
)

type AppsHandler struct {
	*ginbase.CRUDHandler[models.App, appsdtos.CreateAppRequest, appsdtos.PatchAppRequest, appsdtos.AppResponse]
}

func NewAppsHandler(svc *appsservice.Service) *AppsHandler {
	return &AppsHandler{
		CRUDHandler: &ginbase.CRUDHandler[models.App, appsdtos.CreateAppRequest, appsdtos.PatchAppRequest, appsdtos.AppResponse]{
			Service:         svc,
			ToResponse:      appsdtos.ToAppResponse,
			NotFoundError:   appsservice.ErrNotFound,
			ValidationError: appsservice.ErrValidation,
		},
	}
}
