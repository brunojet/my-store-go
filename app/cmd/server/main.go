package main

import (
	"log"
	"net/http"

	core "github.com/brunojet/my-store-go/app/core"
	httpruntime "github.com/brunojet/my-store-go/app/infra/http"
	"github.com/brunojet/my-store-go/app/infra/http/contracts"
	"github.com/brunojet/my-store-go/app/infra/observability"
	"github.com/brunojet/my-store-go/app/infra/persistence"
	storeprovider "github.com/brunojet/my-store-go/app/store-provider"
)

func main() {
	cfg := configFromEnv()

	dbMgr, err := persistence.NewDatabaseManager(cfg.Database)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}

	db, err := dbMgr.OpenAndMigrate(core.Register)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer func() {
		if err := dbMgr.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	}()

	obsMgr := observability.NewObservabilityManager(cfg.Observability)
	mws, err := obsMgr.Open()
	if err != nil {
		log.Fatalf("obs init: %v", err)
	}
	defer func() {
		if err := obsMgr.Close(); err != nil {
			log.Printf("close obs: %v", err)
		}
	}()

	httpParams := cfg.HTTP
	httpParams.ObservabilityMiddlewares = mws
	httpMgr, err := httpruntime.NewHTTPManager(httpParams)
	if err != nil {
		log.Fatalf("http init: %v", err)
	}

	httpRuntime, err := httpMgr.OpenAndRegister(func(r contracts.Router) error {
		r.GET("/health", func(c contracts.Context) {
			c.String(200, "ok")
		})
		storeprovider.Register(r, db)
		return nil
	})
	if err != nil {
		log.Fatalf("http register: %v", err)
	}

	srv := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           httpRuntime.Handler,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
	}
	log.Printf("server listening on %s", cfg.Server.Addr)
	log.Fatal(srv.ListenAndServe())
}
