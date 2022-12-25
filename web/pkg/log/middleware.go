package log

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
)

func Logger(l *LoggerWrapper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				http := &HTTP{
					StatusCode: ww.Status(),
					Method:     r.Method,
				}
				r = r.WithContext(context.WithValue(r.Context(), HTTP{}, http))
				l.Infof(r.Context(), "http://%s%s", r.Host, r.RequestURI)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
