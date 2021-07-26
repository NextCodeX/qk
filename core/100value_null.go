package core

type NULLValue struct {}

func newNULLValue() *NULLValue {
	return &NULLValue{}
}

func (null *NULLValue) val() interface{} {
	return "null"
}
func (null *NULLValue) isNULL() bool {
	return true
}
func (null *NULLValue) isInt() bool {
	return false
}
func (null *NULLValue) isFloat() bool {
	return false
}
func (null *NULLValue) isBoolean() bool {
	return false
}
func (null *NULLValue) isString() bool {
	return false
}
func (null *NULLValue) isAny() bool {
	return false
}
func (null *NULLValue) isClass() bool {
	return false
}
func (null *NULLValue) isJsonArray() bool {
	return false
}
func (null *NULLValue) isJsonObject() bool {
	return false
}