package observability

import (
	"os"
	"strings"

	"github.com/brunojet/my-store-go/app/infra/observability/contracts"
)

// HTTPSelection is the env/config-driven selection of HTTP driver + observability config.
//
// It exists to keep the composition root thin while avoiding a dependency from
// infra/observability -> infra/http.
type HTTPSelection struct {
	Driver string
	Config Config
}

// HTTPDriverFromEnv returns the normalized HTTP driver selection from env.
//
// Env:
//   - HTTP_DRIVER: defaults to "gin" when empty.
func HTTPDriverFromEnv() string {
	return strings.TrimSpace(strings.ToLower(os.Getenv("HTTP_DRIVER")))
}

// SelectFromEnv selects the HTTP driver and observability config from env.
func SelectFromEnv() HTTPSelection {
	return HTTPSelection{
		Driver: HTTPDriverFromEnv(),
		Config: ConfigFromEnv(),
	}
}

// GetHttpMiddlewares returns a single middleware package for the selected driver.
//
// Exactly one of the returned slices is expected to be populated.
func (s HTTPSelection) GetHttpMiddlewares() contracts.HTTPMiddlewares {
	return HTTPMiddlewares(s.Driver, s.Config)
}

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

// HTTPMiddlewaresFromEnv reads observability config from env and builds middleware chain
// for the given driver.
func HTTPMiddlewaresFromEnv(driver string) contracts.HTTPMiddlewares {
	return HTTPMiddlewares(driver, ConfigFromEnv())
}

// HTTPMiddlewaresFromHTTPDriverEnv reads HTTP_DRIVER and observability config from env
// and builds middleware chain for the selected driver.
func HTTPMiddlewaresFromHTTPDriverEnv() contracts.HTTPMiddlewares {
	return HTTPMiddlewares(HTTPDriverFromEnv(), ConfigFromEnv())
}
