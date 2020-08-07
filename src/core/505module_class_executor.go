package core

import "reflect"

type ClassDefine struct {
	name string // class name
	fields map[string]*reflect.Type
	methods map[string]*FunctionExecutor
}

