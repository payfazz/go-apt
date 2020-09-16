package fazzrouter

import (
	"log"
	"net/http"

	"github.com/payfazz/go-middleware/common/kv"
)

type patternKeyType struct{}

var patternKey patternKeyType

func InjectPattern(pattern string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, req *http.Request) {
			log.Println("inject pattern:", pattern)
			next(writer, kv.EnsureKVAndSet(req, patternKey, pattern))
		}
	}
}

func GetPattern(req *http.Request) string {
	return kv.MustGet(req, patternKey).(string)
}
