package core

import (
    "bytes"
    "strings"
)

type StringValue struct {
    chars []rune
    goValue string
    ClassObject
}

func newStringValue(raw string) Value {
    var chs []rune
    for _, ch := range raw {
        chs = append(chs, ch)
    }
    str := &StringValue{goValue: raw, chars: chs}
    str.initAsClass("String", &str)
    return str
}

func (str *StringValue) val() interface{} {
    return str.goValue
}
func (str *StringValue) isString() bool {
    return true
}

func (str *StringValue) String() string {
    return str.goValue
}



func (str *StringValue) indexs() []interface{} {
    var indexs []interface{}

    for i := 0; i < len(str.chars); i++ {
        indexs = append(indexs, i)
    }
    return indexs
}

func (str *StringValue) getItem(index interface{}) Value {
    i := index.(int)
    return newQKValue(string(str.chars[i]))
}

func (str *StringValue) getChar(index int) string {
    return string(str.chars[index])
}

func (str *StringValue) sub(start, end int) string {
    var buf bytes.Buffer
    for _, ch := range str.chars[start:end] {
        buf.WriteRune(ch)
    }
    return buf.String()
}

func (str *StringValue) Bytes() []byte {
    return []byte(str.goValue)
}

func (str *StringValue) Size() int {
    return len(str.chars)
}

func (str *StringValue) Index(subStr string) int {
    return strings.Index(str.goValue, subStr)
}

func (str *StringValue) LastIndex(subStr string) int {
    return strings.LastIndex(str.goValue, subStr)
}
func (str *StringValue) Trim() string {
    return strings.TrimSpace(str.goValue)
}
func (str *StringValue) Replace(old, newStr string) string {
    return strings.ReplaceAll(str.goValue, old, newStr)
}
func (str *StringValue) Contain(subStr string) bool {
    return strings.Contains(str.goValue, subStr)
}
func (str *StringValue) Lower() string {
    return strings.ToLower(str.goValue)
}
func (str *StringValue) Upper() string {
    return strings.ToUpper(str.goValue)
}
func (str *StringValue) LowerFirst() string {
    return strings.ToLower(string(str.chars[0])) + string(str.chars[1:])
}
func (str *StringValue) UpperFirst() string {
    return strings.ToUpper(string(str.chars[0])) + string(str.chars[1:])
}
func (str *StringValue) ToTitle() string {
    return strings.ToTitle(str.goValue)
}
func (str *StringValue) Title() string {
    return strings.Title(str.goValue)
}
func (str *StringValue) HasPrefix(prefix string) bool {
    return strings.HasPrefix(str.goValue, prefix)
}
func (str *StringValue) HasSuffix(suffix string) bool {
    return strings.HasSuffix(str.goValue, suffix)
}
func (str *StringValue) Split(seperator string) []string {
    return strings.Split(str.goValue, seperator)
}
func (str *StringValue) Eic(target string) bool {
    return str.goValue == target || strings.ToLower(str.goValue) == strings.ToLower(target)
}
