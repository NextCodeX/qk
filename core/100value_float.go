package core


type FloatValue struct {
	goValue float64
	ValueAdapter
}

func newFloatValue(raw float64) Value {
	return &FloatValue{goValue: raw}
}

func (fval *FloatValue) val() interface{} {
	return fval.goValue
}

func (fval *FloatValue) isFloat() bool {
	return true
}