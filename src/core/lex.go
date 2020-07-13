package core

import (
    "fmt"
)



func ParseTokens(bs []byte) []Token {
    // 提取原始token列表
    ts := parse4PrimaryTokens(bs)
    // 语法预处理
    //ts := syntaxPreHandle(ts)
    // 提取'++', '--'等运算符
    ts = parse4OperatorTokens(ts)
    // 去掉无用的';', 合并token生成函数调用token(Fcall), 方法调用token(Mtcall)等复合token
    ts = parse4ComplexTokens(ts)
    return ts
}

// 提取多符号运算符(>=, <=...)
func parse4OperatorTokens(ts []Token) []Token {
    var res []Token
    for _, token := range ts {
        last, lastExist := lastToken(res)

        currentIsEqual := token.assertSymbol("=")
        condEqualMerge := lastExist && last.assertSymbols("=", ">", "<", "+", "-", "*", "/", "%")
        condEqual := currentIsEqual && condEqualMerge

        currentIsOr := token.assertSymbol("|")
        condOrMerge := lastExist && last.assertSymbols("|")
        condOr := currentIsOr && condOrMerge

        currentIsAnd := token.assertSymbol("&")
        condAndMerge := lastExist && last.assertSymbols("&")
        condAnd := currentIsAnd && condAndMerge

        currentIsAdd := token.assertSymbol("+")
        condAddMerge := lastExist && last.assertSymbols("+")
        condAdd := currentIsAdd && condAddMerge

        currentIsSub := token.assertSymbol("-")
        condSubMerge := lastExist && last.assertSymbols("-")
        condSub := currentIsSub && condSubMerge

        if condEqual || condAnd || condOr || condAdd || condSub {
            res = tailMerge(res, token)
            if (condAdd || condSub) && canAddSubSelf(res) {
                newTokenType := getTokenType(condAdd)
                res = tailMerge2(res, newTokenType)
            }
            continue
        }

        res = append(res, token)
    }
    return res
}

func tailMerge2(ts []Token, tokenType TokenType) []Token {
    size := len(ts)

    tailIndex := size - 2
    tail := ts[tailIndex]
    tail.t = tail.t | tokenType

    res := ts[:size-1]
    res[tailIndex] = tail
    return res
}

func getTokenType(condAdd bool) TokenType {
    var newTokenType TokenType
    if condAdd {
        newTokenType = AddSelf
    } else {
        newTokenType = SubSelf
    }
    return newTokenType
}

func canAddSubSelf(ts []Token) bool {
    size := len(ts)
    if size>=2 && ts[size-2].isIdentifier() {
        return true
    }
    return false
}

func tailMerge(ts []Token, t Token) []Token {
    size := len(ts)
    tail := ts[size-1]
    tail.str = fmt.Sprintf("%v%v", tail.str, t.str)
    ts[size - 1] = tail
    return ts
}

// 该函数用于： 去掉无用的';', 合并token生成函数调用token(Fcall), 方法调用token(Mtcall)等复合token
func parse4ComplexTokens(ts []Token) []Token {
    var res []Token
    size := len(ts)
    for i:=0; i<size; {
        token := ts[i]
        pre, preExist := preToken(i, ts)
        next, nextExist := nextToken(i, ts)
        var t Token
        var nextIndex int

        // 处理无用分号
        if token.str == ";" && ((preExist && pre.assertSymbols("{","}")) || (nextExist && next.assertSymbols("{", ";"))) {
            goto end_current_iterate
        }

        // 捕获数组的字面值Token
        if token.assertSymbol("[") && preExist && pre.assertSymbols("=", "(") {
            t, nextIndex = extractArrayLiteral(i+1, ts)
            if nextIndex > i {
                res = append(res, t)
                i = nextIndex
                goto next_loop
            }
        }
        // 捕获对象的字面值Token
        if token.assertSymbol("{") && preExist && pre.assertSymbols("=", "(") {
            t, nextIndex = extractObjectLiteral(i+1, ts)
            if nextIndex > i {
                res = append(res, t)
                res = append(res, symbolToken(";"))
                i = nextIndex
                goto next_loop
            }
        }

        if !token.isIdentifier() || !nextExist {
            goto token_collect
        }

        // 捕获Attribute类型token
        t, nextIndex = extractAttribute(i, ts)
        if nextIndex > i {
            // 捕获Mtcall类型token
            if nextIndex < size && ts[nextIndex].assertSymbol("(") {
                t, nextIndex = extractMethodCall(i, ts)
            }
            res = append(res, t)
            i = nextIndex
            goto next_loop
        }

        // 捕获Fcall类型token
        t, nextIndex = extractFunctionCall(i, ts)
        if nextIndex > i {
            // 标记Fdef类型token
            if nextIndex < size && ts[nextIndex].assertSymbol("{") {
                t.t = Fdef | t.t
            }
            res = append(res, t)
            i = nextIndex
            goto next_loop
        }

        // 捕获Element类型token
        t, nextIndex = extractElement(i, ts)
        if nextIndex > i {
            res = append(res, t)
            i = nextIndex
            goto next_loop
        }

        // token 原样返回
        token_collect:
        res = append(res, token)

        end_current_iterate:
        i++
        next_loop:
    }
    return res
}

