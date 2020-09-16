package fazzrouter

import (
	"net/http"

	"github.com/payfazz/go-middleware/common/kv"
)

type patternKeyType struct{}

var patternKey patternKeyType

func InjectPattern(pattern string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, req *http.Request) {
			kv.EnsureKVAndSet(req, patternKey, pattern)
			next(writer, req)
		}
	}
}

func GetPattern(req *http.Request) string {
	return kv.MustGet(req, patternKey).(string)
}
