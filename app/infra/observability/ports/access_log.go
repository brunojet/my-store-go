package ports

import (
	"fmt"
	"strings"
	"time"
)

// SelectRoute prefers a low-cardinality route template when available.
//
// For example:
//   - template: "/apps/:id" (Gin) or "/apps/{id}" (Chi)
//   - actual:   "/apps/123"
//
// When template is empty, it falls back to actual.
func SelectRoute(template string, actual string) string {
	if strings.TrimSpace(template) != "" {
		return template
	}
	return actual
}

// FormatAccessLog returns the canonical access log line used by infra.
func FormatAccessLog(rid string, method string, route string, status int, latency time.Duration) string {
	return fmt.Sprintf(
		"http request rid=%s method=%s path=%s status=%d latency=%s",
		rid,
		method,
		route,
		status,
		latency,
	)
}
