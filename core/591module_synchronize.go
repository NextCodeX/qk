package core

import "sync"

var globalLock = &sync.Mutex{}

// 执行同步操作
func (fns *InternalFunctionSet) Sync(fn Function) {
	globalLock.Lock()
	defer globalLock.Unlock()

	fn.execute()
}

// 新建一个锁对象
func (fns *InternalFunctionSet) NewLock() Value {
	obj := &QKLock{&sync.Mutex{}}
	return newClass("QKLock", &obj)
}

type QKLock struct {
	mux *sync.Mutex
}

func (lk QKLock) Lock() {
	lk.mux.Lock()
}

func (lk QKLock) Unlock() {
	lk.mux.Unlock()
}
