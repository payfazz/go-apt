package formatter

import (
	"fmt"
	"github.com/leekchan/accounting"
	"github.com/microcosm-cc/bluemonday"
	"github.com/satori/go.uuid"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ToLowerFirst is a function that will change the first character of a string into a lowercase letter.
func ToLowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return str
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
	return strconv.Itoa(int(param))
}

// IntegerToString used to get string value from integer.
func IntegerToString(param int) string {
	return strconv.Itoa(param)
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
