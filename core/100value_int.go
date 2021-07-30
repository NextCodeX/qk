package core

type IntValue struct {
	goValue int64
	ValueAdapter
}

func newIntValue(raw int64) Value {
	return &IntValue{goValue: raw}
}

func (ival *IntValue) val() interface{} {
	return ival.goValue
}

func (ival *IntValue) isInt() bool {
	return true
}