package core

type MachineState int

// 状态机解析原始token用
const (
	stateIdentifier MachineState = 1 << iota
	stateStrLiteral
	stateInt
	stateDot
	stateFloat
	stateSymbol
	stateSpace

	statePreComment // 表示即将进入注释状态
	stateSingleLineComment  // 单行注释状态
	stateMutliLineComment  // 多行注释状态

	stateNormal
)

// 提取原始token, 并去掉注释
func parse4PrimaryTokens(bs []byte) []Token {
	lexer := newLexer(bs)
	return lexer.run()
}

type Lexer struct {
	bs []byte // 入参,尚未处理的字节流
	preState MachineState // 状态机上一个状态
	state MachineState // 状态机当前状态
	currentByte byte // 遍历处理的当前字节
	currentIndex int // 遍历处理的当前索引
	tmpBytes []byte // 用于暂存长token字符的变量
	ts []Token // 状态机处理后的token列表
	lineIndex int //行索引，用于定位错误。（每个token都会记录自己的行索引）
}

func newLexer(bs []byte) *Lexer {
	return &Lexer{state:stateNormal, bs:bs, lineIndex:1}
}

func (lexer *Lexer) run() []Token {
	for i, b := range lexer.bs {
		lexer.currentByte = b
		lexer.currentIndex = i

		//fmt.Println("in loop:", lexer.CurrentByteString(), lexer.CurrentStateName(), tokensString(lexer.ts))
		if lexer.inStateSingleLineComment() && b != '\n' {
			continue
		}
		if lexer.inStateMutliLineComment() && b != '/' {
			continue
		}
		if lexer.inStateStrLiteral() && b != '"' {
			lexer.tmpBytesCollect()
			continue
		}

		switch {
		case (b>='a' && b<='z') || (b>='A' && b<='Z') || b=='_':
			lexer.whenIdentifier()

		case b >= '0' && b <= '9':
			lexer.whenNumber()

		case b == ' ' || b == '\t' || b == '\n':
			lexer.whenSpace()

		case lexer.isSymbol():
			lexer.whenSymBol()

		case b == '"':
			lexer.whenStringLiterial()

		case b == '/':
			lexer.whenDivOrComment()

		case b == '*':
			lexer.whenMultiComment()
		}
	}
	lexer.pushLongToken()
	lexer.pushBoundryToken()

	return lexer.ts
}

func (lexer *Lexer) whenIdentifier() {
	lexer.tmpBytesCollect()
	lexer.setState(stateIdentifier)
}

func (lexer *Lexer) whenNumber() {
	lexer.tmpBytesCollect()

	if lexer.inStateFloat() || lexer.inStateIdentifier() {
		return
	}

	if lexer.inStateDot() {
		lexer.setState(stateFloat)
	} else {
		lexer.setState(stateInt)
	}
}

func (lexer *Lexer) whenSymBol() {
	if lexer.currentByte == '.' && lexer.inStateInt() {
		lexer.tmpBytesCollect()
		lexer.setState(stateDot)
		return
	}
	if lexer.currentByte == '/' {
		lexer.whenDivOrComment()
		return
	}
	if lexer.currentByte == '*' {
		lexer.whenMultiComment()
		return
	}

	lexer.pushLongToken()
	lexer.pushSymbolToken()
	lexer.setState(stateSymbol)
}

func (lexer *Lexer) whenStringLiterial() {
	// 处理字符串字面值
	if len(lexer.tmpBytes) < 1 {
		if lexer.state == stateStrLiteral {
			// 状态机当前状态是stateStrLiteral，且tmpBytes没有值，说遇到空字符串
			lexer.ts = append(lexer.ts, Token{
				lineIndex: lexer.lineIndex,
				str: "",
				t:   Str,
			})
			lexer.setState(stateNormal)
			return
		}

		lexer.setState(stateStrLiteral)
		return
	}

	lastIndex := len(lexer.tmpBytes)-1
	last := lexer.tmpBytes[lastIndex]
	if last != '\\' {
		// 当前字符为'"', 且前一个字符不就转义字符, 则视为字符串结束
		lexer.pushLongToken()
		lexer.setState(stateNormal)
	} else {
		// clear escape character
		lexer.tmpBytes = lexer.tmpBytes[:lastIndex]
		lexer.tmpBytesCollect()
	}
}

func (lexer *Lexer) whenMultiComment()  {
	if lexer.inStatePreComment() {
		// 使状态机进入多行注释状态
		lexer.setState(stateMutliLineComment)
		// 并清除之前添加的'/'token
		lexer.tailTokenClear()
		return
	}

	lexer.pushLongToken()
	lexer.pushSymbolToken()
	lexer.setState(stateSymbol)
}

func (lexer *Lexer) whenDivOrComment() {
	lexer.pushLongToken()

	switch {
	case !lexer.inStateMutliLineComment() && !lexer.inStatePreComment() && !lexer.inStateSingleLineComment():
		// 使状态机进入预注释状态
		lexer.setState(statePreComment)
		// 捕获除法运算符
		lexer.pushSymbolToken()

	case lexer.inStatePreComment():
		// 使状态机进入单行注释状态
		lexer.setState(stateSingleLineComment)
		// 并之前添加的'/'token
		lexer.tailTokenClear()

	case lexer.inStateMutliLineComment() && lexer.bs[lexer.currentIndex-1] == '*':
		// 终结状态机的多行注释状态
		lexer.setState(stateNormal)
	}
}

