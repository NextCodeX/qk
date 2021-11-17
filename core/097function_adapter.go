package core

type FunctionAdapter struct {
	name     string
	instance Function
	StatementAdapter
	ValueAdapter
}

func (f *FunctionAdapter) init(name string, obj Function) {
	f.name = name
	f.instance = obj
	f.StatementAdapter.initStatement(obj)
}
func (f *FunctionAdapter) setArgs(args interface{}) {
	if f1, ok := f.instance.(*InternalFunction); ok {
		f1.setGoArgs(args.([]interface{}))
	} else if f2, ok := f.instance.(*CustomFunction); ok {
		f2.setQkArgs(args.([]Value))
	}
}
func (f *FunctionAdapter) parentFrame() Frame {
	return f.instance.getParent()
}

func (f *FunctionAdapter) varList() Variables {
	return nil
}
func (f *FunctionAdapter) val() interface{} {
	return f.instance
}
func (f *FunctionAdapter) typeName() string {
	return "Function"
}
func (f *FunctionAdapter) isFunction() bool {
	return true
}
func (f *FunctionAdapter) isObject() bool {
	return true
}

func (f *FunctionAdapter) get(key string) Value {
	if key == "type" {
		return callable(func() Value {
			return newQKValue("Function")
		})
	} else {
		return NULL
	}
}

func (f *FunctionAdapter) String() string {
	return f.name + "()"
}
