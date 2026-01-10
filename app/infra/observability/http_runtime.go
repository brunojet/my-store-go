package observability

import (
	"os"
	"strings"

	httpruntime "github.com/brunojet/my-store-go/app/infra/http"
)

// HTTPRuntimeOptions builds httpruntime options (middlewares) for the given driver.
//
// driver is expected to be normalized (trimmed/lowercased). Unknown/empty defaults to "gin".
func HTTPRuntimeOptions(driver string, cfg Config) []httpruntime.Option {
	driver = strings.TrimSpace(strings.ToLower(driver))

	switch driver {
	case "chi":
		return []httpruntime.Option{
			httpruntime.WithNetHTTPMiddlewares(NetHTTPMiddlewares(cfg)...),
		}
	default:
		return []httpruntime.Option{
			httpruntime.WithGinMiddlewares(GinMiddlewares(cfg)...),
		}
	}
}

// HTTPRuntimeOptionsFromEnv is a convenience helper that reads the observability config
// from env and builds the runtime options for the given driver.
func HTTPRuntimeOptionsFromEnv(driver string) []httpruntime.Option {
	return HTTPRuntimeOptions(driver, ConfigFromEnv())
}

// HTTPDriverFromEnv returns the normalized HTTP driver selection from env.
//
// Env:
//   - HTTP_DRIVER: defaults to "gin" when empty.
func HTTPDriverFromEnv() string {
	return strings.TrimSpace(strings.ToLower(os.Getenv("HTTP_DRIVER")))
}

// SelectFromEnv builds runtime options using HTTP_DRIVER and observability env.
func SelectFromEnv() []httpruntime.Option {
	return HTTPRuntimeOptions(HTTPDriverFromEnv(), ConfigFromEnv())
}
