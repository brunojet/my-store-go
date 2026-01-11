package httpmw

import (
	"fmt"
	"net/http"
	"time"

	coretel "github.com/brunojet/my-store-go/app/core/telemetry"
	"github.com/go-chi/chi/v5"
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

			route := chi.RouteContext(r.Context()).RoutePattern()
			if route == "" {
				route = r.URL.Path
			}

			ctx, span := p.Tracer().Start(
				r.Context(),
				"http "+r.Method+" "+route,
				coretel.Field{Key: "method", Value: r.Method},
				coretel.Field{Key: "route", Value: route},
			)

			next.ServeHTTP(sr, r.WithContext(ctx))

			lat := time.Since(start)
			status := sr.status

			p.Metrics().Inc(
				"http.server.requests",
				1,
				coretel.Field{Key: "method", Value: r.Method},
				coretel.Field{Key: "route", Value: route},
				coretel.Field{Key: "status", Value: status},
			)
			p.Metrics().ObserveDuration(
				"http.server.duration",
				lat,
				coretel.Field{Key: "method", Value: r.Method},
				coretel.Field{Key: "route", Value: route},
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
