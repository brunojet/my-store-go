package observability

import (
	"os"
	"strings"

	"github.com/brunojet/my-store-go/app/infra/observability/adapters/ginmw"
	"github.com/brunojet/my-store-go/app/infra/observability/adapters/httpmw"
	"github.com/brunojet/my-store-go/app/infra/observability/contracts"
	infrotel "github.com/brunojet/my-store-go/app/infra/observability/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5/middleware"
)

type Config = contracts.Config

func defaultConfig() Config {
	return Config{RequestID: true, Telemetry: true, Recovery: true}
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

// ConfigFromEnv builds an observability Config from env.
//
// Env (defaults shown):
//   - OBS_REQUEST_ID=true
//   - OBS_TELEMETRY=true (alias: OBS_ACCESS_LOG)
//   - OBS_RECOVERY=true
func ConfigFromEnv() Config {
	cfg := defaultConfig()
	cfg.RequestID = envBool("OBS_REQUEST_ID", cfg.RequestID)

	// Backward-compatible: OBS_ACCESS_LOG previously controlled request logs.
	// If OBS_TELEMETRY is explicitly set, it wins.
	if _, ok := envString("OBS_TELEMETRY"); ok {
		cfg.Telemetry = envBool("OBS_TELEMETRY", cfg.Telemetry)
	} else {
		cfg.Telemetry = envBool("OBS_ACCESS_LOG", cfg.Telemetry)
	}

	cfg.Recovery = envBool("OBS_RECOVERY", cfg.Recovery)
	return cfg
}

// GinMiddlewares returns the configured Gin middleware chain.
func GinMiddlewares(cfg Config) []contracts.GinMiddleware {
	tel := infrotel.NewStdProvider()

	var mws []contracts.GinMiddleware
	if cfg.RequestID {
		mws = append(mws, ginmw.RequestID())
	}
	if cfg.Telemetry {
		mws = append(mws, ginmw.Telemetry(tel))
	}
	if cfg.Recovery {
		mws = append(mws, gin.Recovery())
	}
	return mws
}

// NetHTTPMiddlewares returns the configured net/http middleware chain.
func NetHTTPMiddlewares(cfg Config) []contracts.NetHTTPMiddleware {
	tel := infrotel.NewStdProvider()

	var mws []contracts.NetHTTPMiddleware
	if cfg.RequestID {
		mws = append(mws, httpmw.RequestID)
	}
	if cfg.Telemetry {
		mws = append(mws, httpmw.Telemetry(tel))
	}
	if cfg.Recovery {
		mws = append(mws, middleware.Recoverer)
	}
	return mws
}
