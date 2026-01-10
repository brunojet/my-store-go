package observability

import (
	"os"
	"strings"

	"github.com/brunojet/my-store-go/app/infra/observability/contracts"
)

// HTTPDriverFromEnv returns the normalized HTTP driver selection from env.
//
// Env:
//   - HTTP_DRIVER: defaults to "gin" when empty.
func HTTPDriverFromEnv() string {
	return strings.TrimSpace(strings.ToLower(os.Getenv("HTTP_DRIVER")))
}

// HTTPMiddlewares returns the configured middleware chain for the given driver.
//
// driver is expected to be normalized (trimmed/lowercased). Unknown/empty defaults to "gin".
func HTTPMiddlewares(driver string, cfg Config) (gin []contracts.GinMiddleware, nethttp []contracts.NetHTTPMiddleware) {
	driver = strings.TrimSpace(strings.ToLower(driver))

	switch driver {
	case "chi":
		return nil, NetHTTPMiddlewares(cfg)
	default:
		return GinMiddlewares(cfg), nil
	}
}

// HTTPMiddlewaresFromEnv reads observability config from env and builds middleware chain
// for the given driver.
func HTTPMiddlewaresFromEnv(driver string) (gin []contracts.GinMiddleware, nethttp []contracts.NetHTTPMiddleware) {
	return HTTPMiddlewares(driver, ConfigFromEnv())
}

// HTTPMiddlewaresFromHTTPDriverEnv reads HTTP_DRIVER and observability config from env
// and builds middleware chain for the selected driver.
func HTTPMiddlewaresFromHTTPDriverEnv() (gin []contracts.GinMiddleware, nethttp []contracts.NetHTTPMiddleware) {
	return HTTPMiddlewares(HTTPDriverFromEnv(), ConfigFromEnv())
}
