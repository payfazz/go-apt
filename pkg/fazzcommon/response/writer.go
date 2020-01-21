package response

import (
	"fmt"
	"net/http"
)

// Writer is a struct that used to wrap the response writer
type Writer struct {
	http.ResponseWriter
	statusCode int
}

func (w *Writer) Code() string {
	return fmt.Sprint(w.statusCode)
}

// WriteHeader is a function that used to wrap write header default from status writer
func (w *Writer) WriteHeader(status int) {
	w.statusCode = status
	w.ResponseWriter.WriteHeader(status)
}

// Write is a function that used to wrap write function
func (w *Writer) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = 200
	}
	return w.ResponseWriter.Write(b)
}

func WrapWriter(writer http.ResponseWriter) *Writer {
	if _, ok := writer.(*Writer); ok {
		return writer.(*Writer)
	}

	return &Writer{ResponseWriter: writer}
}
