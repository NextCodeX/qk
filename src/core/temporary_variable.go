package core

import (
	"fmt"
	"strings"
)

// 用于统计临时变量名,或计算得到临时变量名
var tmpcount int
// 临时变量前缀
var tmpPrefix = "tmp#"
// 函数调用返回值名
var funcResultName = "#result"

func getTmpVarToken() Token {
	tmpname := getTmpname()
	return varToken(tmpname)
}

func getTmpname() string {
	name := fmt.Sprintf("%v%v", tmpPrefix, tmpcount)
	tmpcount++
	return name
}

func isTmpVar(name string) bool {
	return strings.HasPrefix(name, tmpPrefix)
}
