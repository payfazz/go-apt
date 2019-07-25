package formatter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplaceStrings(t *testing.T) {
	source := "this is a fox jumping over a crude board inside the room"
	olds := []string{"is", "de "}
	news := []string{"on", " de"}
	result := ReplaceStrings(source, olds, news)
	require.Equal(t, "thon on a fox jumping over a cru deboard insi dethe room", result)

	olds = []string{"is", "de "}
	news = []string{"on"}
	result = ReplaceStrings(source, olds, news)
	require.Equal(t, "this is a fox jumping over a crude board inside the room", result)
}

func TestToLowerFirst(t *testing.T) {
	result := ToLowerFirst("TESTING")
	require.Equal(t, "tESTING", result)
}

func TestEmptyToLowerFirst(t *testing.T) {
	result := ToLowerFirst("")
	require.Equal(t, "", result)
}

func TestLeftPad2Len(t *testing.T) {
	result := LeftPad2Len("1", "0", 4)
	require.Equal(t, "0001", result)

	result = LeftPad2Len("12", "x", 5)
	require.Equal(t, "xxx12", result)
}

func TestStringToInteger(t *testing.T) {
	result := StringToInteger("10")
	require.Equal(t, 10, result)
}

func TestStringToFloat(t *testing.T) {
	result := StringToFloat("10")
	require.Equal(t, float64(10), result)
}

func TestFloatToString(t *testing.T) {
	result := FloatToString(10.5)
	require.Equal(t, "10.50000", result)
}

func TestIntegerToString(t *testing.T) {
	result := IntegerToString(10)
	require.Equal(t, "10", result)
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
	result := StringToInt64("10")
	require.Equal(t, int64(10), result)
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
