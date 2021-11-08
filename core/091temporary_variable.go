package core

import (
	"fmt"
	"strings"
)

// 用于统计临时变量名,或计算得到临时变量名
var tmpcount int
// 临时变量前缀
const tmpPrefix = "tmp#"
// 临时变量池名称
const tmpVarsKey = "tmpVars#"

func getTmpVarToken() Token {
	tmpname := getTmpname()
	return newNameToken(tmpname, -1)
}

func getTmpname() string {
	name := fmt.Sprintf("%v%v", tmpPrefix, tmpcount)
	tmpcount++
	return name
}

func isTmpVar(name string) bool {
	return strings.HasPrefix(name, tmpPrefix)
}
