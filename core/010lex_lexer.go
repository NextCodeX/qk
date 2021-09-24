package core

type MachineState int

// 状态机解析原始token用
const (
	stateIdentifier MachineState = 1 << iota // 标识符状态
	stateStrLiteral
	stateDynamicStrLiteral
	stateInt
	stateDot  // 小数点状态
	stateFloat
	stateSymbol // 表示各种运算符，分隔符
	stateSpace // 空白符

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
	bs []byte             // 入参,尚未处理的字节流
	preState MachineState // 状态机上一个状态
	state MachineState    // 状态机当前状态
	currentByte byte      // 遍历处理的当前字节
	currentIndex int      // 遍历处理的当前索引
	tmpBytes []byte       // 用于暂存长token字符的变量
	ts []Token        // 状态机处理后的token列表
	lineIndex int         //行索引，用于定位错误。（每个token都会记录自己的行索引）
}

func newLexer(bs []byte) *Lexer {
	return &Lexer{state:stateNormal, bs:bs, lineIndex:1}
}

func (lexer *Lexer) run() []Token {
	for i, b := range lexer.bs {
		lexer.currentByte = b
		lexer.currentIndex = i

		if lexer.inStateSingleLineComment() && b != '\n' {
			continue
		}
		if lexer.inStateMultiLineComment() && b != '/' {
			if b == '\n' { // 遇到'\n'，源码行数加一
				lexer.lineIndex ++
			}
			continue
		}
		if (lexer.inStateStrLiteral() && b != '"') || (lexer.inStateDynamicStrLiteral() && b != '`') {
			if b == '\n' { // 遇到'\n'，源码行数加一
				lexer.lineIndex ++
			}
			lexer.tmpBytesCollect()
			continue
		}

		switch {
		case (b>='a' && b<='z') || (b>='A' && b<='Z') || b=='_':
			lexer.whenIdentifier()
		case b >= '0' && b <= '9':
			lexer.whenNumber()
		case b == ' ' || b == '\r' || b == '\t' || b == '\n':
			lexer.whenSpace()
		case lexer.isSymbol():
			lexer.whenSymbol()
		case b == '"':
			lexer.whenStringLiterial()
		case b == '`':
			lexer.whenDynamicStringLiterial()
		default:
			errorf("line%v: unexpected character %v", lexer.lineIndex, b)
		}
	}
	lexer.pushLongToken()

	return lexer.ts
}

func (lexer *Lexer) whenIdentifier() {
	lexer.tmpBytesCollect()
	lexer.setState(stateIdentifier)
}

func (lexer *Lexer) whenNumber() {
	lexer.tmpBytesCollect()

	if lexer.inStateFloat() || lexer.inStateIdentifier() {
		// 处于浮点数状态或标识符状态退出函数，防止状态被覆盖。
		return
	}

	if lexer.inStateDot() {
		// 进入浮点数状态
		lexer.setState(stateFloat)
	} else {
		lexer.setState(stateInt)
	}
}

func (lexer *Lexer) whenSymbol() {
	if lexer.currentByte == '.' && lexer.inStateInt() {
		// 小数点处理
		lexer.tmpBytesCollect()
		lexer.setState(stateDot)
		return
	}

	if lexer.currentByte == '/' {
		// 除法运算符 或 注释的处理
		lexer.whenDivOrComment()
		return
	}
	if lexer.currentByte == '*' {
		// 乘法运算符 或 注释的处理
		lexer.whenMultiComment()
		return
	}

	// 如没有特例， 则把当前字符转成符号token添加至最终token列表中即可。
	// 并从缓存字符列表中提取多字符token, 修改状态机状态。
	lexer.pushLongToken()
	lexer.pushSymbolToken()
	lexer.setState(stateSymbol)
}

func (lexer *Lexer) whenDynamicStringLiterial() {
	// 处理字符串字面值
	if len(lexer.tmpBytes) < 1 {
		if lexer.inStateDynamicStrLiteral() {
			// 状态机当前状态是stateDynamicStrLiteral，且tmpBytes没有值，说遇到空字符串
			lexer.ts = append(lexer.ts, &TokenImpl{
				lineIndex: lexer.lineIndex,
				str: "",
				t:   DynamicStr,
			})
			lexer.setState(stateNormal)
		} else {
			// 状态机当前状态不是字符串字面值状态， 使用状态机进入字符串字面值状态
			lexer.setState(stateDynamicStrLiteral)
		}
		return
	}

	lastIndex := len(lexer.tmpBytes)-1
	lastChar := lexer.tmpBytes[lastIndex]
	if lastChar != '\\' {
		// 当前字符为'`', 且前一个字符不就转义字符, 则视为字符串结束
		lexer.pushLongToken()
		lexer.setState(stateNormal)
	} else {
		// 清除'`'的转义字符
		lexer.tmpBytes = lexer.tmpBytes[:lastIndex]
		lexer.tmpBytesCollect()
	}
}

