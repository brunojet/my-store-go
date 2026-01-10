package observability

import (
	"github.com/brunojet/my-store-go/app/infra/observability/adapters/ginmw"
	"github.com/brunojet/my-store-go/app/infra/observability/adapters/httpmw"
	"github.com/brunojet/my-store-go/app/infra/observability/contracts"
	infrotel "github.com/brunojet/my-store-go/app/infra/observability/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5/middleware"
)

type Config = contracts.Config

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
