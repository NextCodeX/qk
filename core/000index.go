package core

import (
	"fmt"
)

var mainFunc = newFuncWithoutTokens("main")

// 脚本解析执行
func Run(bs []byte) {
	defer func() {
		// 全局异常处理
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()


	// 词法分析
	ts := ParseTokens(bs)
	// printTokensByLine(ts)

	// 语法分析(解析)
	mainFunc.setRaw(ts)
	Compile(mainFunc)

	// 程序执行
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
	expr.setStack(stack)
	qkValue := expr.execute()
	return qkValue
}

// 从字节流中提取token列表。
func ParseTokens(bs []byte) []Token {
	// 提取原始token列表
	ts := parse4PrimaryTokens(bs)

	// 语法预处理
	// 提取'++', '--'等运算符以及负数表达式
	ts = parse4OperatorTokens(ts)
	// 去掉无用的';', 合并token生成函数调用token(Fcall), 方法调用token(Mtcall)等复合token
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

