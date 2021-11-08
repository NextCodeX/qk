package core

import (
	"fmt"
	"strconv"
)

type PrimaryToken struct {
	TokenAdapter
}
func (t *PrimaryToken) toExpr() PrimaryExpression {
	return newConstPrimaryExpression(newQKValue(t.val))
}

type IntToken struct {
	PrimaryToken
}
func newIntToken(raw string, row int) Token {
	i, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		panic(i)
	}

	t := &IntToken{}
	t.lineIndex = row
	t.val = i
	t.typName = "Integer"
	return t
}


type FloatToken struct {
	PrimaryToken
}
func newFloatToken(raw string, row int) Token {
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		panic(f)
	}

	t := &FloatToken{}
	t.lineIndex = row
	t.val = f
	t.typName = "Float"
	return t
}


type BoolToken struct {
	PrimaryToken
}
func newBoolToken(raw bool, row int) Token {
	t := &BoolToken{}
	t.lineIndex = row
	t.val = raw
	t.typName = "Boolean"
	return t
}


type StrToken struct {
	dynamic bool
	TokenAdapter
}
func newStrToken(raw string, row int, dynamic bool) Token {
	t := &StrToken{}
	t.dynamic = dynamic
	t.lineIndex = row
	t.val = raw
	t.typName = "String"
	return t
}
func (t *StrToken) toExpr() PrimaryExpression {
	if t.dynamic {
		return newDynamicStrPrimaryExpression(t.val.(string))
	}
	return newConstPrimaryExpression(newQKValue(t.val))
}
func (t *StrToken) String() string {
	if t.dynamic {
		return fmt.Sprintf("`%v`", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

type NullToken struct {
	PrimaryToken
}
func newNullToken(row int) Token {
	t := &NullToken{}
	t.lineIndex = row
	t.val = NULL
	t.typName = "NULL"
	return t
}
