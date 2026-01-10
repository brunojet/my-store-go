package http

import (
	"fmt"
	stdhttp "net/http"
	"os"
	"strings"

	ginadapter "github.com/brunojet/my-store-go/app/infra/http/adapters/gin"
	"github.com/brunojet/my-store-go/app/infra/http/adapters/nethttp"
	"github.com/brunojet/my-store-go/app/infra/http/contracts"
	"github.com/brunojet/my-store-go/app/infra/observability/ginmw"
	"github.com/brunojet/my-store-go/app/infra/observability/httpmw"
	infrotel "github.com/brunojet/my-store-go/app/infra/observability/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5/middleware"
)

// HttpRuntime wires an HTTP framework into a framework-agnostic Router contract.
//
// Router is where modules register routes. Handler is what net/http servers serve.
//
// This package is infra: it owns framework selection/initialization.
type HttpRuntime struct {
	Router  contracts.Router
	Handler stdhttp.Handler
}

type ObservabilityConfig struct {
	RequestID bool
	Telemetry bool
	Recovery  bool
}

type HttpRuntimeConfig struct {
	Observability ObservabilityConfig
}

func defaultRuntimeConfig() HttpRuntimeConfig {
	return HttpRuntimeConfig{
		Observability: ObservabilityConfig{
			RequestID: true,
			Telemetry: true,
			Recovery:  true,
		},
	}
}

func envString(name string) (string, bool) {
	v, ok := os.LookupEnv(name)
	if !ok {
		return "", false
	}
	return strings.TrimSpace(v), true
}

func envBool(name string, def bool) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(name)))
	if v == "" {
		return def
	}
	switch v {
	case "1", "true", "t", "yes", "y", "on", "enabled", "enable":
		return true
	case "0", "false", "f", "no", "n", "off", "disabled", "disable":
		return false
	default:
		return def
	}
}

// ConfigFromEnv builds a RuntimeConfig from env.
//
// Env (defaults shown):
//   - OBS_REQUEST_ID=true
//   - OBS_TELEMETRY=true (alias: OBS_ACCESS_LOG)
//   - OBS_RECOVERY=true
func ConfigFromEnv() HttpRuntimeConfig {
	cfg := defaultRuntimeConfig()
	cfg.Observability.RequestID = envBool("OBS_REQUEST_ID", cfg.Observability.RequestID)

	// Backward-compatible: OBS_ACCESS_LOG previously controlled request logs.
	// If OBS_TELEMETRY is explicitly set, it wins.
	if _, ok := envString("OBS_TELEMETRY"); ok {
		cfg.Observability.Telemetry = envBool("OBS_TELEMETRY", cfg.Observability.Telemetry)
	} else {
		cfg.Observability.Telemetry = envBool("OBS_ACCESS_LOG", cfg.Observability.Telemetry)
	}

	cfg.Observability.Recovery = envBool("OBS_RECOVERY", cfg.Observability.Recovery)
	return cfg
}

// NewGin configures the Gin runtime from the provided config.
func NewGin(cfg HttpRuntimeConfig) *HttpRuntime {
	tel := infrotel.NewStdProvider()
	var mws []gin.HandlerFunc
	if cfg.Observability.RequestID {
		mws = append(mws, ginmw.RequestID())
	}
	if cfg.Observability.Telemetry {
		mws = append(mws, ginmw.Telemetry(tel))
	}
	if cfg.Observability.Recovery {
		mws = append(mws, gin.Recovery())
	}

	router, handler := ginadapter.NewRuntime(ginadapter.WithMiddlewares(mws...))
	return &HttpRuntime{Router: router, Handler: handler}
}

// NewChi configures the chi (net/http) runtime from the provided config.
func NewChi(cfg HttpRuntimeConfig) *HttpRuntime {
	tel := infrotel.NewStdProvider()
	var mws []func(stdhttp.Handler) stdhttp.Handler
	if cfg.Observability.RequestID {
		mws = append(mws, httpmw.RequestID)
	}
	if cfg.Observability.Telemetry {
		mws = append(mws, httpmw.Telemetry(tel))
	}
	if cfg.Observability.Recovery {
		mws = append(mws, middleware.Recoverer)
	}

	router, handler := nethttp.NewRuntime(nethttp.WithMiddlewares(mws...))
	return &HttpRuntime{Router: router, Handler: handler}
}

// SelectFromEnv selects and initializes an HTTP framework based on env.
//
// Env:
//   - HTTP_DRIVER: defaults to "gin". Supported: "gin", "chi".
func SelectFromEnv() (*HttpRuntime, error) {
	cfg := ConfigFromEnv()
	driver := strings.TrimSpace(strings.ToLower(os.Getenv("HTTP_DRIVER")))
	switch driver {
	case "", "gin":
		return NewGin(cfg), nil
	case "chi":
		return NewChi(cfg), nil
	}
	return nil, fmt.Errorf("unsupported HTTP_DRIVER %q", driver)
}
