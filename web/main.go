package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	r.Use(middleware.Logger)
	r.Use(DataDogTracer())

	r.Get("/", welcome())

	http.ListenAndServe(":3000", r)
}

func welcome() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest("GET", "http://web2-server:3001/web2", nil)
		if err != nil {
			log.Printf("Unable to create request: %v", err)
			return
		}
		req.Header = r.Header

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("bad response: %v", err)
			return
		}

		w.Write([]byte("it works"))
	}
}
