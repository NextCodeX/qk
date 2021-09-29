package core

import (
	"github.com/robfig/cron"
)

// 定时任务管理器
var qkcron = cron.New()
var cronRunning = false

func (fns *InternalFunctionSet) Cron(expr string, fn Function) {
	qkcron.AddFunc(expr, func() {
		fn.execute()
	})

	cronStart()
}

// 定时任务开启
func cronStart() {
	if cronRunning {
		return
	}

	cronRunning = true
	//goroutineWaiter.Add(1)
	go func() {
		goroutineManager.incr()
		defer goroutineManager.decr()

		qkcron.Run()
	}()
}
