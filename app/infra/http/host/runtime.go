package host

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/brunojet/my-store-go/app/infra/http/contracts"
	ginadapter "github.com/brunojet/my-store-go/app/infra/http/gin"
	"github.com/brunojet/my-store-go/app/infra/http/nethttp"
	"github.com/brunojet/my-store-go/app/infra/observability/ginmw"
	"github.com/brunojet/my-store-go/app/infra/observability/httpmw"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5/middleware"
)

// Runtime wires an HTTP framework into a framework-agnostic Router contract.
//
// Router is where modules register routes. Handler is what net/http servers serve.
//
// This package is infra: it owns framework selection/initialization.
type Runtime struct {
	Router  contracts.Router
	Handler http.Handler
}

type ObservabilityConfig struct {
	RequestID bool
	AccessLog bool
	Recovery  bool
}

type RuntimeConfig struct {
	Observability ObservabilityConfig
}

func defaultRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		Observability: ObservabilityConfig{
			RequestID: true,
			AccessLog: true,
			Recovery:  true,
		},
	}
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
//   - OBS_ACCESS_LOG=true
//   - OBS_RECOVERY=true
func ConfigFromEnv() RuntimeConfig {
	cfg := defaultRuntimeConfig()
	cfg.Observability.RequestID = envBool("OBS_REQUEST_ID", cfg.Observability.RequestID)
	cfg.Observability.AccessLog = envBool("OBS_ACCESS_LOG", cfg.Observability.AccessLog)
	cfg.Observability.Recovery = envBool("OBS_RECOVERY", cfg.Observability.Recovery)
	return cfg
}

// NewGin configures the Gin runtime from the provided config.
func NewGin(cfg RuntimeConfig) *Runtime {
	var mws []gin.HandlerFunc
	if cfg.Observability.RequestID {
		mws = append(mws, ginmw.RequestID())
	}
	if cfg.Observability.AccessLog {
		mws = append(mws, ginmw.AccessLog())
	}
	if cfg.Observability.Recovery {
		mws = append(mws, gin.Recovery())
	}

	router, handler := ginadapter.NewRuntime(ginadapter.WithMiddlewares(mws...))
	return &Runtime{
		Router:  router,
		Handler: handler,
	}
}

// NewChi configures the chi (net/http) runtime from the provided config.
func NewChi(cfg RuntimeConfig) *Runtime {
	var mws []func(http.Handler) http.Handler
	if cfg.Observability.RequestID {
		mws = append(mws, httpmw.RequestID)
	}
	if cfg.Observability.AccessLog {
		mws = append(mws, httpmw.AccessLog)
	}
	if cfg.Observability.Recovery {
		mws = append(mws, middleware.Recoverer)
	}

	router, handler := nethttp.NewRuntime(nethttp.WithMiddlewares(mws...))
	return &Runtime{
		Router:  router,
		Handler: handler,
	}
}

// SelectFromEnv selects and initializes an HTTP framework based on env.
//
// Env:
//   - HTTP_DRIVER: defaults to "gin". Supported: "gin", "chi".
func SelectFromEnv() (*Runtime, error) {
	cfg := ConfigFromEnv()
	driver := strings.TrimSpace(strings.ToLower(os.Getenv("HTTP_DRIVER")))
	if driver == "" || driver == "gin" {
		return NewGin(cfg), nil
	}
	if driver == "chi" {
		return NewChi(cfg), nil
	}
	return nil, fmt.Errorf("unsupported HTTP_DRIVER %q", driver)
}