func extractArrayLiteral(currentIndex int, ts []Token) (t Token, nextIndex int) {
    size := len(ts)
    scopeOpenCount := 1
    var elems []Token
    for i := currentIndex; i < size; i++ {
        token := ts[i]
        if token.assertSymbol("]") {
            scopeOpenCount --
            nextIndex = i + 1
            break
        }
        if token.isSymbol() && !match(token.str, ",") {
            msg := printCurrentPositionTokens(ts, i)
            panic("extract ArrayLiteral Exception, illegal character:" + msg)
        }
        elems = append(elems, token)
    }
    if scopeOpenCount > 0 {
        panic("extract ArrayLiteral Exception: no match final character \"]\"")
    }
    t = Token{
        str:    "[]",
        t:      ArrLiteral | Complex,
        ts:     elems,
    }
    return t, nextIndex
}

func extractObjectLiteral(currentIndex int, ts []Token) (t Token, nextIndex int) {
    size := len(ts)
    scopeOpenCount := 1
    var elems []Token
    for i := currentIndex; i < size; i++ {
        token := ts[i]
        if token.assertSymbol("{") {
            scopeOpenCount ++
        }
        if token.assertSymbol("}") {
            scopeOpenCount --
            if scopeOpenCount == 0 {
                nextIndex = i + 1
                break
            }
        }
        if token.isSymbol() && !match(token.str,",", ":", "[", "]", "{", "}") {
            msg := printCurrentPositionTokens(ts, i)
            panic("extract element ObjectLiteral, illegal character: " + msg + " -type " + token.TokenTypeName())
        }
        elems = append(elems, token)
    }
    if scopeOpenCount > 0 {
        panic("extract element ObjectLiteral: no match final character \"}\"")
    }
    t = Token{
        str:    "{}",
        t:      ObjLiteral | Complex,
        ts:     elems,
    }
    return t, nextIndex
}

func extractElement(currentIndex int, ts []Token) (t Token, nextIndex int) {
    size := len(ts)
    // 检测不符合元素定义直接返回
    if size - currentIndex < 3 || !ts[currentIndex+1].assertSymbol("[") {
        return
    }
    var indexs []Token
    extractElementIndexTokens(currentIndex+2, ts, &nextIndex, &indexs)

    t = Token {
        str:    ts[currentIndex].str,
        t:      Element | Complex,
        ts:     indexs,
    }
    return t, nextIndex
}

func extractElementIndexTokens(currentIndex int, ts []Token, nextIndex *int, indexs *[]Token) {
    size := len(ts)
    scopeOpenCount := 1
    for i := currentIndex; i < size; i++ {
        token := ts[i]
        if token.assertSymbol("]") {
            scopeOpenCount --
            *nextIndex = i + 1
            break
        }
        if isIllegalElementIndexToken(token) {
            panic("extract element index Exception, illegal character:"+token.str)
        }
        *indexs = append(*indexs, token)
    }
    if scopeOpenCount > 0 {
        panic("extract element index Exception: no match final character \"]\"")
    }
    if *nextIndex < size && ts[*nextIndex].assertSymbol("[") {
        *indexs = append(*indexs, symbolToken(","))
        extractElementIndexTokens(*nextIndex+1, ts, nextIndex, indexs)
    }
}

func extractFunctionCall(currentIndex int, ts []Token) (t Token, nextIndex int) {
    size := len(ts)
    // 检测不符合函数调用定义直接返回
    if size - currentIndex < 3 || !ts[currentIndex+1].assertSymbol("(") {
        return
    }

    args, nextIndex := getCallArgsTokens(currentIndex + 2, ts)

    t = Token{
        str:    ts[currentIndex].str,
        t:      Fcall | Complex,
        ts:     args,
    }
    return t, nextIndex
}

func extractMethodCall(currentIndex int, ts []Token) (t Token, nextIndex int) {
    args, nextIndex := getCallArgsTokens(currentIndex + 4, ts)

    t = Token{
        str:    ts[currentIndex+2].str,
        t:      Mtcall | Complex,
        caller: ts[currentIndex].str,
        ts:     args,
    }
    return t, nextIndex
}

