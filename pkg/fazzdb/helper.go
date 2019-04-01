package fazzdb

import (
	"unicode"
)

// toLowerFirst is a function that will change the first character of a string into a lowercase letter
func toLowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return str
}