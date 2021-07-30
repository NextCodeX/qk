package core

import (
	"bytes"
	"fmt"
	"strings"
)

type TokenType int
const (
	Identifier  TokenType = 1 << iota // 标识符
	Str // 字符串类型
	DynamicStr // 动态字符串类型
	Int // 整数类型
	Float  // 浮点类型
	Symbol  // 符号

	Fcall  // function call 函数调用
	Fdef  // function definition 函数 定义
	Mtcall // method call 方法调用
	Attribute // 对象属性
	ArrLiteral // 数组字面值
	ObjLiteral // 对象字面值
	Element // 元素, 用于指示对象或数组的元素
	Complex // 用于标记复合类型token
	ChainCall // 链式调用
	SubList // 列表截取
	Expr // 表示这个Token包含了一个Expression
	Not  // 用于标记需要非处理的一元表达式

	AddSelf // 自增一
	SubSelf // 自减一
)

type Token interface {
	isIdentifier() bool
	isStr() bool
	isDynamicStr() bool
	isInt() bool
	isFloat() bool
	isSymbol() bool
	isFdef() bool
	isFcall() bool
	isAttribute() bool
	isMtcall() bool
	isArrLiteral() bool
	isObjLiteral() bool
	isElement() bool
	isComplex() bool
	isSubList() bool
	isChainCall() bool
	isExpr() bool
	isNot() bool
	isAddSelf() bool
	isSubSelf() bool
	assertIdentifier(s string) bool
	assertSymbol(s string) bool
	assertSymbols(ss ...string) bool

	getLineIndex() int
	setLineIndex(index int)
	getEndLineIndex() int
	setEndLineIndex(i int)

	raw() string
	setRaw(raw string)
	typ() TokenType
	setTyp(t TokenType)
	addType(t TokenType)
	notFlag() bool // 是否非处理
	setNotFlag(flag bool)

	priority() int
	equal(t Token) bool
	lower(t Token) bool
	upper(t Token) bool
	String() string
	toJSONString() string
	TokenTypeName() string
	lineIndexString() string

	tokens() []Token
	setTokens(ts []Token)
	tokensAppend(t Token)

	chainTokenList() []Token
	setChainTokenList(ts []Token)
	chainTokenListAppend(t Token)

	setStartExprTokens(ts []Token)
	startExprTokens() []Token

	setEndExprTokens(ts []Token)
	endExprTokens() []Token
}


type TokenImpl struct {
	lineIndex int // token的首行索引
	endLineIndex int // 当token跨行时, 存的尾行索引
	str string // 原始字符串
	t TokenType // 类型
	// token为元素, 函数调用, 函数定义类型时, 存的是参数   ;
	// token为数组字面值, 对象字面值类型时, 存的是字面值内容
	ts []Token
	chainTokens []Token
	startExpr []Token
	endExpr []Token
	not bool // 是否进行非处理
}

func newToken(raw string, t TokenType) Token {
	return &TokenImpl{str: raw,t:t}
}

func symbolToken(s string) Token {
	return &TokenImpl{str: s, t:Symbol}
}

func varToken(s string) Token {
	return &TokenImpl{str: s, t:Identifier}
}

func (tk *TokenImpl) tokens() []Token {
	return tk.ts
}
func (tk *TokenImpl) setTokens(ts []Token) {
	tk.ts = ts
}
func (tk *TokenImpl) tokensAppend(t Token) {
	tk.ts = append(tk.ts, t)
}

func (tk *TokenImpl) chainTokenList() []Token {
	return tk.chainTokens
}
func (tk *TokenImpl) setChainTokenList(ts []Token) {
	tk.chainTokens = ts
}
func (tk *TokenImpl) chainTokenListAppend(t Token) {
	tk.chainTokens = append(tk.chainTokens, t)
}

func (tk *TokenImpl) setStartExprTokens(ts []Token) {
	tk.startExpr = ts
}
func (tk *TokenImpl) startExprTokens() []Token {
	return tk.startExpr
}

func (tk *TokenImpl) setEndExprTokens(ts []Token) {
	tk.endExpr = ts
}
func (tk *TokenImpl) endExprTokens() []Token {
	return tk.endExpr
}

func (tk *TokenImpl) notFlag() bool {
	return tk.not
}
func (tk *TokenImpl) setNotFlag(flag bool) {
	tk.not = flag
}
func (tk *TokenImpl) raw() string {
	return tk.str
}
func (tk *TokenImpl) setRaw(raw string) {
	tk.str = raw
}
func (tk *TokenImpl) typ() TokenType {
	return tk.t
}
func (tk *TokenImpl) setTyp(t TokenType) {
	tk.t = t
}
func (tk *TokenImpl) addType(t TokenType) {
	tk.t = tk.t | t
}

func (tk *TokenImpl) getLineIndex() int {
	return tk.lineIndex
}
func (tk *TokenImpl) setLineIndex(index int) {
	tk.lineIndex = index
}
func (tk *TokenImpl) getEndLineIndex() int {
	return tk.endLineIndex
}
func (tk *TokenImpl) setEndLineIndex(i int) {
	tk.endLineIndex = i
}

