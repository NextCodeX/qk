package core

import (
	"fmt"
	"os"
	"bytes"
)


func assert(flag bool, msg ...interface{})  {
	if flag {
		runtimeExcption(msg)
	}
}

func runtimeExcption(raw ...interface{}){
	var msg []interface{}
	for _, item := range raw {
		if err, ok := item.(error); ok && err != nil {
			msg = append(msg, err.Error())
			continue
		}
		msg = append(msg, item)
	}
	if DEBUG_MODE {
		panic(fmt.Sprintln(msg...))
		return
	}
	fmt.Println(msg...)
	os.Exit(2)
}

func printExprTokens(exprTokensList [][]Token) {
	var buf bytes.Buffer
	for _, ts := range exprTokensList {
		buf.WriteString(tokensString(ts))
		buf.WriteString("\n")
	}
	fmt.Println(buf.String())
}


func printFunc() {
	doPrintFunc(mainFunc)
	for _, fn := range funcList {
		doPrintFunc(fn)
	}
}

func doPrintFunc(fn *Function) {
	fmt.Println("######################", fn.name, len(fn.block))
	for i, stmt := range fn.block {
		fmt.Printf("num: %v line %v: \n %v \n", len(stmt.raw), i, stmt)
	}
}

func printTokensByLine(tokens []Token) {
	for i, token := range tokens {
		fmt.Printf("count %v-%v: [%v] -> %v \n", i, token.lineIndexString(), token.String(), token.TokenTypeName())
	}
}