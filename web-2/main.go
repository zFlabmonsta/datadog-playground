package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	r.Get("/web2", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("header: %v", r.Header)
		w.Write([]byte("welcome web2"))
	})

	http.ListenAndServe(":3001", r)
}
