package core

import (
	"fmt"
	"sync"
)

var goroutineWaiter = &sync.WaitGroup{}

func (fns *InternalFunctionSet) Async(args []interface{}) Function {
	if len(args) < 1 {
		runtimeExcption("async(): the first parameter is required.")
	}
	fn, ok := args[0].(Function)
	if !ok {
		runtimeExcption("async(): the first parameter type must be Function.")
	}

	var params []Value
	for _, arg := range args[1:] {
		params = append(params, newQKValue(arg))
	}
	fn.setArgs(params)

	res := make(chan Value, 1)
	goroutineWaiter.Add(1)
	go func() {
		defer goroutineWaiter.Done()
		defer func() {
			// 子协程的全局异常处理
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		val := fn.execute()
		if val == nil {
			res <- nil
		} else {
			res <- val.value()
		}
	}()

	return newAnonymousFunc(func() Value {
		return <-res
	})
}