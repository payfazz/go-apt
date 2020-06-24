package fazzdb

import (
	"fmt"
	"strings"
)

// JSONPath is type for json path slice (item type must be string or integer)
type JSONPath []interface{}

// JSONFieldOp add Json path to operator
func JSONFieldOp(path JSONPath, op Operator) Operator {
	pathStr := make([]string, len(path))
	for i := range path {
		pathStr[i] = fmt.Sprint(path[i])
	}
	result := fmt.Sprintf("%s '{%s}' %s", OP_JSON_GET_PATH_TEXT, strings.Join(pathStr, ","), op)
	return Operator(result)

}
