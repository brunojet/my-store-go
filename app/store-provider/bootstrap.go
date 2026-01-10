package storeprovider

import (
	repositories "github.com/brunojet/my-store-go/app/core/repositories"
	"github.com/brunojet/my-store-go/app/infra/http/contracts"
	routers "github.com/brunojet/my-store-go/app/store-provider/routers"
	"github.com/brunojet/my-store-go/app/store-provider/services"
	"gorm.io/gorm"
)

// Register wires store-provider routes and services into the given router.
func Register(root contracts.Router, db *gorm.DB) {
	appsSvc := services.NewService(repositories.New(db))
	v1 := root.Group("/v1")
	routers.RegisterAppsRoutes(v1, appsSvc)
}
