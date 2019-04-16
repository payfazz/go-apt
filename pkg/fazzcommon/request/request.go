package request

import (
	"encoding/json"
	"net/http"
)

// ParseJson is a function to decode json into assigned variable i
func ParseJson(r *http.Request, i interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(i)
}
