package request

import (
	"encoding/json"
	"net/http"
)

func ParseJson(r *http.Request, i interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(i)
}
