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
	Int // 整数类型
	Float  // 浮点类型
	Symbol  // 符号

	Fcall  // 函数调用
	Fdef  // 函数 定义
	Mtcall // 方法调用
	Attribute // 对象属性
	ArrLiteral // 数组字面值
	ObjLiteral // 对象字面值
	Element // 元素，用于指示对象或数组的元素
	Complex // 用于标记复合类型token

	AddSelf // 自增一
	SubSelf // 自减一
)

type Token struct {
	lineIndex int // token的首行索引
	endLineIndex int // 当token跨行时, 存的尾行索引
	str string // token字符串值
	t TokenType // 类型
	caller string // token为方法, 属性类型时才有的调用者变量名
	// token为元素, 函数调用, 函数定义类型时, 存的是参数   ;
	// token为数组字面值, 对象字面值类型时, 存的是字面值内容
	ts []Token
}

func newToken(raw string, t TokenType) Token {
	return Token{str:raw,t:t}
}

func symbolToken(s string) Token {
	return Token{str:s, t:Symbol}
}

func symbolTokenWithLineIndex(s string, lineIndex int) Token {
	return Token{str:s, t:Symbol, lineIndex:lineIndex}
}

func varToken(s string) Token {
	return Token{str:s, t:Identifier}
}

func (this *Token) isIdentifier() bool {
	return (this.t & Identifier) == Identifier
}

func (this *Token) isStr() bool {
	return (this.t & Str) == Str
}

func (this *Token) isInt() bool {
	return (this.t & Int) == Int
}

func (this *Token) isFloat() bool {
	return (this.t & Float) == Float
}

func (this *Token) isSymbol() bool {
	return (this.t & Symbol) == Symbol
}

func (this *Token) isFdef() bool {
	return (this.t & Fdef) == Fdef
}

func (this *Token) isFcall() bool {
	return (this.t & Fcall) == Fcall
}

func (this *Token) isAttribute() bool {
	return (this.t & Attribute) == Attribute
}

func (this *Token) isMtcall() bool {
	return (this.t & Mtcall) == Mtcall
}

func (this *Token) isArrLiteral() bool {
	return (this.t & ArrLiteral) == ArrLiteral
}

func (this *Token) isObjLiteral() bool {
	return (this.t & ObjLiteral) == ObjLiteral
}

func (this *Token) isElement() bool {
	return (this.t & Element) == Element
}

func (this *Token) isComplex() bool {
	return (this.t & Complex) == Complex
}

func (this *Token) isAddSelf() bool {
	return (this.t & AddSelf) == AddSelf
}

func (this *Token) isSubSelf() bool {
	return (this.t & SubSelf) == SubSelf
}

func (this *Token) assertSymbol(s string) bool {
	return this.isSymbol() && this.str == s
}

func (this *Token) assertSymbols(ss ...string) bool {
	if !this.isSymbol(){
		return false
	}
	for _, s := range ss {
		if s == this.str {
			return true
		}
	}
	return false
}

