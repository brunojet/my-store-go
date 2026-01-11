package main

import (
	"time"

	"github.com/brunojet/my-store-go/app/infra/config/adapters/env"
	cfgcontracts "github.com/brunojet/my-store-go/app/infra/config/contracts"
	cfgports "github.com/brunojet/my-store-go/app/infra/config/ports"
	"github.com/brunojet/my-store-go/app/infra/http"
	httptypes "github.com/brunojet/my-store-go/app/infra/http/types"
	"github.com/brunojet/my-store-go/app/infra/observability"
	"github.com/brunojet/my-store-go/app/infra/persistence"
	"github.com/brunojet/my-store-go/app/infra/server"
)

type appConfig struct {
	Server        server.ServerParams
	HTTP          http.HTTPParams
	Observability observability.ObservabilityParams
	Database      persistence.DatabaseParams
}

func configFromEnv() appConfig {
	src := env.New()
	httpParams := httpParamsFromSource(src)
	return appConfig{
		Server:        serverConfigFromSource(src),
		HTTP:          httpParams,
		Observability: observabilityParamsFromSource(src, string(httpParams.Driver)),
		Database:      databaseParamsFromSource(src),
	}
}

// Env (Server):
//   - PORT (default: 8080)
//   - READ_HEADER_TIMEOUT (default: 5s)
func serverConfigFromSource(src cfgcontracts.Source) server.ServerParams {
	port := cfgports.Trimmed(src, "PORT")
	if port == "" {
		port = "8080"
	}
	return server.ServerParams{
		Addr:              ":" + port,
		ReadHeaderTimeout: cfgports.Duration(src, "READ_HEADER_TIMEOUT", 5*time.Second),
	}
}

// Env (Observability):
//   - OBS_REQUEST_ID (default: true)
//   - OBS_TELEMETRY (default: true)
//   - OBS_ACCESS_LOG (backward-compatible alias for OBS_TELEMETRY)
//   - OBS_RECOVERY (default: true)
func observabilityParamsFromSource(src cfgcontracts.Source, httpDriver string) observability.ObservabilityParams {
	cfg := observability.ObservabilityConfig{RequestID: true, Telemetry: true, Recovery: true}

	cfg.RequestID = cfgports.Bool(src, "OBS_REQUEST_ID", cfg.RequestID)

	// Backward-compatible: OBS_ACCESS_LOG previously controlled request logs.
	// If OBS_TELEMETRY is explicitly set, it wins.
	if _, ok := src.Lookup("OBS_TELEMETRY"); ok {
		cfg.Telemetry = cfgports.Bool(src, "OBS_TELEMETRY", cfg.Telemetry)
	} else {
		cfg.Telemetry = cfgports.Bool(src, "OBS_ACCESS_LOG", cfg.Telemetry)
	}

	cfg.Recovery = cfgports.Bool(src, "OBS_RECOVERY", cfg.Recovery)

	return observability.ObservabilityParams{Driver: httpDriver, Config: cfg}
}

// Env (HTTP):
//   - HTTP_DRIVER (default: "gin")
//   - CORS_ENABLED (default: true)
//   - CORS_ALLOW_ORIGINS (default: "*")
//   - CORS_ALLOW_METHODS (default: "GET,POST,PATCH,DELETE,OPTIONS")
//   - CORS_ALLOW_HEADERS (default: "Content-Type,Authorization")
//   - CORS_EXPOSE_HEADERS (default: "")
//   - CORS_ALLOW_CREDENTIALS (default: false)
//   - CORS_MAX_AGE (default: 10m)
func httpParamsFromSource(src cfgcontracts.Source) http.HTTPParams {
	// Keep defaults aligned with infra/http/config.go.
	cors := httptypes.CORSConfig{
		Enabled:          cfgports.Bool(src, "CORS_ENABLED", true),
		AllowOrigins:     cfgports.SplitCSV(cfgports.Trimmed(src, "CORS_ALLOW_ORIGINS"), []string{"*"}),
		AllowMethods:     cfgports.SplitCSVUpper(cfgports.Trimmed(src, "CORS_ALLOW_METHODS"), []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}),
		AllowHeaders:     cfgports.SplitCSV(cfgports.Trimmed(src, "CORS_ALLOW_HEADERS"), []string{"Content-Type", "Authorization"}),
		ExposeHeaders:    cfgports.SplitCSV(cfgports.Trimmed(src, "CORS_EXPOSE_HEADERS"), nil),
		AllowCredentials: cfgports.Bool(src, "CORS_ALLOW_CREDENTIALS", false),
		MaxAge:           cfgports.Duration(src, "CORS_MAX_AGE", 10*time.Minute),
	}

	driver := http.NormalizeHTTPDriver(cfgports.Trimmed(src, "HTTP_DRIVER"))

	return http.HTTPParams{Driver: driver, CORS: cors, Register: true}
}

// Env:
//   - DB_DRIVER (default: "mysql")
//   - DB_DSN (driver-specific)
//   - DB_HOST
//   - DB_PORT
//   - DB_NAME (or DB_DATABASE)
//   - DB_USER
//   - DB_PASSWORD
//   - DB_SCHEMA
//   - DB_OPTIONS (csv: "k=v,k2=v2")
//   - DB_MIGRATE (default: false)
func databaseParamsFromSource(src cfgcontracts.Source) persistence.DatabaseParams {
	driver := cfgports.LowerTrimmed(src, "DB_DRIVER")
	if driver == "" {
		driver = "mysql"
	}

	database := cfgports.Trimmed(src, "DB_NAME")
	if database == "" {
		database = cfgports.Trimmed(src, "DB_DATABASE")
	}

	params := persistence.DatabaseParams{
		Driver:   driver,
		DSN:      cfgports.Trimmed(src, "DB_DSN"),
		Host:     cfgports.Trimmed(src, "DB_HOST"),
		Port:     cfgports.Int(src, "DB_PORT", 0),
		Database: database,
		Schema:   cfgports.Trimmed(src, "DB_SCHEMA"),
		User:     cfgports.Trimmed(src, "DB_USER"),
		Password: cfgports.Trimmed(src, "DB_PASSWORD"),
		Options:  cfgports.ParseKeyValueCSV(cfgports.Trimmed(src, "DB_OPTIONS")),
		Migrate:  cfgports.Bool(src, "DB_MIGRATE", false),
	}

	return params
}
