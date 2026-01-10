package ginmw

import (
	"log"
	"time"

	"github.com/brunojet/my-store-go/app/infra/observability/requestid"
	"github.com/gin-gonic/gin"
)

// AccessLog writes a minimal structured-ish log line per request.
//
// It is intentionally tiny and dependency-free (uses stdlib log).
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		lat := time.Since(start)
		rid, _ := requestid.From(c.Request.Context())
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		log.Printf("http request rid=%s method=%s path=%s status=%d latency=%s", rid, method, path, status, lat)
	}
}
