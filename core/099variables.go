package core

// 变量池
type Variables map[string]Value

func newVariables() Variables {
	return make(map[string]Value)
}

func (vs Variables) isEmpty() bool {
	return vs == nil || len(vs) < 1
}

func (vs Variables) add(name string, v Value) {
	vs[name] = v
}

func (vs Variables) get(name string) Value {
	if vs.isEmpty() {
		return nil
	}
	res, ok := vs[name]
	if ok {
		return res
	}
	return nil
}