func getCallArgsTokens(currentIndex int, ts []Token) (args []Token, nextIndex int) {
    size := len(ts)
    scopeOpenCount := 1
    for i := currentIndex; i < size; i++ {
        token := ts[i]
        if token.assertSymbol("(") {
            scopeOpenCount ++
        }
        if token.assertSymbol(")") {
            scopeOpenCount --
            if scopeOpenCount == 0 {
                nextIndex = i + 1
                break
            }
        }
        if isIllegalFcallArgsToken(token) {
            msg := printCurrentPositionTokens(ts, i)
            panic("extract call args Exception, illegal character:"+msg)
        }
        args = append(args, token)
    }
    if scopeOpenCount > 0 {
        panic("extract call args Exception: no match final character \")\"")
    }
    return args, nextIndex
}

// 元素索引里的非法符号
// 函数，方法调用时参数列表里的非法符号
func isIllegalElementIndexToken(t Token) bool {
    if !t.isSymbol() {
        return false
    }
    switch t.str {
    case "{", "}", ",", ";", "[", "=":
        return true
    }
    return false
}

// 函数，方法调用时参数列表里的非法符号
func isIllegalFcallArgsToken(t Token) bool {
    if !t.isSymbol() {
        return false
    }
    switch t.str {
    case "{", "}", ";", "=":
        return true
    }
    return false
}

func extractAttribute(currentIndex int, ts []Token) (t Token, nextIndex int) {
    size := len(ts)
    if size - currentIndex < 3 {
        return
    }
    if !ts[currentIndex+1].assertSymbol(".")  || !ts[currentIndex+2].isIdentifier() {
        return
    }

    token := Token{
        str:    ts[currentIndex+2].str,
        t:      Attribute | Complex,
        caller: ts[currentIndex].str,
    }
    return token, currentIndex+3
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

func parse4PrimaryTokens(bs []byte) []Token {
    var tokens []Token
    var tmp []byte
    state := stateNormal
    for _, b := range bs {

        if state == stateStrLiteral && b != '"' {
            tmp = append(tmp, b)
            continue
        }

        switch {
        case (b>='a' && b<='z') || (b>='A' && b<='Z') || b=='_':
            tmp = append(tmp, b)
            state = stateStr
        case b>='0' && b<='9':
            tmp = append(tmp, b)
            if state == stateDot {
                state = stateFloat
            }else{
                state = stateInt
            }

        case b==' ' || b=='\t' || b =='\n':
            longTokenSave(b, state, &tmp, &tokens)

            if b=='\n' {
                addBoundry(&tokens)
            }
            state = stateSpace

        case isSymbol(b):
            if b == '.' && state==stateInt {
                tmp = append(tmp, b)
                state = stateDot
                break
            }

            longTokenSave(b, state, &tmp, &tokens)

            symbol := symbolToken(string(b))
            last, lastExist := lastToken(tokens)
            if symbol.assertSymbol("}") && lastExist && last.assertSymbol(";") {
                // 去掉无用的";"
                tokens[len(tokens)-1] = symbol
            } else {
                tokens = append(tokens, symbol)
            }
            state = stateSymbol

        case b == '"':
            if len(tmp) < 1 {
                state = stateStrLiteral
            } else {
                if tmp[len(tmp)-1] != '\\' {
                    longTokenSave(b, state, &tmp, &tokens)
                    state = stateNormal
                } else {
                    tmp = append(tmp, b)
                }
            }

        }
    }
    longTokenSave('\n', state, &tmp, &tokens)
    addBoundry(&tokens)

    return tokens
}

func addBoundry(ts *[]Token) {
    size := len(*ts)
    if size>0 && (*ts)[size-1].assertSymbols("{", ",", "}") {
        // 防止添加无用的";"
        return
    }

    *ts = append(*ts, Token{
        str: ";",
        t:   Symbol,
    })
}

func longTokenSave(b byte, state int, tmp *[]byte, tokens *[]Token) {
    s := string(*tmp)
    if len(s) < 1 {
        return
    }

    var tokenType TokenType
    if state == stateFloat {
        tokenType = Float
    }
    if state == stateInt && b != '.' {
        tokenType = Int
    }
    if state == stateStr {
        tokenType = Identifier
    }
    if state == stateStrLiteral {
        tokenType = Str
    }
    *tokens = append(*tokens, Token{
        str: s,
        t:   tokenType,
    })
    // 重置临时变量
    *tmp = nil
}

func isSymbol(b byte) bool {
    switch b {
    case '.': fallthrough
    case ':': fallthrough
    case '(': fallthrough
    case ')': fallthrough
    case '[': fallthrough
    case ']': fallthrough
    case '{': fallthrough
    case '}': fallthrough
    case ';': fallthrough
    case ',': fallthrough
    case '=': fallthrough
    case '+': fallthrough
    case '-': fallthrough
    case '*': fallthrough
    case '/': fallthrough
    case '%': fallthrough
    case '>': fallthrough
    case '<':
        return true
    }
    return false
}





