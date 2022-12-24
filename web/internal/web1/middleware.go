package web1

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func DataDogTracer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			span, _ := tracer.StartSpanFromContext(r.Context(), "web.request", tracer.ResourceName(r.RequestURI))
			defer span.Finish()

			err := tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(r.Header))
			if err != nil {
				log.Printf("Failed to inject tracer: %v", err)
			}

			r = r.WithContext(context.WithValue(r.Context(), "datadogTraceID", r.Header.Get("X-Datadog-Trace-Id")))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func InjectCompanyContext() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r.Header.Add("subdomain", "orange")
			r = r.WithContext(context.WithValue(r.Context(), "subdomain", r.Header.Get("subdomain")))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
