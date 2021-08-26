package core

type NULLValue struct {
	ClassObject
}

func newNULLValue() Value {
	nl := &NULLValue{}
	nl.ClassObject.initAsClass("NULL", &nl)
	return nl
}

func (null *NULLValue) val() interface{} {
	return "null"
}
func (null *NULLValue) isNULL() bool {
	return true
}

func (null *NULLValue) String() string {
	return "null"
}
