package fazzcloud

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	POST   string = "POST"
	PATCH  string = "PATCH"
	UPDATE string = "UPDATE"
	JSON   string = "application/json"
	XFORM  string = "application/x-www-form-urlencoded"
	FORM   string = "multipart/form-data"
)

// HTTPClient is a struct that contain host and http client that will be used for http call
type HTTPClient struct {
	host       string
	httpClient *http.Client
}

func (hr *HTTPClient) getURL(path string) string {
	return fmt.Sprintf("%s/%s", hr.host, path)
}

func (hr *HTTPClient) readResponse(response *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode < http.StatusOK || response.StatusCode > 299 {
		return body, errors.New(string(body))
	}
	return body, err
}

// Get is a function to get the data from http call.
func (hr *HTTPClient) Get(path string, params *map[string]string, headers *map[string]string) (int, []byte, error) {
	url := hr.getURL(path)

	req, err := http.NewRequest("GET", url, nil)

	// set headers from params
	if headers != nil && len(*headers) != 0 {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}

	// set request params from params
	if params != nil && len(*params) != 0 {
		q := req.URL.Query()
		for key, value := range *params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	response, err := hr.httpClient.Do(req)

	// parse response from response into bytes
	resp, err := hr.readResponse(response, err)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// return response to caller
	return response.StatusCode, resp, err
}

// Post is a function that used to post the data into another http url.
func (hr *HTTPClient) Send(path string, method string, contentType string, params []byte, headers *map[string]string) (int, []byte, error) {
	url := hr.getURL(path)

	// send request and params
	req, err := http.NewRequest(method, url, bytes.NewBuffer(params))

	// set headers from params
	req.Header.Set("Content-Type", contentType)
	if headers != nil && len(*headers) != 0 {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	response, err := hr.httpClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// read and parse responses
	resp, err := hr.readResponse(response, err)
	return response.StatusCode, resp, err
}

// Delete is a function that used to send delete http verb.
func (hr *HTTPClient) Delete(path string, headers *map[string]string) (int, []byte, error) {
	url := hr.getURL(path)

	req, err := http.NewRequest("DELETE", url, nil)

	// set headers from params
	if headers != nil && len(*headers) != 0 {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	response, err := hr.httpClient.Do(req)

	// parse response from response into bytes
	resp, err := hr.readResponse(response, err)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// return response to caller
	return response.StatusCode, resp, err
}

// NewHTTPClient is a constructor function that used to http call.
func NewHTTPClient(host string) *HTTPClient {
	timeout := time.Duration(5 * time.Second)
	httpClient := &http.Client{
		Timeout: timeout,
	}
	return &HTTPClient{host: host, httpClient: httpClient}
}
