package core

type BooleanValue struct {
	goValue bool
	ValueAdapter
}

func newBooleanValue(raw bool) Value {
	return &BooleanValue{goValue: raw}
}

func (boolVal *BooleanValue) val() interface{} {
	return boolVal.goValue
}

func (boolVal *BooleanValue) isBoolean() bool {
	return true
}