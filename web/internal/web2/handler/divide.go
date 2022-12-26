package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zFlabmonsta/datadog-playground/internal/web2/math"
)

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
		a := stringToFloat32(r.URL.Query().Get("a"))
		b := stringToFloat32(r.URL.Query().Get("b"))

		result, err := math.Divide(float32(a), float32(b))
		if isDivisible(err) {
			h.log.Errorf(r.Context(), "handler(): cannot get result: %w", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		answer := fmt.Sprintf("Answer: %f", result)
		w.Write([]byte(answer))
	}
}

func stringToFloat32(s string) float32 {
	n, _ := strconv.ParseFloat(s, 32)
	return float32(n)
}

func isDivisible(err error) bool {
	return errors.Is(err, math.ErrorDivisibleZero)
}
