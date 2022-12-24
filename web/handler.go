package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func welcome() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest("GET", "http://web2-server:3001/web2", nil)
		if err != nil {
			log.Printf("Unable to create request: %v", err)
			return
		}

		req.Header = r.Header
		req.Header.Add("subdomain", "orange")

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("bad response: %v", err)
			return
		}

		w.Write([]byte("it works"))
	}
}
