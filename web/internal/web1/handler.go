package web1

import (
	"context"
	"fmt"
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

func (h *handler) Welcome() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		a := r.URL.Query().Get("a")
		b := r.URL.Query().Get("b")
		req, err := http.NewRequest("GET", fmt.Sprintf("http://web2-server:3001/web2?a=%v&b=%v", a, b), nil)
		if err != nil {
			h.log.Errorf(ctx, "Welcome(): Unable to create request: %w", err)
			return
		}

		req.Header = r.Header
		req.Header.Add("subdomain", "orange")

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			h.log.Errorf(ctx, "Welcome():bad response: %w", err)
			return
		}

		w.Write([]byte("Hello Welcome!"))
	}
}
