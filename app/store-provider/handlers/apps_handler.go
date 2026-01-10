package handlers

import (
	appsdtos "github.com/brunojet/my-store-go/app/core/dtos"
	"github.com/brunojet/my-store-go/app/core/models"
	coretel "github.com/brunojet/my-store-go/app/core/telemetry"
	httpbase "github.com/brunojet/my-store-go/app/infra/http/base"
	appsservice "github.com/brunojet/my-store-go/app/store-provider/services/apps"
)

type AppsHandler struct {
	*httpbase.CRUDHandler[models.App, appsdtos.CreateAppRequest, appsdtos.PatchAppRequest, appsdtos.AppResponse]
}

func NewAppsHandler(svc *appsservice.Service, tel coretel.Provider) *AppsHandler {
	return &AppsHandler{
		CRUDHandler: &httpbase.CRUDHandler[models.App, appsdtos.CreateAppRequest, appsdtos.PatchAppRequest, appsdtos.AppResponse]{
			Name:            "apps",
			Telemetry:       tel,
			Service:         svc,
			ToResponse:      appsdtos.ToAppResponse,
			NotFoundError:   appsservice.ErrNotFound,
			ValidationError: appsservice.ErrValidation,
		},
	}
}
