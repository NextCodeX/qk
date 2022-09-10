package core

import (
	"bytes"
	"strconv"
	"strings"
	"unicode/utf8"
)

type StringValue struct {
	chars   []rune
	goValue string
	ClassObject
}

func newStringValue(raw string) Value {
	chs := make([]rune, utf8.RuneCountInString(raw))
	var i uint32 = 0
	for _, ch := range raw {
		chs[i] = ch
		i++
	}
	str := &StringValue{goValue: raw, chars: chs}
	str.initAsClass("String", &str)
	return str
}

func (this *StringValue) val() interface{} {
	return this.goValue
}
func (this *StringValue) isString() bool {
	return true
}

func (this *StringValue) String() string {
	return this.goValue
}

func (this *StringValue) indexs() []interface{} {
	size := len(this.chars)
	indexs := make([]interface{}, size)

	for i := 0; i < size; i++ {
		indexs[i] = i
	}
	return indexs
}

func (this *StringValue) getItem(index interface{}) Value {
	i := index.(int)
	return newQKValue(string(this.chars[i]))
}

func (this *StringValue) getChar(index int) string {
	return string(this.chars[index])
}
func (this *StringValue) At(index int) string {
	return this.getChar(index)
}

func (this *StringValue) sub(start, end int) string {
	var buf bytes.Buffer
	for _, ch := range this.chars[start:end] {
		buf.WriteRune(ch)
	}
	return buf.String()
}

func (this *StringValue) Int() int {
	res, err := strconv.Atoi(this.goValue)
	if err != nil {
		return -1
	}
	return res
}
func (this *StringValue) Float() float64 {
	res, err := strconv.ParseFloat(this.goValue, 64)
	if err != nil {
		return -1
	}
	return res
}
func (this *StringValue) Number() interface{} {
	return strToNumber(this.goValue)
}
func (this *StringValue) Bool() bool {
	res, err := strconv.ParseBool(this.goValue)
	if err != nil {
		runtimeExcption(err)
	}
	return res
}

func (this *StringValue) Bytes() []byte {
	return []byte(this.goValue)
}

func (this *StringValue) Size() int {
	return len(this.chars)
}

func (this *StringValue) Index(subStr string) int {
	return strings.Index(this.goValue, subStr)
}
func (this *StringValue) LastIndex(subStr string) int {
	return strings.LastIndex(this.goValue, subStr)
}
func (this *StringValue) Trim() string {
	return strings.TrimSpace(this.goValue)
}
func (this *StringValue) Replace(old, newStr string) string {
	return strings.ReplaceAll(this.goValue, old, newStr)
}
func (this *StringValue) Repl(old, newStr string) string {
	return strings.ReplaceAll(this.goValue, old, newStr)
}
func (this *StringValue) Clear(target string) string {
	return strings.ReplaceAll(this.goValue, target, "")
}
func (this *StringValue) Contains(subStr string) bool {
	return strings.Contains(this.goValue, subStr)
}
func (this *StringValue) Has(subStr string) bool {
	return strings.Contains(this.goValue, subStr)
}
func (this *StringValue) Lower() string {
	return strings.ToLower(this.goValue)
}
func (this *StringValue) Upper() string {
	return strings.ToUpper(this.goValue)
}
func (this *StringValue) LowerFirst() string {
	return strings.ToLower(string(this.chars[0])) + string(this.chars[1:])
}
func (this *StringValue) UpperFirst() string {
	return strings.ToUpper(string(this.chars[0])) + string(this.chars[1:])
}
func (this *StringValue) ToTitle() string {
	return strings.ToTitle(this.goValue)
}
func (this *StringValue) Title() string {
	return strings.Title(this.goValue)
}
func (this *StringValue) HasPrefix(prefix string) bool {
	return strings.HasPrefix(this.goValue, prefix)
}
func (this *StringValue) HasSuffix(suffix string) bool {
	return strings.HasSuffix(this.goValue, suffix)
}
func (this *StringValue) Split(seperator string, rawFlag bool) []string {
	ss := strings.Split(this.goValue, seperator)
	if rawFlag {
		return ss
	}
	for i, s := range ss {
		ss[i] = strings.TrimSpace(s)
	}
	return ss
}
func (this *StringValue) Is(target string) bool {
	return this.goValue == target || strings.ToLower(this.goValue) == strings.ToLower(target)
}

func (this *StringValue) Match(tmpl string) bool {
	return regMatch(tmpl, this.goValue)
}
func (this *StringValue) Find(tmpl string) []interface{} {
	return regFind(tmpl, this.goValue)
}
