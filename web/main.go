package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	mwtracer "github.com/zFlabmonsta/datadog-playground/pkg/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
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
		logger.Errorf("error starting profiler: %v", err)
	}
	defer profiler.Stop()

	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger: logger,
	})

	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Use(mwtracer.DataDogTracer())

	r.Get("/", welcome())

	http.ListenAndServe(":3000", r)
}

func welcome() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest("GET", "http://web2-server:3001/web2", nil)
		if err != nil {
			log.Errorf("Unable to create request: %v", err)
			return
		}
		req.Header = r.Header

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Errorf("bad response: %v", err)
			return
		}

		w.Write([]byte("it works"))
	}
}
