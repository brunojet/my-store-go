package types

import "time"

// CORSConfig configures Cross-Origin Resource Sharing (CORS) behavior.
//
// This is a contract between the infra/http wiring layer and the HTTP framework
// adapters (gin, net/http).
type CORSConfig struct {
	Enabled          bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}
