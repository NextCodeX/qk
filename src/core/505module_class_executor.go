package core

import "reflect"

// 类定义
type ClassDefine struct {
	name string // class name
	fields map[string]*reflect.Type
	methods map[string]*FunctionExecutor
}

// 类对象
type ClassExecutor struct {
	define *ClassDefine
	fields map[string]*reflect.Value
}

