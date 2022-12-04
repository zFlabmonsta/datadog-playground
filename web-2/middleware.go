package main

import (
	"log"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func DataDogTracer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Printf("r.Header: %v", r.Header)
			sctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header))
			if err != nil {
				log.Printf("no span context %v", err)
				return
			}

			span := tracer.StartSpan("web.request", tracer.ChildOf(sctx))
			defer span.Finish()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
