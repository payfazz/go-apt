package fazzrouter

import (
	"net/http"

	"github.com/payfazz/go-middleware/common/kv"
)

type patternKeyType struct{}

var patternKey patternKeyType

func InjectPattern(pattern string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			kv.Set(request, patternKey, pattern)
			next(writer, request)
		}
	}
}

func GetPattern(request *http.Request) string {
	return kv.MustGet(request, patternKey).(string)
}
