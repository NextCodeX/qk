package core

type SymbolToken struct {
	TokenAdapter
}

func newSymbolToken(raw string, row int) Token {
	t := &SymbolToken{}
	t.lineIndex = row
	t.val = raw
	t.typName = "Symbol"
	return t
}

func (t *SymbolToken) get() string {
	return t.val.(string)
}
func (t *SymbolToken) set(s string) {
	t.val = s
}

func (t *SymbolToken) isSymbol() bool {
	return true
}

func (t *SymbolToken) assertSymbol(s string) bool {
	return t.val.(string) == s
}

func (t *SymbolToken) assertSymbols(ss ...string) bool {
	self := t.val.(string)
	for _, s := range ss {
		if s == self {
			return true
		}
	}
	return false
}

// 获取运算符优先级
// （注：运算符的优先级，值越小，优先级越高）
func (t *SymbolToken) priority() int {
	switch t.val.(string) {
	//case "(", ")", "[","]", ".":
	//    return 1
	//case "!", "+", "-", "++", "--":
	//	  ! +(正)  -(负)   ++ -- , 结合性：从右向左
	//	  return 2
	case "*", "/", "%":
		return 3
	case "+", "-":
		// +(加) -(减)
		return 4
	case "<<", ">>", ">>>":
		return 5
	case "<", "<=", ">", ">=":
		return 6
	case "==", "!=":
		return 7
	case "&":
		// (按位与)
		return 8
	case "^":
		return 9
	case "|":
		return 10
	case "&&":
		return 11
	case "||":
		return 12
	case "?:":
		//  结合性：从右向左
		return 13
	case "=", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", " =", "<<=", ">>=", ">>>=":
		// 结合性：从右向左
		return 14
	}
	return -1
}
