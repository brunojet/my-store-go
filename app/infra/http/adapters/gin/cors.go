package base

import (
	"net/http"

	httpcontracts "github.com/brunojet/my-store-go/app/infra/http/contracts"
	httpports "github.com/brunojet/my-store-go/app/infra/http/ports"
	"github.com/gin-gonic/gin"
)

// CORS returns a Gin middleware that applies CORS headers and handles preflight OPTIONS.
func CORS(cfg httpcontracts.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if h, ok := httpports.BuildCORSHeaders(cfg, origin); ok {
			c.Header("Access-Control-Allow-Origin", h.AllowOrigin)
			if h.VaryOrigin {
				c.Header("Vary", "Origin")
			}
			c.Header("Access-Control-Allow-Methods", h.AllowMethods)
			c.Header("Access-Control-Allow-Headers", h.AllowHeaders)
			if h.ExposeHeaders != "" {
				c.Header("Access-Control-Expose-Headers", h.ExposeHeaders)
			}
			if h.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			if h.MaxAgeSeconds > 0 {
				c.Header("Access-Control-Max-Age", h.MaxAgeString())
			}
		}

		if c.Request.Method == http.MethodOptions {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}
		c.Next()
	}
}
