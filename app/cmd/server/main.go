package main

import (
	"log"
	"net/http"
	"os"
	"time"

	coremodels "github.com/brunojet/my-store-go/app/core/models"
	appsrepo "github.com/brunojet/my-store-go/app/core/repositories"
	"github.com/brunojet/my-store-go/app/infra/persistence"
	appsrouters "github.com/brunojet/my-store-go/app/store-provider/routers"
	appsservice "github.com/brunojet/my-store-go/app/store-provider/services/apps"
	"github.com/gin-gonic/gin"
)

func main() {
	conn, err := persistence.SelectConnectorFromEnv()
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	db, err := conn.Open()
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	}()

	if err := db.AutoMigrate(coremodels.MigratableModels()...); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	appsSvc := appsservice.NewService(appsrepo.New(db))

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	v1 := r.Group("/v1")
	appsrouters.RegisterAppsRoutes(v1, appsSvc)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("server listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
