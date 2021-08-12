package core

import (
	"strings"
	"time"
)


// current timestamp (Microsecond)
func (fns *InternalFunctionSet) Timestamp() interface{} {
	return time.Now().UnixNano() / 1e6
}

// current timestamp （Nanosecond）
func (fns *InternalFunctionSet) Nanosecond() interface{} {
	return time.Now().UnixNano()
}

// current datetime qk object
func (fns *InternalFunctionSet) Now() Value {
	dt := &Datetime{time.Now()}
	return newClass("date", &dt)
}

func (fns *InternalFunctionSet) NewDate(ms int64) Value {
	dt := &Datetime{time.Unix(0, ms * 1e6)}
	return newClass("date", &dt)
}

func (fns *InternalFunctionSet) NewDate1(ns int64) Value {
	dt := &Datetime{time.Unix(0, ns)}
	return newClass("date", &dt)
}

func (fns *InternalFunctionSet) ParseDate(dateStr string, tmpl string) Value {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		runtimeExcption(err)
	}
	tmpl = convertDateFmt(tmpl)
	d, err := time.ParseInLocation(tmpl, dateStr, loc)
	if err != nil {
		runtimeExcption(err)
	}
	dt := &Datetime{d}
	return newClass("date", &dt)
}

const CommonDatetimeFormat = "2006-01-02 15:04:05"

type Datetime struct {
	val       time.Time
}

// date format: y-M-d H:m:s:S
func (dt *Datetime) Format(tmpl string) string {
	tmpl = convertDateFmt(tmpl)
	return dt.val.Format(tmpl)
}

func (dt *Datetime) Year() int {
	return dt.val.Year()
}

func (dt *Datetime) Month() int {
	return int(dt.val.Month())
}

func (dt *Datetime) Day() int {
	return dt.val.Day()
}

func (dt *Datetime) Hour() int {
	return dt.val.Hour()
}

func (dt *Datetime) Minute() int {
	return dt.val.Minute()
}

func (dt *Datetime) Second() int {
	return dt.val.Second()
}

func (dt *Datetime) Ms() int64 {
	return dt.val.UnixNano() / 1e6
}

func (dt *Datetime) Ns() int64 {
	return dt.val.UnixNano()
}

func (dt *Datetime) String() string {
	return dt.val.Format(CommonDatetimeFormat)
}


func convertDateFmt(tmpl string) string {
	tmpl = strings.Replace(tmpl, "y", "2006", 1)
	tmpl = strings.Replace(tmpl, "M", "01", 1)
	tmpl = strings.Replace(tmpl, "d", "02", 1)
	tmpl = strings.Replace(tmpl, "h", "15", 1)
	tmpl = strings.Replace(tmpl, "m", "04", 1)
	tmpl = strings.Replace(tmpl, "s", "05", 1)
	tmpl = strings.Replace(tmpl, "S", "000", 1)
	return tmpl
}
