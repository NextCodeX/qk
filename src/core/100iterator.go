package core

type Iterator interface {
	indexs() []interface{}
	getItem(interface{}) *Value
}

func toIterator(v *Value) Iterator {
	if v.isArrayValue() {
		return v.jsonArr
	}
	if v.isObjectValue() {
		return v.jsonObj
	}
	return nil
}

type ForPlusInfo struct {
	indexName string
	itemName string
	iterator string
}

func newForPlusInfo(indexName, itemName, iterator string) *ForPlusInfo {
	return &ForPlusInfo{indexName, itemName, iterator }
}
