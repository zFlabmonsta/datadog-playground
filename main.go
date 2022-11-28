package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	tracer.Start(
		tracer.WithService("dd-toy"),
		tracer.WithEnv("development"),
	)
	defer tracer.Stop()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(DataDogTracer())
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":3000", r)
}

func DataDogTracer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			span := tracer.StartSpan("web.request", tracer.ResourceName(r.RequestURI))
			defer span.Finish()

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
