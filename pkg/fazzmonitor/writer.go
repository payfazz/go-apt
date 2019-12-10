package fazzmonitor

import "net/http"

// Writer is a struct that used to wrap the response writer
type Writer struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader is a function that used to wrap write header default from status writer
func (w *Writer) WriteHeader(status int) {
	w.StatusCode = status
	w.ResponseWriter.WriteHeader(status)
}

// Write is a function that used to wrap write function
func (w *Writer) Write(b []byte) (int, error) {
	if w.StatusCode == 0 {
		w.StatusCode = 200
	}
	n, err := w.ResponseWriter.Write(b)
	return n, err
}
