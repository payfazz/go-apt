package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ParseJson is a function to decode json into assigned variable i
func ParseJson(r *http.Request, i interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(i)
}

// ParseJsonWithRaw is a function to decode json into assigned variable i and return all payload as string
func ParseJsonWithRaw(r *http.Request, i interface{}) (string, error) {
	return parseJsonWithRaw(r, i)
}

// ParseRaw is a function to decode json into string
func ParseRaw(r *http.Request) (string, error) {
	return parseJsonWithRaw(r, nil)
}

// ParseMiddlewareJson is a function to decode json inside middleware while still
// keeping the request body for controller processing
func ParseMiddlewareJson(r *http.Request, i interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if nil != err {
		return err
	}

	if 0 == len(b) {
		return nil
	}

	err = r.Body.Close()
	if nil != err {
		return err
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	return json.Unmarshal(b, i)
}

// ParseQueryParam is a function to parse url query and validate required
func ParseQueryParam(r *http.Request, param map[string]string) (map[string]string, error) {
	var missingFields []string
	result := map[string]string{}

	for k, v := range param {
		param := r.URL.Query().Get(k)
		switch v {
		case "required":
			if "" == param {
				missingFields = append(missingFields, k)
			}
		}
		result[k] = param
	}

	if len(missingFields) > 0 {
		return nil, errors.New(fmt.Sprintf("missing required fields (%s)", strings.Join(missingFields, ", ")))
	}

	return result, nil
}

// QueryParamToJson is a function to convert query param payload to json string
func QueryParamToJson(r *http.Request) string {
	result, _ := json.Marshal(r.URL.Query())
	return string(result)
}

func parseJsonWithRaw(r *http.Request, i interface{}) (string, error) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if nil != err {
		return "", err
	}

	if 0 == len(b) {
		return "", err
	}

	if nil == i {
		return string(b), err
	}

	err = json.Unmarshal(b, i)
	if nil != err {
		return "", err
	}

	return string(b), nil
}
