package main

import (
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/zFlabmonsta/datadog-playground/internal/web2"
	"github.com/zFlabmonsta/datadog-playground/pkg"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	tracer.Start()
	defer tracer.Stop()

	err := profiler.Start(
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	)
	if err != nil {
		log.Fatalf("error starting profiler: %v", err)
	}
	defer profiler.Stop()

	r := chi.NewRouter()
	r.Use(web2.DataDogTracer())
	r.Use(web2.CompanyContext())

	logger := pkg.NewLogger(log.New(), &log.JSONFormatter{})
	r.Use(web2.Logger(logger))

	r.Get("/web2", func(w http.ResponseWriter, r *http.Request) {
		_, err := web2.Divide(10, 0)
		if err != nil {
			logger.Errorf(r.Context(), "handler(): cannot get result: %w", err)
		}
		w.Write([]byte("welcome we are dividing by 0 resulting in ERROR"))
	})

	http.ListenAndServe(":3001", r)
}
