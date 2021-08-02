package core

// 用于实现多种形式的方法调用。
type Object interface {
	get(key string) Value
}
