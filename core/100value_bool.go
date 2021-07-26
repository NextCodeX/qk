package core

type BooleanValue struct {
	goValue bool
}

func newBooleanValue(raw bool) *BooleanValue {
	return &BooleanValue{raw}
}

func (boolVal *BooleanValue) val() interface{} {
	return boolVal.goValue
}
func (boolVal *BooleanValue) isNULL() bool {
	return false
}
func (boolVal *BooleanValue) isInt() bool {
	return false
}
func (boolVal *BooleanValue) isFloat() bool {
	return false
}
func (boolVal *BooleanValue) isBoolean() bool {
	return true
}
func (boolVal *BooleanValue) isString() bool {
	return false
}
func (boolVal *BooleanValue) isAny() bool {
	return false
}
func (boolVal *BooleanValue) isClass() bool {
	return false
}
func (boolVal *BooleanValue) isJsonArray() bool {
	return false
}
func (boolVal *BooleanValue) isJsonObject() bool {
	return false
}