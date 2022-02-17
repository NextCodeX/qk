package core

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JSONArray interface {
	Size() int
	add(elem Value)
	set(index int, elem Value)
	getElem(index int) Value
	sub(start, end int) Value
	checkOutofIndex(index int) bool
	values() []Value
	String() string
	toJSONArrayString() string
	Iterator
	Value
}

type JSONArrayImpl struct {
	valList []Value
	ClassObject
}

func array(v []Value) JSONArray {
	arr := &JSONArrayImpl{valList: v}
	arr.ClassObject.initAsClass("JSONArray", &arr)
	return arr
}

func emptyArray() JSONArray {
	vals := make([]Value, 0, 16)
	return array(vals)
}

func (arr *JSONArrayImpl) Size() int {
	return len(arr.valList)
}

func (arr *JSONArrayImpl) add(elem Value) {
	arr.valList = append(arr.valList, elem)
}

func (arr *JSONArrayImpl) set(index int, elem Value) {
	arr.valList[index] = elem
}

func (arr *JSONArrayImpl) getElem(index int) Value {
	return arr.valList[index]
}

func (arr *JSONArrayImpl) sub(start, end int) Value {
	return array(arr.valList[start:end])
}

func (arr *JSONArrayImpl) checkOutofIndex(index int) bool {
	return index < 0 || index >= len(arr.valList)
}

func (arr *JSONArrayImpl) values() []Value {
	return arr.valList
}

func (arr *JSONArrayImpl) String() string {
	return arr.toJSONArrayString()
}

func (arr *JSONArrayImpl) Pretty() string {
	uglyBody := arr.toJSONArrayString()
	var out bytes.Buffer
	err := json.Indent(&out, []byte(uglyBody), "", "  ")
	if err != nil {
		panic(err)
	}
	return out.String()
}
func (arr *JSONArrayImpl) Pr() {
	fmt.Println(arr.Pretty())
}

func (arr *JSONArrayImpl) toJSONArrayString() string {
	var res bytes.Buffer
	res.WriteString("[")
	for i, item := range arr.valList {
		rawVal := toJsonStrVal(item)
		if i < 1 {
			res.WriteString(fmt.Sprintf("%v", rawVal))
		} else {
			res.WriteString(fmt.Sprintf(", %v", rawVal))
		}
	}
	res.WriteString("]")
	return res.String()
}
func (arr *JSONArrayImpl) Sort() Value {
	size := len(arr.valList)
	tmpMap := make(map[string]Value)
	tmpArr := make([]string, 0, size)
	for _, v := range arr.valList {
		vStr := v.String()
		tmpMap[vStr] = v
		tmpArr = append(tmpArr, vStr)
	}
	quickSort(tmpArr, 0, size-1)
	res := make([]Value, 0, size)
	for _, str := range tmpArr {
		res = append(res, tmpMap[str])
	}
	return array(res)
}

func quickSort(arr []string, lower, upper int) {
	if upper <= lower {
		return
	}

	i := lower - 1
	pivot := arr[upper]
	for j := lower; j < upper; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[upper] = arr[upper], arr[i+1]
	partitionIndex := i + 1

	quickSort(arr, lower, partitionIndex-1)
	quickSort(arr, partitionIndex+1, upper)
}
func (arr *JSONArrayImpl) Reverse() Value {
	size := len(arr.valList)
	tmpArr := make([]interface{}, 0, size)
	for i := size - 1; i >= 0; i-- {
		tmpArr = append(tmpArr, arr.valList[i])
	}
	return newQKValue(tmpArr)
}

func (arr *JSONArrayImpl) indexs() []interface{} {
	var res []interface{}
	for i := range arr.valList {
		res = append(res, i)
	}
	return res
}

func (arr *JSONArrayImpl) getItem(index interface{}) Value {
	i := index.(int)
	return arr.valList[i]
}

func (arr *JSONArrayImpl) Add(args []interface{}) {
	for _, arg := range args {
		arr.add(newQKValue(arg))
	}
}
func (arr *JSONArrayImpl) Remove(index int) {
	assert(arr.checkOutofIndex(index), "array out of index")
	newList := make([]Value, 0, arr.Size())
	newList = append(newList, arr.valList[:index]...)
	if index+1 < arr.Size() {
		newList = append(newList, arr.valList[index+1:]...)
	}
	arr.valList = newList
}
func (arr *JSONArrayImpl) Join(seperator string) string {
	vals := arr.values()
	var res bytes.Buffer
	for i, val := range vals {
		if i > 0 {
			res.WriteString(seperator)
		}
		valStr := fmt.Sprintf("%v", val.val())
		res.WriteString(valStr)
	}
	return res.String()
}

func (arr *JSONArrayImpl) val() interface{} {
	return arr
}
func (arr *JSONArrayImpl) isJsonArray() bool {
	return true
}
