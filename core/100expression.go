package core

type Expression interface {
    raw() []Token  // 获取表达式原始Token列表
    setRaw(ts []Token) // 设置表达式原始Token列表

    setStack(stack Function) // 设置表达式的stack
    getStack() Function // 获取表达式的stack

    execute() Value // 执行表达式
}