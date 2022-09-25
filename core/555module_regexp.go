package core

import (
	"fmt"
	"regexp"
)

// 正则匹配检测
func regMatch(tmpl, source string) bool {
	matcher := regexp.MustCompile(tmpl)

	return matcher.MatchString(source)
}

// 正则查找
func regFind(tmpl, source string) []interface{} {
	matcher := regexp.MustCompile(tmpl)
	list := matcher.FindAllStringSubmatch(source, -1)

	var res []interface{}
	if len(list) < 1 {
		return res
	}
	for _, item := range list {
		if len(item) == 1 {
			res = append(res, item[0])
			continue
		}
		var subList []interface{}
		for _, subItem := range item {
			subList = append(subList, subItem)
		}
		res = append(res, subList)
	}
	return res
}

// 新建一个正则对象
func (this *InternalFunctionSet) Regexp(pattern string) Value {
	return newRegexp(pattern)
}
func newRegexp(pattern string) Value {
	matcher := regexp.MustCompile(pattern)
	obj := &QKRegexp{matcher}
	return newClass("QKRegexp", &obj)
}

type QKRegexp struct {
	raw *regexp.Regexp
}

// 测试字符串是否匹配正则
func (this *QKRegexp) Match(src string) bool {
	return this.raw.MatchString(src)
}

// 查找所有匹配模式符合的字符串
func (this *QKRegexp) Find(src string) []string {
	return this.raw.FindAllString(src, -1)
}
func (this *QKRegexp) F(src string) []string {
	return this.Find(src)
}

// 查找匹配模式的字符串，返回左侧第一个匹配的结果。
func (this *QKRegexp) FindOne(src string) string {
	return this.raw.FindString(src)
}
func (this *QKRegexp) Ff(src string) string {
	return this.FindOne(src)
}

// 查找所有匹配字符串的起止位置
func (this *QKRegexp) FindIndex(src string) Value {
	var res []Value
	tmp := this.raw.FindAllStringIndex(src, -1)
	for _, item := range tmp {
		var pos []Value
		for _, subItem := range item {
			pos = append(pos, newQKValue(subItem))
		}
		res = append(res, array(pos))
	}
	return array(res)
}
func (this *QKRegexp) Fi(src string) Value {
	return this.FindIndex(src)
}

// 查找第一个匹配字符串的起止位置
func (this *QKRegexp) FindOneIndex(src string) Value {
	var res []Value
	tmp := this.raw.FindStringIndex(src)
	for _, item := range tmp {
		res = append(res, newQKValue(item))
	}
	return array(res)
}
func (this *QKRegexp) Ffi(src string) Value {
	return this.FindOneIndex(src)
}

// 查找子匹配符合的所有字符串
func (this *QKRegexp) FindSub(src string) Value {
	subs := this.raw.FindAllStringSubmatch(src, -1)
	var res []Value
	for _, sub := range subs {
		var items []Value
		for _, item := range sub {
			items = append(items, newQKValue(item))
		}
		res = append(res, array(items))
	}
	return array(res)
}
func (this *QKRegexp) Fs(src string) Value {
	return this.FindSub(src)
}

// 查找子匹配符合的第一个字符串
func (this *QKRegexp) FindOneSub(src string) Value {
	subs := this.raw.FindStringSubmatch(src)
	var res []Value
	for _, sub := range subs {
		res = append(res, newQKValue(sub))
	}
	return array(res)
}
func (this *QKRegexp) Ffs(src string) Value {
	return this.FindOneSub(src)
}

// 通过原始字符串替换
func (this *QKRegexp) ReplaceByStr(src, repl string) string {
	return this.raw.ReplaceAllLiteralString(src, repl)
}
func (this *QKRegexp) Rs(src, repl string) string {
	return this.ReplaceByStr(src, repl)
}

// 通过正则字符串替换
func (this *QKRegexp) ReplaceByReg(src, repl string) string {
	return this.raw.ReplaceAllString(src, repl)
}
func (this *QKRegexp) Rr(src, repl string) string {
	return this.ReplaceByReg(src, repl)
}

// 通过函数替换
func (this *QKRegexp) ReplaceByFunc(src string, replFunc Function) string {
	return this.raw.ReplaceAllStringFunc(src, func(old string) string {
		args := make([]Value, 0, 1)
		args = append(args, newQKValue(old))
		replFunc.setArgs(args)
		execRes := replFunc.execute()
		return fmt.Sprint(execRes.value())
	})
}
func (this *QKRegexp) Rf(src string, replFunc Function) string {
	return this.ReplaceByFunc(src, replFunc)
}

// 分割
func (this *QKRegexp) Split(src string) []string {
	return this.raw.Split(src, -1)
}
