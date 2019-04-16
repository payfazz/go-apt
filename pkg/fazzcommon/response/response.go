package response

import (
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/fazzcommon/httpError"
	"github.com/payfazz/go-apt/pkg/fazzcommon/value/content"
	"github.com/payfazz/go-apt/pkg/fazzcommon/value/header"
	"net/http"
)

type Basic struct {
	Message string `json:"message"`
}

func Json(w http.ResponseWriter, data interface{}, statusCode int) {
	parseHeader(w, statusCode, content.Json, data)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func Text(w http.ResponseWriter, msg string, statusCode int) {
	Json(w, Basic{Message: msg}, statusCode)
}

func Error(w http.ResponseWriter, err httpError.HttpErrorInterface) {
	Json(w, err, err.GetCode())
}

func parseHeader(w http.ResponseWriter, statusCode int, contentType string, data interface{}) {
	w.Header().Set(header.ContentType, contentType)
	w.WriteHeader(statusCode)
}