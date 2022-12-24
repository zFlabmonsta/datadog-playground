package web2

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func DataDogTracer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header))
			if err != nil {
				log.Printf("no span context %v", err)
				return
			}

			span := tracer.StartSpan("web.request", tracer.ChildOf(sctx))
			defer span.Finish()

			r = r.WithContext(context.WithValue(r.Context(), "datadogTraceID", r.Header.Get("X-Datadog-Trace-Id")))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func CompanyContext() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "subdomain", r.Header.Get("subdomain")))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
