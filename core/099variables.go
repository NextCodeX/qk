package core

import "sync"

// 变量池
type Variables interface {
	add(name string, v Value)
	get(name string) Value
}

type VariablesImpl struct {
	pool sync.Map
}

func newVariables() Variables {
	obj := &VariablesImpl{}
	return obj
}

func (vs *VariablesImpl) add(name string, v Value) {
	for i := 0; i < 5; i++ {
		vs.pool.Store(name, v)

		if _, ok := vs.pool.Load(name); ok {
			return
		}
	}
}

func (vs *VariablesImpl) get(name string) Value {
	res, ok := vs.pool.Load(name)
	if !ok {
		return nil
	}
	return res.(Value)
}
