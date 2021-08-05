package core

import (
	"regexp"
)


func (fns *InternalFunctionSet) Match(tmpl, source string) []interface{} {
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