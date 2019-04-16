package response

import (
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/fazzcommon/httpError"
	"github.com/payfazz/go-apt/pkg/fazzcommon/value/content"
	"github.com/payfazz/go-apt/pkg/fazzcommon/value/header"
	"net/http"
)

// Basic is a struct to contain default response message
type Basic struct {
	Message string `json:"message"`
}

// Json is a function to return json object with given data and statusCode
func Json(w http.ResponseWriter, data interface{}, statusCode int) {
	parseHeader(w, statusCode, content.Json)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// Text is a function to return raw text and statusCode
func Text(w http.ResponseWriter, msg string, statusCode int) {
	Json(w, Basic{Message: msg}, statusCode)
}

// Error is a function to return http error
func Error(w http.ResponseWriter, err httpError.HttpErrorInterface) {
	Json(w, err, err.GetCode())
}

// parseHeader is a function to parse content data and add it to response header
func parseHeader(w http.ResponseWriter, statusCode int, contentType string) {
	w.Header().Set(header.ContentType, contentType)
	w.WriteHeader(statusCode)
}