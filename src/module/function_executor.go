package module

import (
	"reflect"
	"fmt"
)

type FunctionExecutor struct {
	name string
	ins []reflect.Type
	outs []reflect.Type
	obj reflect.Value
}


