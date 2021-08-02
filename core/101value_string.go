package core

import "bytes"

type StringValue struct {
    chars []rune
    goValue string
    ValueAdapter
}

func newStringValue(raw string) Value {
    var chs []rune
    for _, ch := range raw {
        chs = append(chs, ch)
    }
    return &StringValue{goValue: raw, chars: chs}
}

func (str *StringValue) val() interface{} {
    return str.goValue
}

func (str *StringValue) size() int {
    return len(str.chars)
}

func (str *StringValue) sub(start, end int) string {
    var buf bytes.Buffer
    for _, ch := range str.chars[start:end] {
        buf.WriteRune(ch)
    }
    return buf.String()
}

func (str *StringValue) isString() bool {
    return true
}