// 获取运算符优先级
// （注：运算符的优先级，值越小，优先级越高）
func (this *Token) priority() int {
	res := -1

	if !this.isSymbol() {
		return res
	}

	switch {
	//case this.assertSymbols("(", ")", "[","]", "."):
	//    res = 1
	//case this.assertSymbols("!", "+", "-", " ", "++", "--"):
	//! +(正)  -(负)   ++ -- , 结合性：从右向左
	//res = 2
	case this.assertSymbols("*", "/", "%"):
		res = 3
	case this.assertSymbols("+", "-"):
		// +(加) -(减)
		res = 4
	case this.assertSymbols("<<", ">>", ">>>"):
		res = 5
	case this.assertSymbols("<", "<=", ">", ">="):
		res = 6
	case this.assertSymbols("==", "!="):
		res = 7
	case this.assertSymbols("&"):
		// (按位与)
		res = 8
	case this.assertSymbols("^"):
		res = 9
	case this.assertSymbols("|"):
		res = 10
	case this.assertSymbols("&&"):
		res = 11
	case this.assertSymbols("||"):
		res = 12
	case this.assertSymbols("?:"):
		//  结合性：从右向左
		res = 13
	case this.assertSymbols("=", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", " =", "<<=", ">>=", ">>>="):
		// 结合性：从右向左
		res = 14
	}
	return res
}

func isValidPriorityCompared(t1, t2 *Token) bool {
	if t1.priority() == -1 || t2.priority() == -1 {
		return false
	}
	return true
}

func (this *Token) equal(t *Token) bool {
	return isValidPriorityCompared(this,t) && this.priority() == t.priority()
}

func (this *Token) lower(t *Token) bool {
	return isValidPriorityCompared(this,t) && this.priority() < t.priority()
}

func (this *Token) upper(t *Token) bool {
	return isValidPriorityCompared(this,t) && this.priority() > t.priority()
}

func (this *Token) String() string {
	if this.isArrLiteral() || this.isObjLiteral() {
		return this.toJSONString()
	}

	if this.isFcall() || this.isFdef() {
		var buf bytes.Buffer
		buf.WriteString(this.str)
		buf.WriteString("(")
		if this.ts != nil {
			for _, token := range this.ts {
				buf.WriteString(token.String())
			}
		}
		buf.WriteString(")")
		return buf.String()
	}

	if this.isAttribute() || this.isMtcall() {
		var buf bytes.Buffer
		buf.WriteString(this.caller)
		buf.WriteString(".")
		buf.WriteString(this.str)
		if this.isMtcall() {
			buf.WriteString("(")
			if this.ts != nil {
				for _, token := range this.ts {
					buf.WriteString(token.String())
				}
			}
			buf.WriteString(")")
		}
		return buf.String()
	}

	if this.isElement() {
		var buf bytes.Buffer
		buf.WriteString(this.str)
		buf.WriteString("[")
		if this.ts != nil {
			for _, token := range this.ts {
				buf.WriteString(token.String())
			}
		}
		buf.WriteString("]")
		return buf.String()
	}

	if this.isStr() {
		return fmt.Sprintf(`"%v"`, this.str)
	}

	if this.isAddSelf() {
		return fmt.Sprintf(`%v ++`, this.str)
	}

	if this.isSubSelf() {
		return fmt.Sprintf(`%v --`, this.str)
	}

	return this.str
}

func (this *Token) toJSONString() string {
	if this.isArrLiteral() {
		var buf bytes.Buffer
		buf.WriteString("[")
		if this.ts != nil {
			for _, token := range this.ts {
				if token.isStr() {
					buf.WriteString(fmt.Sprintf(`"%v"`, token.str))
				} else {
					buf.WriteString(token.str)
				}
			}
		}
		buf.WriteString("]")
		return buf.String()
	}

	if this.isObjLiteral() {
		var buf bytes.Buffer
		buf.WriteString("{")
		if this.ts != nil {
			for _, token := range this.ts {
				if token.isStr() || token.isIdentifier() {
					buf.WriteString(fmt.Sprintf(`"%v"`, token.str))
				} else {
					buf.WriteString(token.str)
				}
			}
		}
		buf.WriteString("}")
		return buf.String()
	}
	return ""
}

func (t *Token) TokenTypeName() string {
	var buf bytes.Buffer
	if t.isStr() {
		buf.WriteString( "string, ")
	}
	if t.isIdentifier() {
		buf.WriteString( "identifier, ")
	}
	if t.isInt() {
		buf.WriteString( "int, ")
	}
	if t.isFloat() {
		buf.WriteString( "float, ")
	}
	if t.isSymbol() {
		buf.WriteString( "symbol, ")
	}
	if t.isFdef() {
		buf.WriteString("function define, ")
	}
	if t.isFcall() {
		buf.WriteString("function call, ")
	}
	if t.isMtcall() {
		buf.WriteString("method call, ")
	}
	if t.isAttribute() {
		buf.WriteString("attribute, ")
	}
	if t.isArrLiteral() {
		buf.WriteString("array literal, ")
	}
	if t.isObjLiteral() {
		buf.WriteString("object literal, ")
	}
	if t.isElement() {
		buf.WriteString("element, ")
	}
	if t.isComplex() {
		buf.WriteString("complex, ")
	}

	if t.isAddSelf() {
		buf.WriteString("addself, ")
	}
	if t.isSubSelf() {
		buf.WriteString("subself, ")
	}
	if buf.Len() == 0 {
		return "undefined"
	}
	return strings.TrimRight(strings.TrimSpace(buf.String()), ",")
}

func (t *Token) lineIndexString() string {
	var res bytes.Buffer
	res.WriteString(fmt.Sprintf("line: %v", t.lineIndex))
	if t.endLineIndex > t.lineIndex {
		res.WriteString(fmt.Sprintf(", %v", t.endLineIndex))
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

func tokensString(ts []Token) string {
	var buf bytes.Buffer
	for _, t := range ts {
		buf.WriteString(t.String() + "  ")
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
