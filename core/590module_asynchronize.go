package core

import (
	"runtime"
	"time"
)

//var goroutineWaiter = &sync.WaitGroup{}

var goroutineManager = newGoroutineManager()
var void = struct{}{}

type GoroutineManager struct {
	count chan struct{}
}

func newGoroutineManager() *GoroutineManager {
	count := make(chan struct{}, runtime.NumCPU()*4)
	return &GoroutineManager{count: count}
}

func (gm *GoroutineManager) incr() {
	gm.count <- void
}

func (gm *GoroutineManager) decr() {
	<-gm.count
}

func (gm *GoroutineManager) wait() {
	duration := time.Duration(2) * time.Millisecond
	ticker := time.NewTicker(duration)
	for range ticker.C {
		if len(gm.count) < 1 {
			return
		}
	}
}

func (fns *InternalFunctionSet) Async(args []interface{}) Function {
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

	go func(action Function) {
		defer catch()
		goroutineManager.incr()
		defer goroutineManager.decr()

		val := action.execute()
		//fmt.Printf("cur=%p \n", fn.getCurrentStack())
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
