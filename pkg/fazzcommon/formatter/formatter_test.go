package formatter

import (
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

func TestSliceUint8ToString(t *testing.T) {
	SliceUint8ToString([]uint8{128, 128, 128, 128, 128, 128, 128})
}

func TestStringToInt64(t *testing.T) {
	StringToInt64("10")
}

func TestConvertMapToString(t *testing.T) {
	ConvertMapToString(map[string]string{})
	ConvertMapToString(map[string]string{"test": "test", "test2": "test2"})
}

func TestToStringPtr(t *testing.T) {
	ToStringPtr("test")
}

func TestToFloat32Ptr(t *testing.T) {
	ToFloat32Ptr(25.6)
	ToFloat32Ptr(25)
}

func TestToFloat64Ptr(t *testing.T) {
	ToFloat64Ptr(25.6)
	ToFloat64Ptr(25)
}

func TestToIntPtr(t *testing.T) {
	ToIntPtr(25.6)
	ToIntPtr(25)
}

func TestToInt64Ptr(t *testing.T) {
	ToInt64Ptr(25.6)
	ToInt64Ptr(25)
}

func TestSliceJoins(t *testing.T) {
	SliceJoins([]string{"hallo", "hallo2"}, "-")
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
