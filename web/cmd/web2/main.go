package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/zFlabmonsta/datadog-playground/internal/web2"
	"github.com/zFlabmonsta/datadog-playground/internal/web2/handler"
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

	handler := handler.NewHandler(logger)
	r.Get("/web2", handler.Divide())

	http.ListenAndServe(":3001", r)
}
