package core


// 字符串列表targets中是否存在指定字符串src
func match(src string, targets ...string) bool {
    for _, target := range targets {
        if src == target {
            return true
        }
    }
    return false
}

// 把一个Token转成Token数组
func tokenArray(tk Token) []Token {
	var res []Token
	res = append(res, tk)
	return res
}

// 在token列表头部插入一个新token
func insert(h Token, ts []Token) []Token {
    res := make([]Token, 0, len(ts)+1)
    res = append(res, h)
    for _, t := range ts {
        res = append(res, t)
    }
    return res
}

// 移除Token列表中, 最后一个Token
func removeTailToken(ts []Token) {
	if len(ts) < 1 {
		return
	}
	ts = ts[:len(ts)-1]
}

// 获取当前索引于当前token列表的上一个token， 并返回上一个token是否存在的判断
func preToken(currentIndex int, ts []Token) (t Token, ok bool) {
    if currentIndex-1 < 0 {
        return
    }
    return ts[currentIndex-1], true
}

// token列表的最后一个token， 并返回最后一个token是否存在的判断
func lastToken(ts []Token) (t Token, ok bool) {
    size := len(ts)
    if size < 1 {
        return
    }
    return ts[size-1], true
}

// token列表的倒数第二个token， 并返回倒数第二个token是否存在的判断
func lastSecondToken(ts []Token) (t Token, ok bool) {
	size := len(ts)
	if size < 2 {
		return
	}
	return ts[size-2], true
}

// 获取当前索引于当前token列表的下一个token， 并返回下一个token是否存在的判断
func nextToken(currentIndex int, ts []Token) (t Token, ok bool) {
    if currentIndex+1>=len(ts) {
        return
    }
    return ts[currentIndex+1], true
}

// token列表的倒数第二个token
func lastSecond(ts []Token) Token {
	return ts[len(ts)-2]
}
// token列表的最后一个token
func last(ts []Token) Token {
    return ts[len(ts)-1]
}

// 获取当前索引于当前token列表的下一个token
func next(ts []Token, i int) Token {
    return ts[i+1]
}

// 判断token列表中是否存在指定符号token
func hasSymbol(ts []Token, ss ...string) bool {
    for i:=0; i<len(ts); i++ {
        t := ts[i]
        for _, s := range ss {
            if t.assertSymbol(s) {
                return true
            }
        }
    }
    return false
}

// 获取指定符号token列表任意Token的下一个索引
func nextSymbolsIndex(ts []Token, currentIndex int, tokenStrs ...string) int {
	for i:=currentIndex; i<len(ts); i++ {
		t := ts[i]
		if t.assertSymbols(tokenStrs...) {
			return i
		}
	}
	return -1
}

// 获取指定符号token的下一个索引
func nextSymbolIndex(ts []Token, currentIndex int, s string) int {
    for i:=currentIndex; i<len(ts); i++ {
        t := ts[i]
        if t.assertSymbol(s) {
            return i
        }
    }
    return -1
}

// 获取指定符号token的下一个索引, 其间或错误符号直接退出
func nextSymbolIndexNotError(ts []Token, currentIndex int, s string, errSymbols ...string) int {
	for i:=currentIndex; i<len(ts); i++ {
		t := ts[i]
		if t.assertSymbols(errSymbols...) {
			return -1
		}
		if t.assertSymbol(s) {
			return i
		}
	}
	return -1
}

// 根据指定分隔符获取程序块的尾索引
func scopeEndIndex(ts []Token, currentIndex int, open, close string) int {
	scopeOpenCount := 0
	size := len(ts)
	for i:=currentIndex; i<size; i++ {
		t := ts[i]
		if t.assertSymbol(open) {
			scopeOpenCount++
		}
		if t.assertSymbol(close) {
			scopeOpenCount--
			if scopeOpenCount == 0 {
				return i
			}
		}
	}
	if scopeOpenCount > 0 {
		msg := printCurrentPositionTokens(ts, currentIndex)
		runtimeExcption("scopeEndIndex: no match final character \""+close+"\"", msg)
	}
	return -1
}

// 消除两边小括号token
func clearParentheses(ts []Token) []Token {
    size := len(ts)
    if size >= 3 && ts[0].assertSymbol("(") && ts[size-1].assertSymbol(")") {
        ts = ts[1 : size-1]
    }
    return ts
}

// 消除两边中括号token
func clearBrackets(ts []Token) []Token {
	size := len(ts)
	if size >= 3 && ts[0].assertSymbol("[") && ts[size-1].assertSymbol("]") {
		ts = ts[1 : size-1]
	}
	return ts
}

// 消除两边大括号token
func clearBraces(ts []Token) []Token {
    size := len(ts)
    if size >= 3 && ts[0].assertSymbol("{") && ts[size-1].assertSymbol("}") {
        ts = ts[1 : size-1]
    }
    return ts
}





