package core

import (
    "fmt"
    "io/ioutil"
    "os/exec"
    "runtime"
    "strings"
)

var (
    funcList = make(map[string]*Function)
    mainFunc = newFunc("main")
)

const DEBUG_MODE = true

func Run() {
    //qkfile := "examples/demo.qk"
    //qkfile := "examples/expr.qk"
    //qkfile := "examples/if-stmt.qk"
    //qkfile := "examples/for-stmt2.qk"
    //qkfile := "examples/foreach-stmt.qk"
    //qkfile := "examples/foritem-stmt.qk"
    //qkfile := "examples/forindex-stmt.qk"
    qkfile := "examples/func.qk"

    bs, _ := ioutil.ReadFile(qkfile)

    // 词法分析
    ts := ParseTokens(bs)
    printTokensByLine(ts)

    // 语法分析
    mainFunc.raw = ts
    Compile(mainFunc)
    printFunc()

    // 解析并执行
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
        fmt.Printf("count %v-%v: [%v] -> %v \n", i, token.lineIndexString(), token.String(), token.TokenTypeName())
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


