package core

import (
    "fmt"
    "os/exec"
    "runtime"
    "strings"
    "io/ioutil"
)

var (
    funcList = make(map[string]*Function)
    mainFunc = newFunc("main")
)

const DEBUG_MODE = true

func Run() {
    qkfile := "examples/string-module-test2.qk"

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

func ParseTokens(bs []byte) []Token {
    // 提取原始token列表
    ts := parse4PrimaryTokens(bs)

    // 语法预处理
    // 提取'++', '--'等运算符
    ts = parse4OperatorTokens(ts)
    // 去掉无用的';', 合并token生成函数调用token(Fcall), 方法调用token(Mtcall)等复合token
    ts = parse4ComplexTokens(ts)
    return ts
}

func Compile(stmts StatementList) {
    if stmts == nil {
        return
    }
    if stmts.isCompiled() {
        return
    }else {
        stmts.setCompiled()
    }
    extractStatement(stmts)
    parseStatementList(stmts.stmts())

    for _, customFunc := range funcList {
        Compile(customFunc)
    }
}

func Interpret() {
    stack := newVariableStack()
    stack.push()
    executeFunctionStatementList(mainFunc.block, stack)
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


