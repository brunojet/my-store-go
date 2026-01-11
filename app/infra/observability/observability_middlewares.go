package observability

import (
	httptypes "github.com/brunojet/my-store-go/app/infra/http/types"
	"github.com/brunojet/my-store-go/app/infra/observability/adapters/ginmw"
	"github.com/brunojet/my-store-go/app/infra/observability/adapters/httpmw"
	infrotel "github.com/brunojet/my-store-go/app/infra/observability/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5/middleware"
)

// ginMiddlewares returns the configured Gin middleware chain.
func ginMiddlewares(cfg ObservabilityConfig) []GinMiddleware {
	tel := infrotel.NewStdProvider()

	var mws []GinMiddleware
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
func netHTTPMiddlewares(cfg ObservabilityConfig) []NetHTTPMiddleware {
	tel := infrotel.NewStdProvider()

	var mws []NetHTTPMiddleware
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

func BuildMiddlewares(driver string, cfg ObservabilityConfig) ObservabilityMiddlewares {
	d := httptypes.NormalizeHTTPDriver(driver)

	switch d {
	case httptypes.HTTPDriverChi:
		return ObservabilityMiddlewares{NetHTTP: netHTTPMiddlewares(cfg)}
	default:
		return ObservabilityMiddlewares{Gin: ginMiddlewares(cfg)}
	}
}
