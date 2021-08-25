package core

import (
	"time"
)

// 延迟执行函数
func (fns *InternalFunctionSet) Delay(fn Function, ms int64) {
	goroutineWaiter.Add(1)
	time.AfterFunc(time.Duration(ms) * time.Millisecond, func() {
		defer goroutineWaiter.Done()
		fn.execute()
	})
}

// 每隔一个时间段，执行一次函数
func (fns *InternalFunctionSet) Interval(args []interface{}) {
	if len(args) < 1 {
		runtimeExcption("interval(fn[, ms][, async]) parameter fn is required ")
	}

	fn, ok := args[0].(Function)
	if !ok {
		runtimeExcption("interval(fn[, ms][, async]) parameter fn type must be Function ")
	}
	duration := time.Second
	if len(args) > 1 {
		ms, ok := args[1].(int64)
		if ok {
			duration = time.Duration(ms) * time.Millisecond
		}
	}
	var async bool
	if len(args) > 2 {
		async = args[2].(bool)
	}
	ticker := time.NewTicker(duration)
	for range ticker.C {
		if async {
			goroutineWaiter.Add(1)
			go func() {
				defer goroutineWaiter.Done()
				fn.execute()
			}()
		} else {
			fn.execute()
		}
	}
}