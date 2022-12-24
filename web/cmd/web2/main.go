package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/zFlabmonsta/datadog-playground/internal/web2"
	"github.com/zFlabmonsta/datadog-playground/pkg/log"
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
		logrus.Fatalf("error starting profiler: %v", err)
	}
	defer profiler.Stop()

	r := chi.NewRouter()
	r.Use(web2.DataDogTracer())
	r.Use(web2.CompanyContext())

	logger := log.NewLogger(logrus.New(), &logrus.JSONFormatter{})
	r.Use(log.Logger(logger))

	r.Get("/web2", func(w http.ResponseWriter, r *http.Request) {
		_, err := web2.Divide(10, 0)
		if err != nil {
			logger.Errorf(r.Context(), "handler(): cannot get result: %w", err)
		}
		w.Write([]byte("welcome we are dividing by 0 resulting in ERROR"))
	})

	http.ListenAndServe(":3001", r)
}
