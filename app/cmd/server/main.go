package main

import (
	"log"
	"net/http"

	core "github.com/brunojet/my-store-go/app/core"
	httpruntime "github.com/brunojet/my-store-go/app/infra/http"
	"github.com/brunojet/my-store-go/app/infra/http/contracts"
	"github.com/brunojet/my-store-go/app/infra/observability"
	"github.com/brunojet/my-store-go/app/infra/persistence"
	server "github.com/brunojet/my-store-go/app/infra/server"
	storeprovider "github.com/brunojet/my-store-go/app/store-provider"
)

func main() {
	serverCfg := server.ConfigFromEnv()
	httpCfg := httpruntime.ConfigFromEnv()
	obsCfg := observability.ConfigFromEnv()
	persistCfg := persistence.ConfigFromEnv()

	conn, err := persistence.Select(persistCfg)
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

	if err := core.Register(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	mws := observability.HTTPMiddlewares(httpCfg.Driver, obsCfg)
	httpRuntime, err := httpruntime.Select(
		httpCfg.Driver,
		httpruntime.WithMiddlewares(mws),
		httpruntime.WithCORS(httpCfg.CORS),
	)

	if err != nil {
		log.Fatalf("http init: %v", err)
	}

	httpRuntime.Router.GET("/health", func(c contracts.Context) {
		c.String(200, "ok")
	})

	storeprovider.Register(httpRuntime.Router, db)

	srv := &http.Server{
		Addr:              serverCfg.Addr,
		Handler:           httpRuntime.Handler,
		ReadHeaderTimeout: serverCfg.ReadHeaderTimeout,
	}
	log.Printf("server listening on %s", serverCfg.Addr)
	log.Fatal(srv.ListenAndServe())
}
