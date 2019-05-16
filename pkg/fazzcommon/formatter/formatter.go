package formatter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/leekchan/accounting"
	"github.com/microcosm-cc/bluemonday"
	uuid "github.com/satori/go.uuid"
)

// ToLowerFirst is a function that will change the first character of a string into a lowercase letter.
func ToLowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return str
}

// SliceJoins is a function to join slice into string with chosen delimiter
func SliceJoins(v interface{}, delimiter string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(v), " ", delimiter, -1), "[]")
}

// LeftPad2Len used to fill string with some formats.
// example: LeftPad2Len(9, "0", 4)
// result: 0009
func LeftPad2Len(str string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + str
	return retStr[(len(retStr) - overallLen):]
}

// StringToInteger used to get int value from string, without try to catch an error.
func StringToInteger(param string) int {
	val, _ := strconv.Atoi(param)
	return val
}

// StringToFloat used to get float64 value from string, without try to catch an error.
func StringToFloat(param string) float64 {
	val, _ := strconv.ParseFloat(param, 10)
	return val
}

// FloatToString used to get string value from float64.
func FloatToString(param float64) string {
	return fmt.Sprintf("%0.5f", param)
}

// IntegerToString used to get string value from integer.
func IntegerToString(param int) string {
	return fmt.Sprintf("%d", param)
}

// SliceUint8ToString used to convert []uint8 to string
func SliceUint8ToString(ui []uint8) string {
	runes := make([]rune, len(ui))
	for i, v := range ui {
		runes[i] = rune(v)
	}
	return string(runes)
}

// StringToInt64 used to get int64 value from string, without try to catch an error.
func StringToInt64(param string) int64 {
	val, _ := strconv.ParseInt(param, 10, 64)
	return val
}

// GenerateStringUUID used to generate UUID (return string).
func GenerateStringUUID() string {
	return fmt.Sprintf("%s", uuid.NewV4())
}

// SanitizePhone used to sanitize Indonesia's phone number from 08 to +628.
func SanitizePhone(input string) string {
	stripExp, _ := regexp.Compile(`[^0-9]+`)
	stripped := stripExp.ReplaceAllString(input, "")

	preExp, _ := regexp.Compile(`^(\+?62|0)([0-9]*)`)
	matches := preExp.FindStringSubmatch(stripped)
	if len(matches) == 0 {
		return ""
	}
	sanitized := fmt.Sprintf("%s%s", "+62", matches[len(matches)-1])
	return sanitized
}

// UnSanitizePhone used to unsanitize Indonesia's phone number from +628 to 08.
func UnSanitizePhone(input string) string {
	stripExp, _ := regexp.Compile(`[^0-9]+`)
	stripped := stripExp.ReplaceAllString(input, "")

	preExp, _ := regexp.Compile(`^(\+?62|0)([0-9]*)`)
	matches := preExp.FindStringSubmatch(stripped)
	if len(matches) == 0 {
		return ""
	}
	sanitized := fmt.Sprintf("%s%s", "0", matches[len(matches)-1])
	return sanitized
}

// CleanString to guard string from any SQL Injection Syntax.
func CleanString(param *string) string {
	if param == nil || *param == "" {
		return ""
	}
	var replacer = strings.NewReplacer("exec", "", "--", "", "DROP", "", "EXEC", "", "drop", "", "'", "", ";", "")
	cleanParam := replacer.Replace(*param)
	p := bluemonday.UGCPolicy()
	return p.Sanitize(cleanParam)
}

// MoneyFormat is a function that used to format the money.
func MoneyFormat(param float64) string {
	ac := accounting.Accounting{Precision: 2}
	return ac.FormatMoney(param)
}

// ConvertMapToString is a function that used to convert map to string
func ConvertMapToString(value map[string]string) string {
	if len(value) < 1 {
		return ""
	}
	str := "map["
	for key, val := range value {
		str = fmt.Sprintf("%s%s:%s", str, key, val)
	}
	return fmt.Sprintf("%s]", str)
}

// ToStringPtr used to return string pointer from param
func ToStringPtr(param interface{}) *string {
	result := fmt.Sprint(param)
	return &result
}

// ToFloat32Ptr used to return float32 pointer from param
func ToFloat32Ptr(param interface{}) *float32 {
	var result float32
	switch i := param.(type) {
	case float32:
		result = i
	case float64:
		result = float32(i)
	case int64:
		result = float32(i)
	case int32:
		result = float32(i)
	case int:
		result = float32(i)
	case uint64:
		result = float32(i)
	case uint32:
		result = float32(i)
	case uint:
		result = float32(i)
	default:
		return nil
	}
	return &result
}

// ToFloat64Ptr used to return float64 pointer from param
func ToFloat64Ptr(param interface{}) *float64 {
	var result float64
	switch i := param.(type) {
	case float64:
		result = i
	case float32:
		result = float64(i)
	case int64:
		result = float64(i)
	case int32:
		result = float64(i)
	case int:
		result = float64(i)
	case uint64:
		result = float64(i)
	case uint32:
		result = float64(i)
	case uint:
		result = float64(i)
	default:
		return nil
	}
	return &result
}

// ToIntPtr used to return int pointer from param
func ToIntPtr(param interface{}) *int {
	var result int
	switch i := param.(type) {
	case float64:
		result = int(i)
	case float32:
		result = int(i)
	case int64:
		result = int(i)
	case int32:
		result = int(i)
	case int:
		result = i
	case uint64:
		result = int(i)
	case uint32:
		result = int(i)
	case uint:
		result = int(i)
	default:
		return nil
	}
	return &result
}

// ToInt64Ptr used to return int64 pointer from param
func ToInt64Ptr(param interface{}) *int64 {
	var result int64
	switch i := param.(type) {
	case float64:
		result = int64(i)
	case float32:
		result = int64(i)
	case int64:
		result = i
	case int32:
		result = int64(i)
	case int:
		result = int64(i)
	case uint64:
		result = int64(i)
	case uint32:
		result = int64(i)
	case uint:
		result = int64(i)
	default:
		return nil
	}
	return &result
}
