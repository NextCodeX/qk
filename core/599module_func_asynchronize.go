package core

import (
	"sync"
)

var goroutineWaiter = &sync.WaitGroup{}

func (fns *InternalFunctionSet) Async(fn Function) Function {
	res := make(chan Value, 1)
	goroutineWaiter.Add(1)
	go func() {
		defer goroutineWaiter.Done()

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