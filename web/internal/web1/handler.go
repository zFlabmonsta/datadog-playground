package web1

import (
	"net/http"

	pkg "github.com/zFlabmonsta/datadog-playground/pkg/log"
)

type handler struct {
	log *pkg.LoggerWrapper
}

func NewHandler(logger *pkg.LoggerWrapper) *handler {
	return &handler{log: logger}
}

func (h *handler) Welcome() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req, err := http.NewRequest("GET", "http://web2-server:3001/web2", nil)
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

		w.Write([]byte("it works"))
	}
}
