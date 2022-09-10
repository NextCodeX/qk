package core

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	CRLF = "\r\n"
	LF   = "\n"
	CR   = "\r"
)

type StringValue struct {
	goValue string
	ClassObject
}

func newStringValue(raw string) Value {
	str := &StringValue{goValue: raw}
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
	size := utf8.RuneCountInString(this.goValue)
	indexs := make([]interface{}, size)
	for i := 0; i < size; i++ {
		indexs[i] = i
	}
	return indexs
}

func (this *StringValue) getItem(index interface{}) Value {
	i := index.(int)
	return newQKValue(getCharAt(this.goValue, i))
}
func getCharAt(raw string, index int) string {
	for ii, ch := range raw {
		if ii == index {
			return string(ch)
		}
	}
	return ""
}
func (this *StringValue) getChar(index int) string {
	return getCharAt(this.goValue, index)
}
func (this *StringValue) At(index int) string {
	return this.getChar(index)
}

func (this *StringValue) sub(start, end int) string {
	var buf bytes.Buffer
	for _, ch := range strToChars(this.goValue)[start:end] {
		buf.WriteRune(ch)
	}
	return buf.String()
}
func (this *StringValue) Left(seperator string) string {
	index := strings.Index(this.goValue, seperator)
	if index < 0 {
		return this.goValue
	}
	return this.goValue[:index]
}
func (this *StringValue) Right(seperator string) string {
	index := strings.Index(this.goValue, seperator)
	if index < 0 {
		return this.goValue
	}
	return this.goValue[index:]
}

func strToChars(raw string) []rune {
	chs := make([]rune, utf8.RuneCountInString(raw))
	var i uint32 = 0
	for _, ch := range raw {
		chs[i] = ch
		i++
	}
	return chs
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

func (this *StringValue) Exec() string {
	return doCmd(this.goValue)
}
func (this *StringValue) Bytes() []byte {
	return []byte(this.goValue)
}
func (this *StringValue) Save(path string) {
	fileSave(path, []byte(this.goValue))
}
func (this *StringValue) Base64() string {
	return base64.StdEncoding.EncodeToString([]byte(this.goValue))
}
func (this *StringValue) Debase64() []byte {
	data, err := base64.StdEncoding.DecodeString(this.goValue)
	if err != nil {
		return nil
	}
	return data
}
func (this *StringValue) Gzip() []byte {
	return gzipEncode([]byte(this.goValue))
}
func (this *StringValue) DeGzip() []byte {
	return gzipDecode([]byte(this.goValue))
}

func (this *StringValue) Size() int {
	return utf8.RuneCountInString(this.goValue)
}

func (this *StringValue) ToJson() any {
	rawStr := strings.TrimSpace(this.goValue)
	if strings.HasPrefix(rawStr, "{") && strings.HasSuffix(rawStr, "}") {
		// Declared an empty map interface
		var result map[string]interface{}

		// Unmarshal or Decode the JSON to the interface.
		err := json.Unmarshal([]byte(rawStr), &result)
		if err != nil {
			return nil
		}
		return result
	}
	if strings.HasPrefix(rawStr, "[") && strings.HasSuffix(rawStr, "]") {
		var arr []any
		err := json.Unmarshal([]byte(rawStr), &arr)
		if err != nil {
			return nil
		}
		return arr
	}
	return nil
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
func (this *StringValue) ToLine() string {
	tmp := strings.ReplaceAll(this.goValue, CRLF, "")
	tmp = strings.ReplaceAll(tmp, CR, "")
	return strings.ReplaceAll(tmp, LF, "")
}
func (this *StringValue) Lines() []string {
	if strings.Contains(this.goValue, CRLF) {
		return trimSpaces(strings.Split(this.goValue, CRLF))
	}
	if strings.Contains(this.goValue, CR) {
		return trimSpaces(strings.Split(this.goValue, CR))
	}
	return trimSpaces(strings.Split(this.goValue, LF))
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
	if len(this.goValue) < 1 {
		return ""
	}
	return strings.ToLower(this.goValue[:1]) + this.goValue[1:]
}
func (this *StringValue) UpperFirst() string {
	if len(this.goValue) < 1 {
		return ""
	}
	return strings.ToUpper(this.goValue[:1]) + this.goValue[1:]
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
	return trimSpaces(ss)
}
func trimSpaces(ss []string) []string {
	for i, s := range ss {
		ss[i] = strings.TrimSpace(s)
	}
	return ss
}
func (this *StringValue) Is(target string) bool {
	return this.goValue == target || strings.ToLower(this.goValue) == strings.ToLower(target)
}
func (this *StringValue) In(targets []string, ignoreCase bool) bool {
	for _, target := range targets {
		if !ignoreCase && this.goValue == target {
			return true
		}
		if ignoreCase && strings.ToLower(this.goValue) == strings.ToLower(target) {
			return true
		}
	}
	return false
}

func (this *StringValue) Match(tmpl string) bool {
	return regMatch(tmpl, this.goValue)
}
func (this *StringValue) Find(tmpl string) []interface{} {
	return regFind(tmpl, this.goValue)
}
