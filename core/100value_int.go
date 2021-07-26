package core

type IntValue struct {
	goValue int64
}

func newIntValue(raw int64) *IntValue {
	return &IntValue{raw}
}

func (ival *IntValue) val() interface{} {
	return ival.goValue
}
func (ival *IntValue) isNULL() bool {
	return false
}
func (ival *IntValue) isInt() bool {
	return true
}
func (ival *IntValue) isFloat() bool {
	return false
}
func (ival *IntValue) isBoolean() bool {
	return false
}
func (ival *IntValue) isString() bool {
	return false
}
func (ival *IntValue) isAny() bool {
	return false
}
func (ival *IntValue) isClass() bool {
	return false
}
func (ival *IntValue) isJsonArray() bool {
	return false
}
func (ival *IntValue) isJsonObject() bool {
	return false
}
