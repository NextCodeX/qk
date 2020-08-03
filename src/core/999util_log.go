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
