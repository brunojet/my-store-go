package observability

import (
	httptypes "github.com/brunojet/my-store-go/app/infra/http/types"
	"github.com/brunojet/my-store-go/app/infra/observability/adapters/ginmw"
	"github.com/brunojet/my-store-go/app/infra/observability/adapters/httpmw"
	infrotel "github.com/brunojet/my-store-go/app/infra/observability/telemetry"
	"github.com/brunojet/my-store-go/app/infra/observability/types"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5/middleware"
)

// ginMiddlewares returns the configured Gin middleware chain.
func ginMiddlewares(cfg types.MiddlewareConfig) []types.GinMiddleware {
	tel := infrotel.NewStdProvider()

	var mws []types.GinMiddleware
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

// netHTTPMiddlewares returns the configured net/http middleware chain.
func netHTTPMiddlewares(cfg types.MiddlewareConfig) []types.NetHTTPMiddleware {
	tel := infrotel.NewStdProvider()

	var mws []types.NetHTTPMiddleware
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

func BuildMiddlewares(driver httptypes.HTTPDriver, cfg types.MiddlewareConfig) types.Middlewares {
	switch driver {
	case httptypes.HTTPDriverChi:
		return types.Middlewares{NetHTTP: netHTTPMiddlewares(cfg)}
	default:
		return types.Middlewares{Gin: ginMiddlewares(cfg)}
	}
}
