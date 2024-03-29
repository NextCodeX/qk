package core

import (
	"bytes"
	"fmt"
)

func assert(flag bool, msg ...interface{}) {
	if flag {
		runtimeExcption(msg)
	}
}

// 报错并退出程序(带格式化)
func errorf(format string, args ...interface{}) {
	var msg []interface{}
	for _, item := range args {
		if err, ok := item.(error); ok && err != nil {
			msg = append(msg, err.Error())
			continue
		}
		msg = append(msg, item)
	}
	panic(fmt.Sprintf(format, msg...))
}

// 报错并退出程序(不带格式化)
func runtimeExcption(raw ...interface{}) {
	var msg []interface{}
	for _, item := range raw {
		if err, ok := item.(error); ok && err != nil {
			msg = append(msg, err.Error())
			continue
		}
		msg = append(msg, item)
	}
	panic(fmt.Sprint(msg...))
}

func printExprTokens(exprTokensList [][]Token) {
	var buf bytes.Buffer
	for _, ts := range exprTokensList {
		buf.WriteString(tokensString(ts))
		buf.WriteString("\n")
	}
	fmt.Println(buf.String())
}

func printTokensByLine(tokens []Token) {
	for i, token := range tokens {
		fmt.Printf("count %v-%v: [%v] -> %v \n", i, token.row(), token.String(), token.typeName())
	}
}

// 用于定位发生的位置
func line(tokens []Token) string {
	var start, end int
	for _, tk := range tokens {
		index := tk.rowIndex()
		if start == 0 && end == 0 {
			start, end = index, index
			continue
		}
		if index < start {
			start = index
		}
		if index > end {
			end = index
		}
	}
	return fmt.Sprintf("line%v:%v", start, end)
}

func showBytes(bs []byte) {
	for _, b := range bs {
		fmt.Printf("%08b ", b)
	}
	fmt.Println()
}

// 全局异常处理
func catch() {
	if err := recover(); err != nil {
		fmt.Println(err)
		//debug.PrintStack()

	}
}
