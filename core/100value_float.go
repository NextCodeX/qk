package core


type FloatValue struct {
	goValue float64
}

func newFloatValue(raw float64) *FloatValue {
	return &FloatValue{raw}
}

func (fval *FloatValue) val() interface{} {
	return fval.goValue
}
func (fval *FloatValue) isNULL() bool {
	return false
}
func (fval *FloatValue) isInt() bool {
	return false
}
func (fval *FloatValue) isFloat() bool {
	return true
}
func (fval *FloatValue) isBoolean() bool {
	return false
}
func (fval *FloatValue) isString() bool {
	return false
}
func (fval *FloatValue) isAny() bool {
	return false
}
func (fval *FloatValue) isClass() bool {
	return false
}
func (fval *FloatValue) isJsonArray() bool {
	return false
}
func (fval *FloatValue) isJsonObject() bool {
	return false
}