package routers

import (
	"github.com/brunojet/my-store-go/app/infra/http/contracts"
	"github.com/brunojet/my-store-go/app/store-provider/handlers"
	"github.com/brunojet/my-store-go/app/store-provider/services"
)

func RegisterAppsRoutes(r contracts.Router, svc *services.AppsService) {
	h := handlers.NewAppsHandler(svc)

	r.GET("/apps", h.List)
	r.GET("/apps/:id", h.Get)
	r.POST("/apps", h.Create)
	r.PATCH("/apps/:id", h.Patch)
	r.DELETE("/apps/:id", h.Delete)
}
