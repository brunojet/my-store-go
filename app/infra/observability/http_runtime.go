package observability

import (
	"strings"

	"github.com/brunojet/my-store-go/app/infra/observability/contracts"
)

// HTTPMiddlewares returns the configured middleware chain for the given driver.
//
// driver is expected to be normalized (trimmed/lowercased). Unknown/empty defaults to "gin".
func HTTPMiddlewares(driver string, cfg Config) contracts.HTTPMiddlewares {
	driver = strings.TrimSpace(strings.ToLower(driver))

	switch driver {
	case "chi":
		return contracts.HTTPMiddlewares{NetHTTP: NetHTTPMiddlewares(cfg)}
	default:
		return contracts.HTTPMiddlewares{Gin: GinMiddlewares(cfg)}
	}
}
