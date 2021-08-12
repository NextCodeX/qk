package core

// 只有实现了Iterator的对象，才可以使用foreach
type Iterator interface {
	indexs() []interface{} // 索引列表
	getItem(interface{}) Value // 根据索引获取值
}
