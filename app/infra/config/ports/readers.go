package ports

import (
	"strconv"
	"strings"
	"time"
)

// Trimmed returns the trimmed value for key, or empty string if not set.
func Trimmed(src Source, key string) string {
	v, _ := src.Lookup(key)
	return strings.TrimSpace(v)
}

// LowerTrimmed returns the trimmed, lowercased value for key, or empty string if not set.
func LowerTrimmed(src Source, key string) string {
	return strings.ToLower(Trimmed(src, key))
}

// Bool returns a boolean parsed from key with a default.
func Bool(src Source, key string, def bool) bool {
	v, ok := src.Lookup(key)
	if !ok {
		return def
	}
	vv := strings.TrimSpace(strings.ToLower(v))
	if vv == "" {
		return def
	}
	switch vv {
	case "1", "true", "t", "yes", "y", "on", "enabled", "enable":
		return true
	case "0", "false", "f", "no", "n", "off", "disabled", "disable":
		return false
	default:
		return def
	}
}

// Duration returns a duration parsed from key with a default.
//
// Supported formats:
//   - Go duration (e.g. "250ms", "5s", "1m")
//   - seconds as int (e.g. "5")
func Duration(src Source, key string, def time.Duration) time.Duration {
	v, ok := src.Lookup(key)
	if !ok {
		return def
	}
	vv := strings.TrimSpace(v)
	if vv == "" {
		return def
	}
	if d, err := time.ParseDuration(vv); err == nil {
		return d
	}
	if n, err := strconv.Atoi(vv); err == nil {
		return time.Duration(n) * time.Second
	}
	return def
}
