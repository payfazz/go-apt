package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/payfazz/go-apt/pkg/fazzmonitor/prometheusclient"
)

// handlerA is a common handler signature.
func handlerA(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// otherHandlerA is another handler signature.
func otherHandlerA(w http.ResponseWriter, r *http.Request) error {
	return errors.New("boom")
}

func otherHandlerAdapter(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("error happened:", err)
		}
	}
}

// mustMethod represents a simple router.
func mustMethod(method string, hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		hf(w, r)
	}
}

// usually you will write this function to define how a router register a handler.
func instrumentingRouter(mux *http.ServeMux, method, pattern string, hf http.HandlerFunc) {
	mux.HandleFunc(pattern, prometheusclient.InstrumentHandlerCounter(pattern, mustMethod(method, hf)).ServeHTTP)
}

func main() {
	mux := http.NewServeMux()
	instrumentingRouter(mux, http.MethodPost, "/users", handlerA)
	instrumentingRouter(mux, http.MethodGet, "/groups", otherHandlerAdapter(otherHandlerA))
	mux.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", mux)
}
