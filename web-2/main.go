package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
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

	log.SetFormatter(&log.JSONFormatter{})
	standardFields := log.Fields{
		"hostname": "localhost",
		"port":     "3001",
		"service":  "datadog-playground2",
	}

	log.WithFields(standardFields)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(DataDogTracer())
	r.Use(CompanyContext())

	r.Get("/web2", func(w http.ResponseWriter, r *http.Request) {
		_, err := divide(10, 0)
		if err != nil {
			Errorf(r.Context(), "handler(): cannot get result: %w", err)
		}
		w.Write([]byte("welcome we are dividing by 0 resulting in ERROR"))
	})

	http.ListenAndServe(":3001", r)
}
