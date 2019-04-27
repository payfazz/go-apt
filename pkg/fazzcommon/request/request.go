package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ParseJson is a function to decode json into assigned variable i
func ParseJson(r *http.Request, i interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(i)
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
		return nil, fmt.Errorf("missing required fields (%s)", strings.Join(missingFields, ", "))
	}

	return result, nil
}

// QueryParamToJson is a function to convert query param payload to json string
func QueryParamToJson(r *http.Request) string {
	result, _ := json.Marshal(r.URL.Query())
	return string(result)
}
