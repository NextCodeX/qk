package core

import "fmt"

func match(src string, targets ...string) bool {
    for _, target := range targets {
        if src == target {
            return true
        }
    }
    return false
}

func exitOnError(err error, a ...string) {
    if err != nil {
        fmt.Println(a)
        panic(err)
    }
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