func (tk *TokenImpl) isIdentifier() bool {
	return (tk.t & Identifier) == Identifier
}

func (tk *TokenImpl) isStr() bool {
	return (tk.t & Str) == Str
}

func (tk *TokenImpl) isDynamicStr() bool {
	return (tk.t & DynamicStr) == DynamicStr
}

func (tk *TokenImpl) isInt() bool {
	return (tk.t & Int) == Int
}

func (tk *TokenImpl) isFloat() bool {
	return (tk.t & Float) == Float
}

func (tk *TokenImpl) isSymbol() bool {
	return (tk.t & Symbol) == Symbol
}

func (tk *TokenImpl) isFdef() bool {
	return (tk.t & Fdef) == Fdef
}

func (tk *TokenImpl) isFcall() bool {
	return (tk.t & Fcall) == Fcall
}

func (tk *TokenImpl) isAttribute() bool {
	return (tk.t & Attribute) == Attribute
}

func (tk *TokenImpl) isMtcall() bool {
	return (tk.t & Mtcall) == Mtcall
}

func (tk *TokenImpl) isArrLiteral() bool {
	return (tk.t & ArrLiteral) == ArrLiteral
}

func (tk *TokenImpl) isObjLiteral() bool {
	return (tk.t & ObjLiteral) == ObjLiteral
}

func (tk *TokenImpl) isElement() bool {
	return (tk.t & Element) == Element
}

func (tk *TokenImpl) isComplex() bool {
	return (tk.t & Complex) == Complex
}

func (tk *TokenImpl) isSubList() bool {
	return (tk.t & SubList) == SubList
}

func (tk *TokenImpl) isChainCall() bool {
	return (tk.t & ChainCall) == ChainCall
}

func (tk *TokenImpl) isExpr() bool {
	return (tk.t & Expr) == Expr
}

func (tk *TokenImpl) isNot() bool {
	return (tk.t & Not) == Not
}

func (tk *TokenImpl) isAddSelf() bool {
	return (tk.t & AddSelf) == AddSelf
}

func (tk *TokenImpl) isSubSelf() bool {
	return (tk.t & SubSelf) == SubSelf
}

func (tk *TokenImpl) assertIdentifier(s string) bool {
	return tk.isIdentifier() && tk.str == s
}

func (tk *TokenImpl) assertSymbol(s string) bool {
	return tk.isSymbol() && tk.str == s
}

func (tk *TokenImpl) assertSymbols(ss ...string) bool {
	if !tk.isSymbol(){
		return false
	}
	for _, s := range ss {
		if s == tk.str {
			return true
		}
	}
	return false
}

