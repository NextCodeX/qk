package core

import (
	"bytes"
	"fmt"
)

type Token interface {
	isSymbol() bool
	assertSymbol(s string) bool
	assertSymbols(ss ...string) bool
	assertKey(k string) bool
	assertKeys(ks ...string) bool
	toExpr() PrimaryExpression
	typeName() string
	row() string
	rowIndex() int
	String() string
}

func toName(t Token) string {
	if nameToken, ok := t.(*NameToken); ok {
		return nameToken.name()
	}

	panic(fmt.Errorf("%v: t [%v] is not NameToken", t.row(), t))
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
