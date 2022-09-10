package core

import (
	"fmt"
)

type AnyValue struct {
	goValue interface{}
	ClassObject
}

func newAnyValue(raw interface{}) Value {
	obj := &AnyValue{goValue: raw}
	obj.initAsClass("Anything", &obj)
	return obj
}

func (this *AnyValue) val() interface{} {
	return this.goValue
}
func (this *AnyValue) isAny() bool {
	return true
}

func (this *AnyValue) Pr() {
	fmt.Println(this.String())
}
func (this *AnyValue) String() string {
	return fmt.Sprint(this.goValue)
}
