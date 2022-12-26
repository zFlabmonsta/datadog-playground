package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/zFlabmonsta/datadog-playground/internal/web1"
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
	r.Use(web1.DataDogTracer())
	r.Use(web1.InjectCompanyContext())

	logger := log.NewLogger(logrus.New(), &logrus.JSONFormatter{})
	r.Use(log.Logger(logger))

	handler := web1.NewHandler(logger)
	r.Get("/divide", handler.Welcome())

	http.ListenAndServe(":3000", r)
}
