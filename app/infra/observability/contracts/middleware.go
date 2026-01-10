package contracts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Config defines which observability features are enabled.
// It is shared across all HTTP adapters.
//
// Defaults are applied by infra/observability.
type Config struct {
	RequestID bool
	Telemetry bool
	Recovery  bool
}

// GinMiddleware is a Gin middleware.
type GinMiddleware = gin.HandlerFunc

// NetHTTPMiddleware is a net/http middleware.
type NetHTTPMiddleware = func(http.Handler) http.Handler

// HTTPMiddlewares is a single "package" of middlewares for the selected HTTP driver.
// Exactly one of Gin or NetHTTP should be populated.
type HTTPMiddlewares struct {
	Gin     []GinMiddleware
	NetHTTP []NetHTTPMiddleware
}
