package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	logger := log.New()
	logger.Formatter = &log.JSONFormatter{}
	logger.Out = os.Stdout

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

	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger: logger,
	})

	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Use(DataDogTracer())

	r.Get("/web2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome web2"))
	})

	http.ListenAndServe(":3001", r)
}
