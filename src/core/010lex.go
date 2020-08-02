package core



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