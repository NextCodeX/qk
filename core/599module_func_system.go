package core

import "time"

func (fns *InternalFunctionSet) Sleep(t int64)  {
	time.Sleep(time.Duration(t) * time.Millisecond)
}