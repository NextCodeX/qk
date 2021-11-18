package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

func intToRunes(raw int) []rune {
	s := fmt.Sprint(raw)
	var res []rune
	for _, ch := range s {
		res = append(res, ch)
	}
	return res
}

// 字符串转数值类型
func strToNumber(raw string) interface{} {
	if numI, errI := strconv.Atoi(raw); errI == nil {
		return numI
	}

	if numF, errF := strconv.ParseFloat(raw, 64); errF == nil {
		return numF
	}

	return -1
}

func toInt(any interface{}) int {
	switch v := any.(type) {
	case int32:
		return int(v)
	case int64:
		return int(v)
	case int:
		return v
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		i, err := strconv.Atoi(v)
		assert(err != nil, err, "Value:", any)
		return i
	case Value:
		return toInt(v.val())
	default:
		runtimeExcption("failed to int value", any)
	}
	return -1
}

// QK Value 转 go 类型bool
func toBoolean(raw Value) bool {
	if raw == nil || raw.isNULL() {
		return false
	}
	if raw.isInt() {
		return raw.val().(int64) != 0
	} else if raw.isFloat() {
		return raw.val().(float64) != 0
	} else if raw.isBoolean() {
		return raw.val().(bool)
	} else if raw.isString() {
		return raw.val().(string) != ""
	} else if raw.isJsonArray() {
		return raw.val() != nil
	} else if raw.isJsonObject() {
		return raw.val() != nil
	} else if raw.isAny() || raw.isObject() {
		return raw.val() != nil
	} else {
		runtimeExcption("toBoolean: unknown value type: ", raw)
		return false
	}
}

func intToBytes(raw int64) []byte {
	data := uint64(raw)
	res := make([]byte, 8)
	for i := 0; i < 8; i++ {
		res[7-i] = uint8(data >> (8 * i))
	}
	return res
}

func bytesToInt(bs []byte) int64 {
	size := len(bs)
	if size > 8 {
		return 0
	}
	var res uint64
	for i := 0; i < size; i++ {
		var num uint64
		if size < 8 {
			num = uint64(bs[size-1-i])
		} else {
			num = uint64(bs[7-i])
		}
		res = res | (num << (i * 8))
	}
	return int64(res)
}

func floatToBytes(raw float64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, raw)
	// 这里可以继续往buf里写, 都存在buf里
	// binary.Write(buf, binary.LittleEndian, uint16(12345))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return buf.Bytes()
}

func bytesToFloat(data []byte) float64 {
	data = fix8bit(data)
	var res float64
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &res)
	// 这里可以继续读出来存在变量里, 这样就可以解码出来很多, 读的次序和变量类型要对
	// binary.Read(buf, binary.LittlEndian, &v2)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return res
}

// 字节数组如果不足8位，则进行补全
func fix8bit(data []byte) []byte {
	bsLen := len(data)
	if bsLen < 8 {
		res := make([]byte, 0, 8)
		for i := 0; i < 8-bsLen; i++ {
			res = append(res, 0)
		}
		for i := 0; i < bsLen; i++ {
			res = append(res, data[i])
		}
		return res
	}

	return data
}
