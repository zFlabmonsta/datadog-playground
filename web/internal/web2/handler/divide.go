package handler

import (
	"context"
	"errors"
	"net/http"
)

var ErrorDivisibleZero = errors.New("divide(): cannot be divided by zero")

type LoggerWrapper interface {
	Errorf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
}

type Handler struct {
	log LoggerWrapper
}

func NewHandler(log LoggerWrapper) *Handler {
	return &Handler{
		log: log,
	}
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
