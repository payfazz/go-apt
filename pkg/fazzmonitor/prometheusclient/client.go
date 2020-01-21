package prometheusclient

import (
	"fmt"
	"net/http"
)

type RoutePattern interface {
	Get(req *http.Request) string
}

type prometheusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (writer *prometheusResponseWriter) Code() string {
	return fmt.Sprint(writer.statusCode)
}

func (writer *prometheusResponseWriter) WriteHeader(statusCode int) {
	writer.statusCode = statusCode
	writer.ResponseWriter.WriteHeader(statusCode)
}

func wrapResponseWriter(writer http.ResponseWriter) *prometheusResponseWriter {
	if _, ok := writer.(*prometheusResponseWriter); ok {
		return writer.(*prometheusResponseWriter)
	}

	return &prometheusResponseWriter{ResponseWriter: writer}
}
