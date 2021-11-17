package core

type Frame interface {
	varList() Variables // 当前作用域的变量列表
	parentFrame() Frame // 父帧
}

type Function interface {
	Frame
	Statement
	Value
	setArgs(args interface{})
}
