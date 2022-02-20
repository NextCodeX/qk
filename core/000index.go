package core

import (
	"fmt"
	"io/ioutil"
	"time"
)

var mainFunc = newMainFunction()

// 启动Quick解释器
func Start() {
	startupTime := time.Now().UnixNano()
	//qkfile := "examples/demo.qk"
	qkfile := findScriptFile()

	bs, err := ioutil.ReadFile(qkfile)
	if err != nil {
		errorf("failed to read %v;\n error info: %v", qkfile, err)
	}
	bs = ignoreFirstLine(bs)
	setRootDir(qkfile)
	Run(bs)

	if IsCost {
		duration := time.Now().UnixNano() - startupTime
		fmt.Printf("\n\nspend: %vns, %.3fms, %.3fs  \n", duration, float64(duration)/1e6, float64(duration)/1e9)
	}
}

func ignoreFirstLine(bs []byte) []byte {
	if len(bs) < 1 {
		return bs
	}
	if bs[0] == '#' {
		for i, ch := range bs {
			if ch == '\n' {
				return bs[i+1:]
			}
		}
	}
	return bs
}

// 脚本解析执行
func Run(bs []byte) {
	defer catch()

	// 词法分析
	ts := ParseTokens(bs)
	//printTokensByLine(ts)

	// 语法分析
	mainFunc.setTokenList(ts)
	Compile(mainFunc)

	// 程序执行
	initInternalVars()
	mainFunc.setInternalVars(internalVars)
	mainFunc.execute()

	// 等待所有协程执行完，再结束程序
	goroutineWaiter.Wait()
}

// 指定变量𣏾, 执行qk代码片段.
func evalScript(src string, stack Function) Value {
	ts := ParseTokens([]byte(src))
	tsLen := len(ts)
	if tsLen < 1 {
		return NULL
	}
	if last(ts).assertSymbol(";") {
		ts = ts[:tsLen-1]
	}
	if len(ts) < 1 {
		return NULL
	}
	expr := extractExpression(ts)
	expr.setParent(stack)
	qkValue := expr.execute()
	return qkValue
}

// 从字节流中提取token列表。
func ParseTokens(bs []byte) []Token {
	// 提取原始token列表(包括提取'++', '--'等运算符以及负数表达式)
	ts := parse4PrimaryTokens(bs)

	// 语法预处理
	// 去掉无用的';'
	// 提取复合token。复合token是指包含嵌套的原始表达式，比如：函数调用，JSONObject字面值
	// 这一步执行完毕，每一个非符号，非关键字Token都会对应一个PrimaryExpression
	ts = parse4ComplexTokens(ts)
	return ts
}

// 将token列表转化成可执行的go数据结构
func Compile(stmt Statement) {
	extractStatement(stmt)
	for _, stmt := range stmt.stmts() {
		stmt.parse()
	}
}
