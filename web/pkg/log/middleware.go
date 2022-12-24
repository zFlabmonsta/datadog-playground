package log

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

func Logger(l *LoggerWrapper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			timeNow := time.Now()
			defer func() {
				fields := Fields{
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
