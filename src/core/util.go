package core

import (
    "fmt"
    "os"
    "bytes"
)

func match(src string, targets ...string) bool {
    for _, target := range targets {
        if src == target {
            return true
        }
    }
    return false
}

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

func insert(h Token, ts []Token) []Token {
    res := make([]Token, 0, len(ts)+1)
    res = append(res, h)
    for _, t := range ts {
        res = append(res, t)
    }
    return res
}

func preToken(currentIndex int, ts []Token) (t Token, ok bool) {
    if currentIndex-1 < 0 {
        return
    }
    return ts[currentIndex-1], true
}

func lastToken(ts []Token) (t Token, ok bool) {
    size := len(ts)
    if size < 1 {
        return
    }
    return ts[size-1], true
}

func nextToken(currentIndex int, ts []Token) (t Token, ok bool) {
    if currentIndex+1>=len(ts) {
        return
    }
    return ts[currentIndex+1], true
}

