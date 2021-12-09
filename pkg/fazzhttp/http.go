package fazzhttp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	GET    string = "GET"
	POST   string = "POST"
	PATCH  string = "PATCH"
	UPDATE string = "UPDATE"
	DELETE string = "DELETE"
	JSON   string = "application/json"
	XFORM  string = "application/x-www-form-urlencoded"
	FORM   string = "multipart/form-data"
)

type Interface interface {
	Get(ctx context.Context, path string, params map[string]string, headers map[string]string) (int, []byte, error)
	Send(ctx context.Context, path string, method string, contentType string, body []byte, headers map[string]string) (int, []byte, error)
	Delete(ctx context.Context, path string, headers map[string]string) (int, []byte, error)
}

// Client is a context compatible version of fazzcloud, used for migration from old non context API.
type Client struct {
	Host       string
	HttpClient *http.Client
}

// Get is a function to get the data from http call.
func (c *Client) Get(ctx context.Context, path string, params map[string]string, headers map[string]string) (int, []byte, error) {
	return c.request(ctx, path, GET, params, nil, headers)
}

// Send is a function that used to send the data into another http url.
func (c *Client) Send(ctx context.Context, path string, method string, contentType string, body []byte, headers map[string]string) (int, []byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = contentType
	return c.request(ctx, path, method, nil, body, headers)
}

// Delete is a function that used to send delete http verb.
func (c *Client) Delete(ctx context.Context, path string, headers map[string]string) (int, []byte, error) {
	return c.request(ctx, path, DELETE, nil, nil, headers)
}

func (c *Client) request(ctx context.Context, path string, method string, params map[string]string, body []byte, headers map[string]string) (int, []byte, error) {
	var req *http.Request
	var err error

	url := c.getURL(path)
	if body == nil {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	}

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// set headers from params
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// set request params from params
	if len(params) != 0 {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	response, err := c.HttpClient.Do(req)
	// parse response from response into bytes
	respBody, err := c.readResponse(response, err)
	if err != nil {
		if response != nil {
			return response.StatusCode, nil, err
		}
		return http.StatusInternalServerError, nil, err
	}

	return response.StatusCode, respBody, nil
}

func (c *Client) getURL(path string) string {
	return fmt.Sprintf("%s/%s", c.Host, path)
}
func (c *Client) readResponse(response *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode < http.StatusOK || response.StatusCode > 299 {
		return nil, errors.New(string(body))
	}

	return body, err
}

// New is a constructor function that used to http call.
func New(host string, httpClient *http.Client) Interface {
	return &Client{Host: host, HttpClient: httpClient}
}
