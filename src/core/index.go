package core

import (
    "fmt"
    "io/ioutil"
    "os/exec"
    "runtime"
    "strings"
)

const DEBUG_MODE = true

func Run() {
    //qkfile := "demo.qk"
    qkfile := "expr.qk"
    bs, _ := ioutil.ReadFile(qkfile)
    ts := ParseTokens(bs)
    printTokensByLine(ts)
    Compile(mainFunc, ts)
    printFunc()
    fmt.Println("================")
    Interpret()
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
        fmt.Printf("count %v-line:%v: [%v] -> %v \n", i, token.lineIndex, token.String(), token.TokenTypeName())
    }
}

// 获取命令所在的路径
func getCmdDir() string {
    cmd := exec.Command( "cmd", "/c", "cd")
    if runtime.GOOS != "windows" {
        cmd = exec.Command("pwd")
    }
    d, err := cmd.CombinedOutput()
    assert(err!=nil, err, "failed to get personal work directory")
    pwd := strings.TrimSpace(string(d))
    return pwd
}


