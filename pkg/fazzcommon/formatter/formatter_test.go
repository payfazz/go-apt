package formatter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestToLowerFirst(t *testing.T) {
	ToLowerFirst("testing")
}

func TestEmptyToLowerFirst(t *testing.T) {
	ToLowerFirst("")
}

func TestLeftPad2Len(t *testing.T) {
	LeftPad2Len("1", "0", 4)
}

func TestStringToInteger(t *testing.T) {
	StringToInteger("10")
}

func TestStringToFloat(t *testing.T) {
	StringToFloat("10")
}

func TestFloatToString(t *testing.T) {
	FloatToString(10.5)
}

func TestIntegerToString(t *testing.T) {
	IntegerToString(10)
}

func TestSanitizePhone(t *testing.T) {
	SanitizePhone("08123456789")
}

func TestEmptySanitizePhone(t *testing.T) {
	SanitizePhone("")
}

func TestGenerateStringUUID(t *testing.T) {
	GenerateStringUUID()
}

func TestUnSanitizePhone(t *testing.T) {
	UnSanitizePhone("+62812345689")
}

func TestEmptyUnSanitizePhone(t *testing.T) {
	UnSanitizePhone("")
}

func TestCleanString(t *testing.T) {
	str := "--test"
	CleanString(&str)
}

func TestEmptyCleanString(t *testing.T) {
	CleanString(nil)
}

func TestMoneyFormat(t *testing.T) {
	MoneyFormat(100000)
}

func TestResponseWithData(t *testing.T) {
	w := httptest.NewRecorder()
	ResponseWithData(w, http.StatusOK, "ok")
}

func TestJSONDecode(t *testing.T) {
	h := httptest.NewRequest("GET", "https://test.com", nil)
	JSONDecode(h, "")
}