func (lexer *Lexer) whenStringLiterial() {
	// 处理字符串字面值
	if len(lexer.tmpBytes) < 1 {
		if lexer.inStateStrLiteral() {
			// 状态机当前状态是stateStrLiteral，且tmpBytes没有值，说遇到空字符串
			lexer.ts = append(lexer.ts, &TokenImpl{
				lineIndex: lexer.lineIndex,
				str: "",
				t:   Str,
			})
			lexer.setState(stateNormal)
			return
		}

		// 状态机当前状态不是字符串字面值状态， 使用状态机进入字符串字面值状态
		lexer.setState(stateStrLiteral)
		return
	}

	lastIndex := len(lexer.tmpBytes)-1
	lastChar := lexer.tmpBytes[lastIndex]
	if lastChar != '\\' {
		// 当前字符为'"', 且前一个字符不是转义字符, 则视为字符串结束
		lexer.pushLongToken()
		lexer.setState(stateNormal)
	} else {
		// 清除'"'的转义字符
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

	// 如果状态机不是预注释状态, 把 '*' 当乘法运算符处理。
	lexer.pushLongToken()
	lexer.pushSymbolToken()
	lexer.setState(stateSymbol)
}

func (lexer *Lexer) whenDivOrComment() {
	lexer.pushLongToken()

	switch {
	case !lexer.inStatePreComment() && !lexer.inStateMultiLineComment() && !lexer.inStateSingleLineComment():
		// 使状态机进入预注释状态
		lexer.setState(statePreComment)
		// 捕获除法运算符
		lexer.pushSymbolToken()

	case lexer.inStatePreComment():
		// 使状态机进入单行注释状态
		lexer.setState(stateSingleLineComment)
		// 并之前添加的'/'token
		lexer.tailTokenClear()

	case lexer.inStateMultiLineComment() && lexer.bs[lexer.currentIndex-1] == '*':
		// 终结状态机的多行注释状态
		lexer.setState(stateNormal)

	default:
		errorf("line%v: unknown case('/')", lexer.lineIndex)
	}
}

func (lexer *Lexer) whenSpace() {
	lexer.pushLongToken()

	if lexer.currentByte == '\n' {
		// 遇到'\n'，源码行数加一
		lexer.lineIndex++

		if lexer.inStateSingleLineComment() {
			// 结束状态机的单行注释状态
			lexer.setState(stateNormal)
			// 添加行结束符
			lexer.pushBoundryToken()
			return
		}
		if lexer.preState == stateMutliLineComment {
			// 如状态机前一个状态是多行注释状态直接返回，不用添加语句分隔符';'token
			return
		}
		lexer.pushBoundryToken()
	}
	lexer.setState(stateSpace)
}

// 把当前字节追加到临时字节列表中
func (lexer *Lexer) tmpBytesCollect() {
	lexer.tmpBytes = append(lexer.tmpBytes, lexer.currentByte)
}

// 设置状态机的状态
func (lexer *Lexer) setState(state MachineState) {
	lexer.preState = lexer.state
	lexer.state = state
}

// 删除最终token列表的最后一个token
func (lexer *Lexer) tailTokenClear() {
	size := len(lexer.ts)
	if size < 1 {
		return
	}
	lexer.ts = (lexer.ts)[:size-1]
}

// 添加符号token至最终token列表
func (lexer *Lexer) pushSymbolToken() {
	last, lastExist := lastToken(lexer.ts)
	if lastExist && lexer.operatorMerge(last) {
		// 存在运算符合并时，直接退出函数
		return
	}

	symbol := symbolToken(string(lexer.currentByte))
	symbol.setLineIndex(lexer.lineIndex)
	if lastExist && last.assertSymbol(";") &&
		(lexer.currentByte == '}' || lexer.currentByte == ']') {
		// 为JSONObject, JSONArray字面值去掉无用的";"
		tailIndex := len(lexer.ts)-1
		lexer.ts[tailIndex] = symbol
	} else {
		// 正常捕获符号token
		lexer.ts = append(lexer.ts, symbol)
	}
}

// 运算符合并
func (lexer *Lexer) operatorMerge(last Token) bool {
	flag := lexer.currentByte == '>' && last.assertSymbol("-") ||
		lexer.currentByte == '=' && last.assertSymbols("=", ">", "<", "+", "-", "*", "/", "%", "!") ||
		lexer.currentByte == '|' && last.assertSymbols("|") ||
		lexer.currentByte == '&' && last.assertSymbols("&") ||
		lexer.currentByte == '+' && last.assertSymbol("+") ||
		lexer.currentByte == '-' && last.assertSymbol("-")

	if flag {
		last.setRaw(last.raw() + string(lexer.currentByte))
	}

	return flag
}

// 添加语句分隔符';'token, 并避免多余的';'
func (lexer *Lexer) pushBoundryToken() {
	size := len(lexer.ts)
	if size>0 && lexer.ts[size-1].assertSymbols("{", ",", "?", ":", "[", ";") {
		// 防止添加无用的";"
		return
	}

	lexer.ts = append(lexer.ts, &TokenImpl{
		lineIndex: lexer.lineIndex,
		str: ";",
		t:   Symbol,
	})
}

// 提取多字符token, 并添加至最终token列表中
func (lexer *Lexer) pushLongToken() {
	if len(lexer.tmpBytes) < 1 {
		return
	}

	var tokenType TokenType
	if lexer.inStateFloat() {
		tokenType = Float
		if lexer.negativeNumberHandler(Float) {
			// 如果成功捕获负数，退出函数
			return
		}

	} else if lexer.inStateInt() {
		tokenType = Int
		if lexer.negativeNumberHandler(Int) {
			// 如果成功捕获负数，退出函数
			return
		}

	} else if lexer.inStateIdentifier() {
		tokenType = Identifier

	} else if lexer.inStateStrLiteral() {
		tokenType = Str

	} else if lexer.inStateDynamicStrLiteral() {
		tokenType = DynamicStr
	} else {}


	lexer.ts = append(lexer.ts, &TokenImpl{
		lineIndex: lexer.lineIndex,
		str: string(lexer.tmpBytes),
		t:   tokenType,
	})
	// 重置临时变量
	lexer.tmpBytes = nil
}

// 捕获负数
func (lexer *Lexer) negativeNumberHandler(tokenType TokenType) bool {
	last, lastExist := lastToken(lexer.ts)
	if lastExist && last.assertSymbol("-") {
		lastSecond, lastSecondExist := lastSecondToken(lexer.ts)
		if !lastSecondExist || lastSecondExist && (lastSecond.assertIdentifier("return") ||
			lastSecond.assertSymbols("+", "-", "*", "/", "=", ",", "(", ":", "[", "->", "{")) {

			last.setRaw(last.raw() + string(lexer.tmpBytes))
			last.setTyp(tokenType)
			// 重置临时变量
			lexer.tmpBytes = nil
			return true
		}
	}
	return false
}

func (lexer *Lexer) isSymbol() bool {
	switch lexer.currentByte {
	case '.', '?', ':', '(', ')', '[', ']', '{', '}', ';', ',', '=', '!', '+', '-', '*', '/', '%', '>', '<', '|', '&', '$':
		return true
	default:
		return false
	}
}


func (lexer *Lexer) inStateIdentifier() bool {
	return (lexer.state & stateIdentifier) == stateIdentifier
}
func (lexer *Lexer) inStateStrLiteral() bool {
	return (lexer.state & stateStrLiteral) == stateStrLiteral
}
func (lexer *Lexer) inStateDynamicStrLiteral() bool {
	return (lexer.state & stateDynamicStrLiteral) == stateDynamicStrLiteral
}
func (lexer *Lexer) inStateInt() bool {
	return (lexer.state & stateInt) == stateInt
}
func (lexer *Lexer) inStateDot() bool {
	return (lexer.state & stateDot) == stateDot
}
func (lexer *Lexer) inStateFloat() bool {
	return (lexer.state & stateFloat) == stateFloat
}
func (lexer *Lexer) inStateSymbol() bool {
	return (lexer.state & stateSymbol) == stateSymbol
}
func (lexer *Lexer) inStateSpace() bool {
	return (lexer.state & stateSpace) == stateSpace
}
func (lexer *Lexer) inStatePreComment() bool {
	return (lexer.state & statePreComment) == statePreComment
}
func (lexer *Lexer) inStateSingleLineComment() bool {
	return (lexer.state & stateSingleLineComment) == stateSingleLineComment
}
func (lexer *Lexer) inStateMultiLineComment() bool {
	return (lexer.state & stateMutliLineComment) == stateMutliLineComment
}
func (lexer *Lexer) inStateNormal() bool {
	return (lexer.state & stateNormal) == stateNormal
}


func (lexer *Lexer) CurrentStateName() string {
	var stateName string
	if lexer.inStateDynamicStrLiteral() {
		stateName += "dynamicStr, "
	}
	if lexer.inStateDot() {
		stateName += "dot, "
	}
	if lexer.inStateFloat() {
		stateName += "float, "
	}
	if lexer.inStateIdentifier() {
		stateName += "identifier, "
	}
	if lexer.inStateInt() {
		stateName += "int, "
	}
	if lexer.inStateMultiLineComment() {
		stateName += "multiLineComment, "
	}
	if lexer.inStateNormal() {
		stateName += "normal, "
	}
	if lexer.inStatePreComment() {
		stateName += "preComment, "
	}
	if lexer.inStateSingleLineComment() {
		stateName += "singleLineComment, "
	}
	if lexer.inStateSpace() {
		stateName += "space, "
	}
	if lexer.inStateDynamicStrLiteral() {
		stateName += "dynamicStrLiteral, "
	}
	if lexer.inStateStrLiteral() {
		stateName += "strLiteral, "
	}
	if lexer.inStateSymbol() {
		stateName += "symbol, "
	}

	return stateName
}

func (lexer *Lexer) CurrentByteString() string {
	if lexer.currentByte == '\n' {
		return `\n`
	}
	return string(lexer.currentByte)
}