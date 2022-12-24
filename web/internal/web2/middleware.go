package web2

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/zFlabmonsta/datadog-playground/pkg"
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
			r = r.WithContext(context.WithValue(r.Context(), "subdomain", "orange"))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func Logger(l *pkg.LoggerWrapper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			timeNow := time.Now()
			defer func() {
				fields := pkg.Fields{
					"duration":    timeNow.String(),
					"status_code": ww.Status(),
					"method":      r.Method,
					"service":     "web-2",
				}
				l.WithFields(fields).Infof(r.Context(), "http://%s%s", r.Host, r.RequestURI)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
