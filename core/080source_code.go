package core

type SourceCode interface {
	tokenList() []Token
	setTokenList(ts []Token)
	source() string
}
type SourceCodeImpl struct {
	tokens   []Token
	errPrint bool
}

func (this *SourceCodeImpl) tokenList() []Token {
	return this.tokens
}
func (this *SourceCodeImpl) setTokenList(ts []Token) {
	this.tokens = ts
}
func (this *SourceCodeImpl) source() string {
	return tokensString(this.tokens)
}
