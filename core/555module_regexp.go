package core

import (
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
func (reg *QKRegexp) Match(src string) bool {
	return reg.raw.MatchString(src)
}

// 查找所有匹配模式符合的字符串
func (reg *QKRegexp) Find(src string) []string {
	return reg.raw.FindAllString(src, -1)
}

// 查找匹配模式的字符串，返回左侧第一个匹配的结果。
func (reg *QKRegexp) FindOne(src string) string {
	return reg.raw.FindString(src)
}

// 查找所有匹配字符串的起止位置
func (reg *QKRegexp) FindIndex(src string) Value {
	var res []Value
	tmp := reg.raw.FindAllStringIndex(src, -1)
	for _, item := range tmp {
		var pos []Value
		for _, subItem := range item {
			pos = append(pos, newQKValue(subItem))
		}
		res = append(res, array(pos))
	}
	return array(res)
}

// 查找第一个匹配字符串的起止位置
func (reg *QKRegexp) FindOneIndex(src string) Value {
	var res []Value
	tmp := reg.raw.FindStringIndex(src)
	for _, item := range tmp {
		res = append(res, newQKValue(item))
	}
	return array(res)
}

// 查找子匹配符合的所有字符串
func (reg *QKRegexp) FindSub(src string) Value {
	subs := reg.raw.FindAllStringSubmatch(src, -1)
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

// 查找子匹配符合的第一个字符串
func (reg *QKRegexp) FindOneSub(src string) Value {
	subs := reg.raw.FindStringSubmatch(src)
	var res []Value
	for _, sub := range subs {
		res = append(res, newQKValue(sub))
	}
	return array(res)
}

// 通过原始字符串替换
func (reg *QKRegexp) ReplaceByStr(src, repl string) string {
	return reg.raw.ReplaceAllLiteralString(src, repl)
}
func (reg *QKRegexp) Rs(src, repl string) string {
	return reg.ReplaceByStr(src, repl)
}

// 通过正则字符串替换
func (reg *QKRegexp) ReplaceByReg(src, repl string) string {
	return reg.raw.ReplaceAllString(src, repl)
}
func (reg *QKRegexp) Rr(src, repl string) string {
	return reg.ReplaceByReg(src, repl)
}

// 通过函数替换
func (reg *QKRegexp) ReplaceByFunc(src string, replFunc Function) string {
	return reg.raw.ReplaceAllStringFunc(src, func(old string) string {
		args := make([]Value, 0, 1)
		args = append(args, newQKValue(old))
		replFunc.setArgs(args)
		execRes := replFunc.execute()
		return goStr(execRes.value())
	})
}
func (reg *QKRegexp) Rf(src string, replFunc Function) string {
	return reg.ReplaceByFunc(src, replFunc)
}

// 分割
func (reg *QKRegexp) Split(src string) []string {
	return reg.raw.Split(src, -1)
}
