package ginmw

import (
	"fmt"
	"net/http"
	"time"

	coretel "github.com/brunojet/my-store-go/app/core/telemetry"
	"github.com/gin-gonic/gin"
)

// Telemetry provides baseline tracing/metrics/logs for every HTTP request.
//
// It is middleware-level "wide" observability: measure everything, then later
// add deeper spans/metrics in specific services/repositories when needed.
func Telemetry(p coretel.Provider) gin.HandlerFunc {
	if p == nil {
		p = coretel.NoopProvider{}
	}
	return func(c *gin.Context) {
		start := time.Now()

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		ctx, span := p.Tracer().Start(
			c.Request.Context(),
			"http "+c.Request.Method+" "+route,
			coretel.Field{Key: "method", Value: c.Request.Method},
			coretel.Field{Key: "route", Value: route},
		)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		status := c.Writer.Status()
		lat := time.Since(start)

		p.Metrics().Inc(
			"http.server.requests",
			1,
			coretel.Field{Key: "method", Value: c.Request.Method},
			coretel.Field{Key: "route", Value: route},
			coretel.Field{Key: "status", Value: status},
		)
		p.Metrics().ObserveDuration(
			"http.server.duration",
			lat,
			coretel.Field{Key: "method", Value: c.Request.Method},
			coretel.Field{Key: "route", Value: route},
			coretel.Field{Key: "status", Value: status},
		)

		var err error
		if status >= http.StatusInternalServerError {
			err = fmt.Errorf("http status %d", status)
		}
		span.End(err, coretel.Field{Key: "status", Value: status}, coretel.Field{Key: "latency", Value: lat.String()})
	}
}
