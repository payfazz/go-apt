package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/payfazz/go-apt/pkg/fazzcommon/httpError"
	"github.com/payfazz/go-apt/pkg/fazzcommon/value/content"
	"github.com/payfazz/go-apt/pkg/fazzcommon/value/header"
	"github.com/payfazz/go-errors"
)

// basicResponse is a struct to contain default response message
type basicResponse struct {
	Message string `json:"message"`
}

// successResponse is a struct to contain default success response
type successResponse struct {
	Success bool `json:"success"`
}

// Json is a function to return json object with given data and statusCode
func Json(w http.ResponseWriter, data interface{}, statusCode int) {
	parseHeader(w, statusCode, content.JSON)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// Text is a function to return raw text and statusCode
func Text(w http.ResponseWriter, msg string, statusCode int) {
	Json(w, basicResponse{Message: msg}, statusCode)
}

// Success is a function to return success status and statusCode
func Success(w http.ResponseWriter, success bool, statusCode int) {
	Json(w, successResponse{Success: success}, statusCode)
}

// NotFound is a function to return 404 statusCode with empty body
func NotFound(w http.ResponseWriter) {
	Json(w, struct{}{}, http.StatusNotFound)
}

// Error is a function to return http error
func Error(w http.ResponseWriter, err error) {
	cause := err
	message := fmt.Sprint("[ERROR] ", err.Error())
	if ge, ok := err.(*errors.Error); ok {
		cause = ge.Cause()
		message = ge.String()
	}

	log.Println(message)

	if _, ok := cause.(httpError.HttpErrorInterface); !ok {
		cause = httpError.InternalServer(cause)
	}

	be := cause.(httpError.HttpErrorInterface)
	Json(w, be, be.GetCode())
}

// ErrorWithLog is a function to return http error and a flag to show / hide log
// to be deprecated because Error will automatically print error message into stderr
// if possible change this to Error
func ErrorWithLog(w http.ResponseWriter, err error, showLog bool) {
	if showLog {
		log.Println("[ERROR]", err.Error())
	}
	Error(w, err)
}

// parseHeader is a function to parse content data and add it to response header
func parseHeader(w http.ResponseWriter, statusCode int, contentType string) {
	w.Header().Set(header.CONTENT_TYPE, contentType)
	w.WriteHeader(statusCode)
}
