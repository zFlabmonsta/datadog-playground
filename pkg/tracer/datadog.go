package tracer

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func DataDogTracer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header))
			if err != nil {
				log.Errorf("%v", err)
				return
			}

			span, _ := tracer.StartSpanFromContext(r.Context(), "web.request", tracer.ResourceName(r.RequestURI), tracer.ChildOf(sctx))
			defer span.Finish()

			err = tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(r.Header))
			if err != nil {
				log.Errorf("%v", err)
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
