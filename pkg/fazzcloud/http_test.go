package fazzcloud

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestHTTPClient_Get(t *testing.T) {
	var params = make(map[string]string)
	var headers = make(map[string]string)

	params["format"] = "json"
	headers["Content-Type"] = "multipart/form-data"

	httpClient := NewHTTPClient("https://api.ipify.org", nil)

	_, _, err := httpClient.Get("", &params, &headers)
	if err != nil {
		t.Fatalf("failed to test http client get")
	}
}

func TestFailedHTTPClient_Get(t *testing.T) {
	httpClient := NewHTTPClient("https:/api.ipify.asdf", nil)

	_, _, err := httpClient.Get("", nil, nil)
	if err == nil {
		t.Fatalf("failed to test failed http client get")
	}

	httpClient = NewHTTPClient("http://dummy.restapiexample.com/api/v1/create", nil)
	_, _, err = httpClient.Get("", nil, nil)
	httpClient.TraceRequest()
	if err == nil {
		t.Fatalf("failed to test failed http client get")
	}
}

func TestHTTPClient_Send(t *testing.T) {
	var headers = make(map[string]string)
	headers["Content-Type"] = JSON

	httpClient := NewHTTPClient("http://dummy.restapiexample.com/api/v1", nil)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := fmt.Sprintf("%d", rand.New(s1))

	tempData := struct {
		Name   string `json:"name"`
		Salary string `json:"salary"`
		Age    string `json:"age"`
	}{
		Name:   r1,
		Salary: "123123",
		Age:    "22",
	}

	params, err := json.Marshal(tempData)
	if err != nil {
		t.Fatalf("cannot marshal data")
	}
	_, _, err = httpClient.Send("create", POST, JSON, params, &headers)
	if err != nil {
		t.Fatalf("failed to test send data")
	}
}

func TestFailedHTTPClient_Send(t *testing.T) {
	httpClient := NewHTTPClient("http://dummy.restapiexample.co.id/api/v1", nil)
	tempData := struct {
		Name   string `json:"name"`
		Salary string `json:"salary"`
		Age    string `json:"age"`
	}{
		Name:   "test",
		Salary: "123123",
		Age:    "22",
	}
	params, err := json.Marshal(tempData)
	if err != nil {
		t.Fatalf("cannot marshal data")
	}
	_, _, err = httpClient.Send("create", UPDATE, JSON, params, nil)
	if err == nil {
		t.Fatalf("failed to failed test send data (UPDATE)")
	}
	_, _, err = httpClient.Send("create", POST, JSON, params, nil)
	_, _, err = httpClient.Send("create", POST, JSON, params, nil)
	if err == nil {
		t.Fatalf("failed to failed test send data (DUPLICATE DATA)")
	}
}

func TestHTTPClient_Delete(t *testing.T) {
	var params = make(map[string]string)
	var headers = make(map[string]string)

	params["format"] = "json"
	headers["Content-Type"] = "application/json"

	httpClient := NewHTTPClient("http://dummy.restapiexample.com/api/v1", nil)

	_, _, err := httpClient.Delete("delete/2", &headers)
	if err != nil {
		t.Fatalf("failed to test http client delete, %s", err)
	}
}

func TestFailedHTTPClient_Delete(t *testing.T) {
	var params = make(map[string]string)
	var headers = make(map[string]string)

	params["format"] = "json"
	headers["Content-Type"] = "application/json"

	httpClient := NewHTTPClient("http://dummy.restapiexample.co.id/api/v1", nil)

	_, _, err := httpClient.Delete("delete/2", &headers)
	if err == nil {
		t.Fatalf("failed to test failed http client delete")
	}
}
