package handler

import (
	"errors"
	"net/http"

	"github.com/zFlabmonsta/datadog-playground/pkg/log"
)

var ErrorDivisibleZero = errors.New("divide(): cannot be divided by zero")

type Handler struct {
	log *log.LoggerWrapper
}

func (h *Handler) Divide() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := divide(10, 0)
		if err != nil {
			h.log.Errorf(r.Context(), "handler(): cannot get result: %w", err)
		}
		w.Write([]byte("welcome we are dividing by 0 resulting in ERROR"))
	}
}

func divide(a, b float32) (float32, error) {
	if b == 0 {
		return 0, ErrorDivisibleZero
	}

	return a / b, nil
}
