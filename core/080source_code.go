package core

type SourceCode interface {
	tokenList() []Token
	setTokenList(ts []Token)
	source() string
}
type SourceCodeImpl struct {
	tokens []Token
}

func (src *SourceCodeImpl) tokenList() []Token {
	return src.tokens
}
func (src *SourceCodeImpl) setTokenList(ts []Token) {
	src.tokens = ts
}
func (src *SourceCodeImpl) source() string {
	return tokensString(src.tokens)
}
