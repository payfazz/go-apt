package main

import (
	"log"
	"net/http"
	"time"

	"github.com/payfazz/go-apt/pkg/fazzrouter"
)

func main() {
	r := fazzrouter.BaseRoute()

	r.Prefix("v1/events", func(r *fazzrouter.Route) {
		r.Post("", func(w http.ResponseWriter, r *http.Request) {
			log.Println("POST")
		})
		r.Get("", func(w http.ResponseWriter, r *http.Request) {
			log.Println("GET")
		})
	})

	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  30 * time.Second,
		Handler:      r.Compile(),
	}
	_ = s.ListenAndServe()
}