func (lexer *Lexer) whenSpace() {
	lexer.pushLongToken()

	if lexer.currentByte == '\n' {
		//fmt.Println("when space:", lexer.CurrentByteString(), lexer.CurrentStateName(), lexer.preState == stateMutliLineComment)

		if lexer.inStateSingleLineComment() {
			// 结束状态机的单行注释状态
			lexer.setState(stateNormal)
			// 添加行结束符
			lexer.pushBoundryToken()
			return
		}
		if lexer.preState == stateMutliLineComment {
			return
		}
		//fmt.Println("before pushBoundryToken when space:", lexer.CurrentByteString(), lexer.CurrentStateName(), lexer.preState == stateMutliLineComment)
		lexer.pushBoundryToken()
		lexer.lineIndex++
	}
	lexer.setState(stateSpace)
}

// 把当前字节追加到临时字节列表中
func (lexer *Lexer) tmpBytesCollect() {
	lexer.tmpBytes = append(lexer.tmpBytes, lexer.currentByte)
}

func (lexer *Lexer) setState(state MachineState)  {
	lexer.preState = lexer.state
	lexer.state = state
}

func (lexer *Lexer) tailTokenClear() {
	size := len(lexer.ts)
	if size < 1 {
		return
	}
	lexer.ts = (lexer.ts)[:size-1]
}


func (lexer *Lexer) pushSymbolToken() {
	symbol := symbolToken(string(lexer.currentByte))
	symbol.lineIndex = lexer.lineIndex

	last, lastExist := lastToken(lexer.ts)
	if symbol.assertSymbol("}") && lastExist && last.assertSymbol(";") {
		// 去掉无用的";"
		tailIndex := len(lexer.ts)-1
		lexer.ts[tailIndex] = symbol
	} else {
		lexer.ts = append(lexer.ts, symbol)
	}
}

func (lexer *Lexer) pushBoundryToken() {
	size := len(lexer.ts)
	if size>0 && lexer.ts[size-1].assertSymbols("{", ",") {
		// 防止添加无用的";",
		// 前一个token为symbol"}", 因为要考虑json对象字面值的情况
		return
	}

	lexer.ts = append(lexer.ts, Token{
		lineIndex: lexer.lineIndex,
		str: ";",
		t:   Symbol,
	})
	//fmt.Println("before append pushBoundryToken:", tokensString(lexer.ts), lexer.CurrentByteString(), lexer.CurrentStateName(), size>0 && lexer.ts[size-1].assertSymbols("{", ",", "}"))
}

func (lexer *Lexer) pushLongToken() {
	if len(lexer.tmpBytes) < 1 {
		return
	}
	s := string(lexer.tmpBytes)

	var tokenType TokenType
	if lexer.inStateFloat() {
		tokenType = Float
	}
	if lexer.inStateInt() {
		if lexer.currentByte == '.' {
			return
		}
		tokenType = Int
	}
	if lexer.inStateIdentifier() {
		tokenType = Identifier
	}
	if lexer.inStateStrLiteral() {
		tokenType = Str
	}

	lexer.ts = append(lexer.ts, Token{
		lineIndex: lexer.lineIndex,
		str: s,
		t:   tokenType,
	})
	// 重置临时变量
	lexer.tmpBytes = nil
}

func (lexer *Lexer) isSymbol() bool {
	switch lexer.currentByte {
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
	case '!': fallthrough
	case '+': fallthrough
	case '-': fallthrough
	case '*': fallthrough
	case '/': fallthrough
	case '%': fallthrough
	case '>': fallthrough
	case '<': fallthrough
	case '|': fallthrough
	case '&':
		return true
	}
	return false
}


func (lexer *Lexer) inStateIdentifier() bool {
	return lexer.state == stateIdentifier
}
func (lexer *Lexer) inStateStrLiteral() bool {
	return lexer.state == stateStrLiteral
}
func (lexer *Lexer) inStateInt() bool {
	return lexer.state == stateInt
}
func (lexer *Lexer) inStateDot() bool {
	return lexer.state == stateDot
}
func (lexer *Lexer) inStateFloat() bool {
	return lexer.state == stateFloat
}
func (lexer *Lexer) inStateSymbol() bool {
	return lexer.state == stateSymbol
}
func (lexer *Lexer) inStateSpace() bool {
	return lexer.state == stateSpace
}
func (lexer *Lexer) inStatePreComment() bool {
	return lexer.state == statePreComment
}
func (lexer *Lexer) inStateSingleLineComment() bool {
	return lexer.state == stateSingleLineComment
}
func (lexer *Lexer) inStateMutliLineComment() bool {
	return lexer.state == stateMutliLineComment
}
func (lexer *Lexer) inStateNormal() bool {
	return lexer.state == stateNormal
}


func (lexer *Lexer) CurrentStateName() string {
	switch {
	case lexer.state == stateIdentifier: return "Identifier"
	case lexer.state == stateStrLiteral: return "StrLiteral"
	case lexer.state == stateInt: return "Int"
	case lexer.state == stateDot: return "Dot"
	case lexer.state == stateFloat: return "Float"
	case lexer.state == stateSymbol: return "Symbol"
	case lexer.state == stateSpace: return "Space"
	case lexer.state == statePreComment: return "PreComment"
	case lexer.state == stateSingleLineComment: return "SingleLineComment"
	case lexer.state == stateMutliLineComment: return "MutliLineComment"
	case lexer.state == stateNormal: return "Normal"
	}
	return "unknow"
}

func (lexer *Lexer) CurrentByteString() string {
	if lexer.currentByte == '\n' {
		return `\n`
	}
	return string(lexer.currentByte)
}