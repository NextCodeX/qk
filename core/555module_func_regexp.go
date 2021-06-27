package core

import (
	"regexp"
)

func (mr *ModuleRegister) RegexpModuleInit() {
	re := &QkRegexp{}
	fre := collectFunctionInfo(&re)
	functionRegister("reg", fre)
}

type QkRegexp struct {}

func (re *QkRegexp) Match(tmpl, source string) []interface{} {
	matcher := regexp.MustCompile(tmpl)
	list := matcher.FindAllStringSubmatch(source, -1)

	var res []interface{}
	if len(list) < 1 {
		return res
	}
	for _, item := range list {
		var subList []interface{}
		for _, subItem := range item {
			subList = append(subList, subItem)
		}
		res = append(res, subList)
	}
	return res
}