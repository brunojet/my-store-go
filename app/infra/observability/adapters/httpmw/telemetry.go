package httpmw

import (
	"fmt"
	"net/http"
	"time"

	coretel "github.com/brunojet/my-store-go/app/core/telemetry"
)

// Telemetry provides baseline tracing/metrics/logs for every HTTP request.
func Telemetry(p coretel.Provider) func(http.Handler) http.Handler {
	if p == nil {
		p = coretel.NoopProvider{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()

			ctx, span := p.Tracer().Start(
				r.Context(),
				"http "+r.Method+" "+r.URL.Path,
				coretel.Field{Key: "method", Value: r.Method},
				coretel.Field{Key: "path", Value: r.URL.Path},
			)

			next.ServeHTTP(sr, r.WithContext(ctx))

			lat := time.Since(start)
			status := sr.status

			p.Metrics().Inc(
				"http.server.requests",
				1,
				coretel.Field{Key: "method", Value: r.Method},
				coretel.Field{Key: "path", Value: r.URL.Path},
				coretel.Field{Key: "status", Value: status},
			)
			p.Metrics().ObserveDuration(
				"http.server.duration",
				lat,
				coretel.Field{Key: "method", Value: r.Method},
				coretel.Field{Key: "path", Value: r.URL.Path},
				coretel.Field{Key: "status", Value: status},
			)

			var err error
			if status >= http.StatusInternalServerError {
				err = fmt.Errorf("http status %d", status)
			}
			span.End(err, coretel.Field{Key: "status", Value: status}, coretel.Field{Key: "latency", Value: lat.String()})
		})
	}
}
