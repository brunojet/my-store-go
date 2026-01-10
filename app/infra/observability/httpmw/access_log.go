package httpmw

import (
	"log"
	"net/http"
	"time"

	"github.com/brunojet/my-store-go/app/infra/observability/requestid"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

// AccessLog writes a minimal log line per request.
func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(sr, r)

		lat := time.Since(start)
		rid, _ := requestid.From(r.Context())

		log.Printf("http request rid=%s method=%s path=%s status=%d latency=%s", rid, r.Method, r.URL.Path, sr.status, lat)
	})
}
