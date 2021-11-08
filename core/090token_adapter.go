package core

import "fmt"

type TokenAdapter struct {
	lineIndex int
	val interface{}
	typName string
}

func (t *TokenAdapter) isSymbol() bool {
	return false
}
func (t *TokenAdapter) assertSymbol(s string) bool {
	return false
}
func (t *TokenAdapter) assertSymbols(ss ...string) bool {
	return false
}

func (t *TokenAdapter) assertKey(k string) bool {
	return false
}
func (t *TokenAdapter) assertKeys(ks ...string) bool {
	return false
}

func (t *TokenAdapter) toExpr() PrimaryExpression {
	return nil
}

func (t *TokenAdapter) typeName() string {
	return t.typName
}

func (t *TokenAdapter) row() string {
	return fmt.Sprint(t.lineIndex)
}
func (t *TokenAdapter) rowIndex() int {
	return t.lineIndex
}
func (t *TokenAdapter) String() string {
	return fmt.Sprint(t.val)
}