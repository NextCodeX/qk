package core

import (
	"sync"
)

var goroutineWaiter sync.WaitGroup

func (this *InternalFunctionSet) Async(args []interface{}) Function {
	if len(args) < 1 {
		runtimeExcption("async(): the first parameter is required.")
	}
	fn, ok := args[0].(Function)
	if !ok {
		runtimeExcption("async(): the first parameter type must be Function.")
	}

	// 设置异步函数参数
	var params []Value
	for _, arg := range args[1:] {
		params = append(params, newQKValue(arg))
	}
	fn.setArgs(params)

	res := make(chan Value, 1)

	goroutineWaiter.Add(1)
	go func(action Function) {
		defer catch()
		defer goroutineWaiter.Done()

		val := action.execute()
		if val == nil {
			res <- nil
		} else {
			res <- val.value()
		}
	}(fn)

	return callable(func() Value {
		return <-res
	})
}
