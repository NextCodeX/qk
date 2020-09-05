package core

// 数据库连接接口
type DBConnection interface {
	Insert(sql string, args...interface{}) interface{}
	Update(sql string, args...interface{}) int64
	GetValue(sql string, args...interface{}) interface{}
	GetRow(sql string, args...interface{}) map[string]interface{}
	GetRows(sql string, args...interface{}) []map[string]interface{}
}
