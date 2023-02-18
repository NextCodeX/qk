package core

import (
	"time"
)

// 延迟执行函数
func (this *InternalFunctionSet) Delay(fn Function, ms int64) {
	goroutineWaiter.Add(1)
	time.AfterFunc(time.Duration(ms)*time.Millisecond, func() {
		defer catch()
		defer goroutineWaiter.Done()

		fn.execute()
	})
}

// 每隔一个时间段，执行一次函数
func (this *InternalFunctionSet) Interval(action Function, ms int64) {
	goroutineWaiter.Add(1)
	go func() {
		defer catch()
		defer goroutineWaiter.Done()
		action.execute()

		duration := time.Duration(ms) * time.Millisecond
		ticker := time.NewTicker(duration)
		for range ticker.C {

			goroutineWaiter.Add(1)
			go func() {
				defer catch()
				defer goroutineWaiter.Done()
				action.execute()
			}()
		}
	}()
}
