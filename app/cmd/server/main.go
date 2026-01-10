package main

import (
	"log"
	"net/http"
	"os"
	"time"

	core "github.com/brunojet/my-store-go/app/core"
	httpruntime "github.com/brunojet/my-store-go/app/infra/http"
	"github.com/brunojet/my-store-go/app/infra/http/contracts"
	"github.com/brunojet/my-store-go/app/infra/persistence"
	storeprovider "github.com/brunojet/my-store-go/app/store-provider"
)

func main() {
	conn, err := persistence.SelectFromEnv()
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

	httpRuntime, err := httpruntime.SelectFromEnv()
	if err != nil {
		log.Fatalf("http init: %v", err)
	}

	httpRuntime.Router.GET("/health", func(c contracts.Context) {
		c.String(200, "ok")
	})

	storeprovider.Register(httpRuntime.Router, db)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           httpRuntime.Handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("server listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
