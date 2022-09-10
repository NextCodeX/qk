package core

import "fmt"

type BooleanValue struct {
	goValue bool
	ClassObject
}

func newBooleanValue(raw bool) Value {
	bl := &BooleanValue{goValue: raw}
	bl.ClassObject.initAsClass("Boolean", &bl)
	return bl
}

func (this *BooleanValue) val() interface{} {
	return this.goValue
}

func (this *BooleanValue) isBoolean() bool {
	return true
}

func (this *BooleanValue) Pr() {
	fmt.Println(this.String())
}
func (this *BooleanValue) String() string {
	return fmt.Sprint(this.goValue)
}
