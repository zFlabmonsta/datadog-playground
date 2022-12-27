package web1

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type LoggerWrapper interface {
	Errorf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
}

type handler struct {
	log LoggerWrapper
}

func NewHandler(logger LoggerWrapper) *handler {
	return &handler{log: logger}
}

func (h *handler) Divide() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		a := r.URL.Query().Get("a")
		b := r.URL.Query().Get("b")
		req, err := http.NewRequest("GET", fmt.Sprintf("http://web2-server:3001/web2?a=%v&b=%v", a, b), nil)
		if err != nil {
			h.log.Errorf(ctx, "Divide(): Unable to create request: %w", err)
			return
		}

		req.Header = r.Header
		req.Header.Add("subdomain", "orange")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			h.log.Errorf(ctx, "Divide(): Client was not able to make call %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			h.log.Errorf(ctx, "Divide(): Was not able to read response body %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if IsBadRequest(resp) {
			w.WriteHeader(resp.StatusCode)
			w.Write(body)
			return
		}

		w.Write([]byte(fmt.Sprintf("Answer: %v", string(body))))
	}
}
