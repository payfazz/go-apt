package assign

func WithDefault(defaultCondition bool, value interface{}, defaultValue interface{}) interface{} {
	if defaultCondition {
		return defaultValue
	}

	return value
}

func String(value string, defaultValue string) string {
	return WithDefault("" == value, value, defaultValue).(string)
}

func StringPtr(value *string, defaultValue *string) *string {
	return WithDefault(nil == value, value, defaultValue).(*string)
}

func IntPtr(value *int, defaultValue *int) *int {
	return WithDefault(nil == value, value, defaultValue).(*int)
}

func Int8Ptr(value *int8, defaultValue *int8) *int8 {
	return WithDefault(nil == value, value, defaultValue).(*int8)
}

func Int16Ptr(value *int16, defaultValue *int16) *int16 {
	return WithDefault(nil == value, value, defaultValue).(*int16)
}

func Int32Ptr(value *int32, defaultValue *int32) *int32 {
	return WithDefault(nil == value, value, defaultValue).(*int32)
}

func Int64Ptr(value *int64, defaultValue *int64) *int64 {
	return WithDefault(nil == value, value, defaultValue).(*int64)
}

func UintPtr(value *uint, defaultValue *uint) *uint {
	return WithDefault(nil == value, value, defaultValue).(*uint)
}

func Uint8Ptr(value *uint8, defaultValue *uint8) *uint8 {
	return WithDefault(nil == value, value, defaultValue).(*uint8)
}

func Uint16Ptr(value *uint16, defaultValue *uint16) *uint16 {
	return WithDefault(nil == value, value, defaultValue).(*uint16)
}

func Uint32Ptr(value *uint32, defaultValue *uint32) *uint32 {
	return WithDefault(nil == value, value, defaultValue).(*uint32)
}

func Uint64Ptr(value *uint64, defaultValue *uint64) *uint64 {
	return WithDefault(nil == value, value, defaultValue).(*uint64)
}

func Float32Ptr(value *float32, defaultValue *float32) *float32 {
	return WithDefault(nil == value, value, defaultValue).(*float32)
}

func Float64Ptr(value *float64, defaultValue *float64) *float64 {
	return WithDefault(nil == value, value, defaultValue).(*float64)
}

func BoolPtr(value *bool, defaultValue *bool) *bool {
	return WithDefault(nil == value, value, defaultValue).(*bool)
}
