package core

import (
	"strings"
	"time"
)

func (mr *ModuleRegister) DateModuleInit() {
	dtc := &DatetimeConstructor{}
	fs := collectFunctionInfo(&dtc)
	functionRegister("", fs)
}

type Datetime struct {
	val       time.Time
	Timestamp int64
	StandardFormat string
}

// date format: y-M-d H:m:s:S
func (dt *Datetime) Format(tmpl string) string {
	tmpl = strings.Replace(tmpl, "y", "2006", 1)
	tmpl = strings.Replace(tmpl, "M", "01", 1)
	tmpl = strings.Replace(tmpl, "d", "02", 1)
	tmpl = strings.Replace(tmpl, "H", "15", 1)
	tmpl = strings.Replace(tmpl, "m", "04", 1)
	tmpl = strings.Replace(tmpl, "s", "05", 1)
	tmpl = strings.Replace(tmpl, "S", "000", 1)
	return dt.val.Format(tmpl)
}

type DatetimeConstructor struct{}

func (dtc *DatetimeConstructor) Now() *ClassExecutor {
	now := time.Now()
	dt := &Datetime{now, now.Unix(), now.Format("2006-01-02 15:04:05")}
	return newClassExecutor("date", dt, &dt)
}

