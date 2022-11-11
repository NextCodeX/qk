package core

import (
	"github.com/robfig/cron"
)

// 定时任务管理器
var qkcron *cron.Cron
var cronRunning = false

func (this *InternalFunctionSet) Cron(expr string, fn Function) {
	if qkcron == nil {
		qkcron = cron.New()
	}

	_ = qkcron.AddFunc(expr, func() {
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
	goroutineWaiter.Add(1)
	go func() {
		defer goroutineWaiter.Done()

		qkcron.Run()
	}()
}
