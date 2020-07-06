package core

import (
    "fmt"
    "io/ioutil"
    "os/exec"
    "runtime"
    "strings"
)

func Run() {
    bs, _ := ioutil.ReadFile("demo.qk")
    ts := Parse(bs)
    //printTokens(ts)
    Compile(mainFunc, ts)
    printFunc()
}

func printFunc() {
    doPrintFunc(mainFunc)
    for _, fn := range funcList {
        doPrintFunc(fn)
    }
}

func doPrintFunc(fn *Function) {
    for i, stmt := range fn.block {
        fmt.Printf("line %v: \n %v \n", i, stmt)
    }
}

func printTokens(tokens []Token) {
    for i, token := range tokens {
        fmt.Printf("line %v: [%v] -> %v \n", i, token.String(), token.TokenTypeName())
    }
}

// 获取命令所在的路径
func getCmdDir() string {
    cmd := exec.Command( "cmd", "/c", "cd")
    if runtime.GOOS != "windows" {
        cmd = exec.Command("pwd")
    }
    d, err := cmd.CombinedOutput()
    exitOnError(err)
    pwd := strings.TrimSpace(string(d))
    return pwd
}


