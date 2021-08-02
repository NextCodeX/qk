package core

import (
	"strings"
	"time"
)

func init() {
	dtc := &DatetimeConstructor{}
	fs := collectFunctionInfo(&dtc)
	functionRegister("", fs)
}

const CommonDatetimeFormat = "2006-01-02 15:04:05"

type Datetime struct {
	val       time.Time
	Time int64
	Stdfmt string
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

func (dt *Datetime) String() string {
	return dt.val.Format(CommonDatetimeFormat)
}

type DatetimeConstructor struct{}

// current datetime qk object
func (dtc *DatetimeConstructor) Now() Value {
	now := time.Now()
	dt := &Datetime{now, now.Unix(), now.Format(CommonDatetimeFormat)}
	return newClassExecutor("date", dt, &dt)
}

// current timestamp
func (dtc *DatetimeConstructor) Timestamp() interface{} {
	return time.Now().Unix()
}

