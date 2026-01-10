package routers

import (
	"github.com/brunojet/my-store-go/app/infra/http/contracts"
	"github.com/brunojet/my-store-go/app/store-provider/handlers"
	appsservice "github.com/brunojet/my-store-go/app/store-provider/services/apps"
)

func RegisterAppsRoutes(r contracts.Router, svc *appsservice.Service) {
	h := handlers.NewAppsHandler(svc)

	r.GET("/apps", h.List)
	r.GET("/apps/:id", h.Get)
	r.POST("/apps", h.Create)
	r.PATCH("/apps/:id", h.Patch)
	r.DELETE("/apps/:id", h.Delete)
}