// 获取运算符优先级
// （注：运算符的优先级，值越小，优先级越高）
func (tk *TokenImpl) priority() int {
	res := -1

	if !tk.isSymbol() {
		return res
	}

	switch {
	//case this.assertSymbols("(", ")", "[","]", "."):
	//    res = 1
	//case this.assertSymbols("!", "+", "-", "++", "--"):
	//! +(正)  -(负)   ++ -- , 结合性：从右向左
	//res = 2
	case tk.assertSymbols("*", "/", "%"):
		res = 3
	case tk.assertSymbols("+", "-"):
		// +(加) -(减)
		res = 4
	case tk.assertSymbols("<<", ">>", ">>>"):
		res = 5
	case tk.assertSymbols("<", "<=", ">", ">="):
		res = 6
	case tk.assertSymbols("==", "!="):
		res = 7
	case tk.assertSymbols("&"):
		// (按位与)
		res = 8
	case tk.assertSymbols("^"):
		res = 9
	case tk.assertSymbols("|"):
		res = 10
	case tk.assertSymbols("&&"):
		res = 11
	case tk.assertSymbols("||"):
		res = 12
	case tk.assertSymbols("?:"):
		//  结合性：从右向左
		res = 13
	case tk.assertSymbols("=", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", " =", "<<=", ">>=", ">>>="):
		// 结合性：从右向左
		res = 14
	}
	return res
}

func isValidPriorityCompared(t1, t2 Token) bool {
	if t1.priority() == -1 || t2.priority() == -1 {
		return false
	}
	return true
}

func (tk *TokenImpl) equal(t Token) bool {
	return isValidPriorityCompared(tk,t) && tk.priority() == t.priority()
}

func (tk *TokenImpl) lower(t Token) bool {
	return isValidPriorityCompared(tk,t) && tk.priority() < t.priority()
}

func (tk *TokenImpl) upper(t Token) bool {
	return isValidPriorityCompared(tk,t) && tk.priority() > t.priority()
}

func (tk *TokenImpl) String() string {
	var res string
	if tk.isChainCall() {
		var buf bytes.Buffer
		tmp := tk.t
		tk.t = ^Not & (^ChainCall) & tk.t
		buf.WriteString(tk.String())
		if tk.chainTokens != nil {
			for _, token := range tk.chainTokens {
				buf.WriteString(".")
				buf.WriteString(token.String())
			}
		}
		tk.t = tmp
		res = buf.String()
	} else if tk.isArrLiteral() || tk.isObjLiteral() {
		res = tk.toJSONString()
	} else if tk.isFcall() || tk.isFdef() {
		var buf bytes.Buffer
		buf.WriteString(tk.str)
		buf.WriteString("(")
		if tk.ts != nil {
			for _, token := range tk.ts {
				buf.WriteString(token.String())
			}
		}
		buf.WriteString(")")
		res = buf.String()
	} else if tk.isAttribute() || tk.isMtcall() {
		var buf bytes.Buffer
		//buf.WriteString(tk.caller)
		buf.WriteString(".")
		buf.WriteString(tk.str)
		if tk.isMtcall() {
			buf.WriteString("(")
			if tk.ts != nil {
				for _, token := range tk.ts {
					buf.WriteString(token.String())
				}
			}
			buf.WriteString(")")
		}
		res = buf.String()
	} else if tk.isElement() {
		var buf bytes.Buffer
		buf.WriteString(tk.str)
		buf.WriteString("[")
		if tk.ts != nil {
			for _, token := range tk.ts {
				buf.WriteString(token.String())
			}
		}
		buf.WriteString("]")
		res = buf.String()
	} else if tk.isExpr() {
		res = "("+tokensString(tk.ts)+")"
	} else if tk.isStr() {
		res = fmt.Sprintf(`"%v"`, tk.str)
	} else if tk.isAddSelf() {
		res = fmt.Sprintf(`%v ++`, tk.str)
	} else if tk.isSubSelf() {
		res = fmt.Sprintf(`%v --`, tk.str)
	} else {
		res = tk.str
	}
	if tk.isNot() {
		res = "!"+res
		if !tk.not {
			res = "!"+res
		}
	}
	return res
}

func (tk *TokenImpl) toJSONString() string {
	if tk.isArrLiteral() {
		var buf bytes.Buffer
		buf.WriteString("[")
		buf.WriteString(tokensString(tk.ts))
		buf.WriteString("]")
		return buf.String()
	}

	if tk.isObjLiteral() {
		var buf bytes.Buffer
		buf.WriteString("{")
		buf.WriteString(tokensString(tk.ts))
		buf.WriteString("}")
		return buf.String()
	}
	return ""
}

func (tk *TokenImpl) TokenTypeName() string {
	var buf bytes.Buffer
	if tk.isStr() {
		buf.WriteString( "string, ")
	}
	if tk.isIdentifier() {
		buf.WriteString( "identifier, ")
	}
	if tk.isChainCall() {
		buf.WriteString( "chain call, ")
	}
	if tk.isInt() {
		buf.WriteString( "int, ")
	}
	if tk.isFloat() {
		buf.WriteString( "float, ")
	}
	if tk.isSymbol() {
		buf.WriteString( "symbol, ")
	}
	if tk.isFdef() {
		buf.WriteString("function define, ")
	}
	if tk.isFcall() {
		buf.WriteString("function call, ")
	}
	if tk.isMtcall() {
		buf.WriteString("method call, ")
	}
	if tk.isAttribute() {
		buf.WriteString("attribute, ")
	}
	if tk.isArrLiteral() {
		buf.WriteString("array literal, ")
	}
	if tk.isObjLiteral() {
		buf.WriteString("object literal, ")
	}
	if tk.isElement() {
		buf.WriteString("element, ")
	}
	if tk.isComplex() {
		buf.WriteString("complex, ")
	}
	if tk.isExpr() {
		buf.WriteString("expression, ")
	}
	if tk.isNot() {
		buf.WriteString("not, ")
	}

	if tk.isAddSelf() {
		buf.WriteString("addself, ")
	}
	if tk.isSubSelf() {
		buf.WriteString("subself, ")
	}
	if buf.Len() == 0 {
		return "undefined"
	}
	return strings.TrimRight(strings.TrimSpace(buf.String()), ",")
}

func (tk *TokenImpl) lineIndexString() string {
	var res bytes.Buffer
	res.WriteString(fmt.Sprintf("line: %v", tk.lineIndex))
	if tk.endLineIndex > tk.lineIndex {
		res.WriteString(fmt.Sprintf(", %v", tk.endLineIndex))
	}
	return res.String()
}

func toString4Tokens(ts []Token, start, end int) string {
	var buf bytes.Buffer
	for i:=start; i<=end; i++ {
		token := ts[i]
		buf.WriteString(token.String()+" ")
	}
	return buf.String()
}

func tokensShow10(ts []Token) string {
	var buf bytes.Buffer
	for i, t := range ts {
		buf.WriteString(t.String() + " ")
		if i >= 10 {
			break
		}
	}
	return buf.String()
}

func tokensString(ts []Token) string {
	var buf bytes.Buffer
	for _, t := range ts {
		buf.WriteString(t.String() + " ")
	}
	return buf.String()
}

func printCurrentPositionTokens(ts []Token, currentIndex int) string {
	size := len(ts)
	start := 0
	if currentIndex > 10 {
		start = currentIndex - 10
	}
	end := currentIndex
	if currentIndex+1 < size {
		end = currentIndex+1
	}
	return toString4Tokens(ts, start, end)
}
