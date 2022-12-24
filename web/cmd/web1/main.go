package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/zFlabmonsta/datadog-playground/internal/web1"
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
	r.Use(web1.DataDogTracer())

	logger := pkg.NewLogger(log.New(), &log.JSONFormatter{})
	r.Use(web2.Logger(logger))

	r.Get("/", web1.Welcome())

	http.ListenAndServe(":3000", r)
}
