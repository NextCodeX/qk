package core

import (
	"fmt"
	"time"
)

func (mr *ModuleRegister) DateModuleInit() {
	now := time.Now()
	dt := &Datetime{now, now.Unix(), now.Format("2006-01-02 15:04:05")}
	fs := collectFieldInfo(dt)
	for k, v := range fs {
		fmt.Println(k, "->", v)
	}
}

type Datetime struct {
	val       time.Time
	Timestamp int64
	Format    string
}

