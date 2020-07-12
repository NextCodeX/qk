package core

import (
    "fmt"
    "os"
)

func match(src string, targets ...string) bool {
    for _, target := range targets {
        if src == target {
            return true
        }
    }
    return false
}

func exitOnError(err error, a ...interface{}) {
    if err != nil {
        fmt.Println(err.Error())
        fmt.Println(a...)
        os.Exit(1)
    }
}

func runtimeExcption(msg ...interface{}){
    if DEBUG_MODE {
        panic(fmt.Sprintln(msg...))
        return
    }
    fmt.Println(msg...)
    os.Exit(2)
}

func insert(h Token, ts []Token) []Token {
    res := make([]Token, 0, len(ts)+1)
    res = append(res, h)
    for _, t := range ts {
        res = append(res, t)
    }
    return res
}

func insert2(t1, t2 Token, ts []Token) []Token {
    res := make([]Token, 0, len(ts)+2)
    res = append(res, t1)
    res = append(res, t2)
    for _, t := range ts {
        res = append(res, t)
    }
    return res
}
