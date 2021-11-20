package core

import (
	"bytes"
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
		panic(err)
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
		panic(err)
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
	if raw {
		t.val = TRUE
	} else {
		t.val = FALSE
	}
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
	if !dynamic {
		raw = t.strEscape(raw)
	}
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
func (t *StrToken) strEscape(raw string) string {
	var buf bytes.Buffer

	chars := []rune(raw)
	size := len(chars)
	for i := 0; i < size; i++ {
		currentChar := chars[i]
		if currentChar != '\\' {
			buf.WriteRune(currentChar)
			continue
		}

		nextChar := chars[i+1]
		switch nextChar {
		case 'a':
			buf.WriteByte('\a') //\a 响铃(BEL)
		case 'b':
			buf.WriteByte('\b') //\b 退格(BS)
		case 'f':
			buf.WriteByte('\f') //\f 换页(FF)
		case 'n':
			buf.WriteByte('\n') //\n 换行(LF)
		case 'r':
			buf.WriteByte('\r') //\r 回车(CR)
		case 't':
			buf.WriteByte('\t') //\t 水平制表(HT)
		case 'v':
			buf.WriteByte('\v') //\v 垂直制表(VT)
		case '\\':
			buf.WriteByte('\\') //\\ 反斜杠
		case '"':
			buf.WriteByte('"')
		default:
			errorf("line%v: invalid string literal %v%v", t.lineIndex, string('\\'), string(nextChar))
		}
		i++
	}
	return buf.String()
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
