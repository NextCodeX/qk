package core

import "fmt"

type FloatValue struct {
	goValue float64
	ClassObject
}

func newFloatValue(raw float64) Value {
	fl := &FloatValue{goValue: raw}
	fl.ClassObject.initAsClass("Float", &fl)
	return fl
}

func (fval *FloatValue) val() interface{} {
	return fval.goValue
}

func (fval *FloatValue) isFloat() bool {
	return true
}

func (fval *FloatValue) Bytes() []byte {
	return floatToBytes(fval.goValue)
}

func (fval *FloatValue) String() string {
	return fmt.Sprint(fval.goValue)
}