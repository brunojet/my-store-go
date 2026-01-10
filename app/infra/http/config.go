package http

import (
	"strings"
	"time"

	"github.com/brunojet/my-store-go/app/infra/config/adapters/env"
	"github.com/brunojet/my-store-go/app/infra/config/ports"
	httpcontracts "github.com/brunojet/my-store-go/app/infra/http/contracts"
)

// Config holds HTTP runtime selection config.
//
// Env:
//   - HTTP_DRIVER (default: "gin")
//   - CORS_ENABLED (default: true)
//   - CORS_ALLOW_ORIGINS (default: "*")
//   - CORS_ALLOW_METHODS (default: "GET,POST,PATCH,DELETE,OPTIONS")
//   - CORS_ALLOW_HEADERS (default: "Content-Type,Authorization")
//   - CORS_EXPOSE_HEADERS (default: "")
//   - CORS_ALLOW_CREDENTIALS (default: false)
//   - CORS_MAX_AGE (default: 10m)
type Config struct {
	Driver string
	CORS   httpcontracts.CORSConfig
}

func defaultCORSConfig() httpcontracts.CORSConfig {
	return httpcontracts.CORSConfig{
		Enabled:          true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    nil,
		AllowCredentials: false,
		MaxAge:           10 * time.Minute,
	}
}

func LoadConfig(src ports.Source) Config {
	cors := defaultCORSConfig()
	cors.Enabled = ports.Bool(src, "CORS_ENABLED", cors.Enabled)
	cors.AllowCredentials = ports.Bool(src, "CORS_ALLOW_CREDENTIALS", cors.AllowCredentials)
	cors.MaxAge = ports.Duration(src, "CORS_MAX_AGE", cors.MaxAge)

	cors.AllowOrigins = splitCSV(ports.Trimmed(src, "CORS_ALLOW_ORIGINS"), cors.AllowOrigins)
	cors.AllowMethods = splitCSVUpper(ports.Trimmed(src, "CORS_ALLOW_METHODS"), cors.AllowMethods)
	cors.AllowHeaders = splitCSV(ports.Trimmed(src, "CORS_ALLOW_HEADERS"), cors.AllowHeaders)
	cors.ExposeHeaders = splitCSV(ports.Trimmed(src, "CORS_EXPOSE_HEADERS"), cors.ExposeHeaders)

	driver := ports.LowerTrimmed(src, "HTTP_DRIVER")
	if driver == "" {
		driver = "gin"
	}
	return Config{Driver: driver, CORS: cors}
}

func ConfigFromEnv() Config {
	return LoadConfig(env.New())
}

func splitCSV(raw string, def []string) []string {
	if strings.TrimSpace(raw) == "" {
		return def
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	if len(out) == 0 {
		return def
	}
	return out
}

func splitCSVUpper(raw string, def []string) []string {
	items := splitCSV(raw, def)
	for i := range items {
		items[i] = strings.ToUpper(items[i])
	}
	return items
}
