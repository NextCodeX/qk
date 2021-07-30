package core

type NULLValue struct {
	ValueAdapter
}

func newNULLValue() Value {
	return &NULLValue{}
}

func (null *NULLValue) val() interface{} {
	return "null"
}
func (null *NULLValue) isNULL() bool {
	return true
}
