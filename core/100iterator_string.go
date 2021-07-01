package core

type StringIterator struct {
	chars []*Value
	indexArr []interface{}
}

func newStringIterator(raw string) *StringIterator {
	var i int
	var indexs []interface{}
	var ss []*Value
	for _, item := range raw {
		indexs = append(indexs, i)
		char := newQkValue(string(item))
		ss = append(ss, char)
		i++
	}
	return &StringIterator{ss, indexs}
}

func (strIterator *StringIterator) indexs() []interface{} {
	return strIterator.indexArr
}

func (strIterator *StringIterator) getItem(index interface{}) *Value {
	i := index.(int)
	return strIterator.chars[i]
}