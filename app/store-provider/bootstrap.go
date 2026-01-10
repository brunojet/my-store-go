package storeprovider

import (
	appsrepo "github.com/brunojet/my-store-go/app/core/repositories"
	"github.com/brunojet/my-store-go/app/infra/http/contracts"
	appsrouters "github.com/brunojet/my-store-go/app/store-provider/routers"
	appsservice "github.com/brunojet/my-store-go/app/store-provider/services/apps"
	"gorm.io/gorm"
)

// Register wires store-provider routes and services into the given router.
func Register(root contracts.Router, db *gorm.DB) {
	appsSvc := appsservice.NewService(appsrepo.New(db))
	v1 := root.Group("/v1")
	appsrouters.RegisterAppsRoutes(v1, appsSvc)
}
