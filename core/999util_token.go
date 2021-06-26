package core



func match(src string, targets ...string) bool {
    for _, target := range targets {
        if src == target {
            return true
        }
    }
    return false
}

func insert(h Token, ts []Token) []Token {
    res := make([]Token, 0, len(ts)+1)
    res = append(res, h)
    for _, t := range ts {
        res = append(res, t)
    }
    return res
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

func lastSecondToken(ts []Token) (t Token, ok bool) {
	size := len(ts)
	if size < 2 {
		return
	}
	return ts[size-2], true
}

func nextToken(currentIndex int, ts []Token) (t Token, ok bool) {
    if currentIndex+1>=len(ts) {
        return
    }
    return ts[currentIndex+1], true
}


func last(ts []Token) *Token {
    return &ts[len(ts)-1]
}

func next(ts []Token, i int) *Token {
    return &ts[i+1]
}


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

func nextSymbolIndex(ts []Token, currentIndex int, s string) int {
    for i:=currentIndex; i<len(ts); i++ {
        t := ts[i]
        if t.assertSymbol(s) {
            return i
        }
    }
    return -1
}

func scopeEndIndex(ts []Token, currentIndex int, open, close string) int {
	scopeOpenCount := 0
	size := len(ts)
	for i:=currentIndex; i<size; i++ {
		t := ts[i]
		if t.assertSymbol("{") {
			scopeOpenCount++
		}
		if t.assertSymbol("}") {
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
	if size >= 3 && ts[0].assertSymbol("(") && ts[size-1].assertSymbol(")") {
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





