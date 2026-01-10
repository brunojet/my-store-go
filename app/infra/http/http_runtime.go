package http

import (
	"fmt"
	stdhttp "net/http"
	"os"
	"strings"

	ginadapter "github.com/brunojet/my-store-go/app/infra/http/adapters/gin"
	"github.com/brunojet/my-store-go/app/infra/http/adapters/nethttp"
	"github.com/brunojet/my-store-go/app/infra/http/contracts"
	"github.com/gin-gonic/gin"
)

type Option func(*runtimeOptions)

type runtimeOptions struct {
	ginMiddlewares     []gin.HandlerFunc
	netHTTPMiddlewares []func(stdhttp.Handler) stdhttp.Handler
}

func defaultRuntimeOptions() runtimeOptions {
	return runtimeOptions{}
}

// WithGinMiddlewares appends custom Gin middlewares.
func WithGinMiddlewares(mw ...gin.HandlerFunc) Option {
	return func(o *runtimeOptions) {
		o.ginMiddlewares = append(o.ginMiddlewares, mw...)
	}
}

// WithNetHTTPMiddlewares appends custom net/http middlewares.
func WithNetHTTPMiddlewares(mw ...func(stdhttp.Handler) stdhttp.Handler) Option {
	return func(o *runtimeOptions) {
		o.netHTTPMiddlewares = append(o.netHTTPMiddlewares, mw...)
	}
}

func buildRuntimeOptions(opts ...Option) runtimeOptions {
	o := defaultRuntimeOptions()
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// HttpRuntime wires an HTTP framework into a framework-agnostic Router contract.
//
// Router is where modules register routes. Handler is what net/http servers serve.
//
// This package is infra: it owns framework selection/initialization.
type HttpRuntime struct {
	Router  contracts.Router
	Handler stdhttp.Handler
}

// NewGin builds the Gin runtime.
func NewGin(opts ...Option) *HttpRuntime {
	o := buildRuntimeOptions(opts...)
	router, handler := ginadapter.NewRuntime(ginadapter.WithMiddlewares(o.ginMiddlewares...))
	return &HttpRuntime{Router: router, Handler: handler}
}

// NewChi builds the chi (net/http) runtime.
func NewChi(opts ...Option) *HttpRuntime {
	o := buildRuntimeOptions(opts...)
	router, handler := nethttp.NewRuntime(nethttp.WithMiddlewares(o.netHTTPMiddlewares...))
	return &HttpRuntime{Router: router, Handler: handler}
}

// SelectFromEnv selects and initializes an HTTP framework based on env.
//
// Env:
//   - HTTP_DRIVER: defaults to "gin". Supported: "gin", "chi".
func SelectFromEnv(opts ...Option) (*HttpRuntime, error) {
	driver := strings.TrimSpace(strings.ToLower(os.Getenv("HTTP_DRIVER")))
	switch driver {
	case "", "gin":
		return NewGin(opts...), nil
	case "chi":
		return NewChi(opts...), nil
	}
	return nil, fmt.Errorf("unsupported HTTP_DRIVER %q", driver)
}
