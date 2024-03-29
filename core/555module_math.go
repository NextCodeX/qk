package core

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// return absolute value of tokenList
func (this *InternalFunctionSet) Abs(raw interface{}) (res interface{}) {
	switch num := raw.(type) {
	case int64:
		res = int64(math.Abs(float64(num)))
	case float64:
		res = math.Abs(num)
	default:
		log.Fatal("abs(number) arg type error:", raw)
	}
	return
}

// Pow returns x**y, the base-x exponential of y.
func (this *InternalFunctionSet) Pow(x, y interface{}) interface{} {
	a, ok := toFloat(x)
	assert(!ok, "pow(number, number) arg error:", x, y)
	b, ok := toFloat(y)
	assert(!ok, "pow(number, number) arg error:", x, y)
	return math.Pow(a, b)
}

// Sqrt returns the square root of x.
func (this *InternalFunctionSet) Sqrt(x interface{}) interface{} {
	a, ok := toFloat(x)
	assert(!ok, "Sqrt(number) arg error:", x)
	return math.Sqrt(a)
}

// 四舍五入取整
func (this *InternalFunctionSet) Round(raw float64) interface{} {
	return math.Round(raw)
}

// float number format
func (this *InternalFunctionSet) Fix(raw float64, bitSize int) interface{} {
	dotFormat := "%." + strconv.Itoa(bitSize) + "f"
	tmp := fmt.Sprintf(dotFormat, raw)
	res, err := strconv.ParseFloat(tmp, 64)
	assert(err != nil, "float fix error", raw, bitSize, err)
	return res
}

// string type to number type
func (this *InternalFunctionSet) Number(raw string) interface{} {
	return strToNumber(raw)
}

func toFloat(num interface{}) (float64, bool) {
	if tmp, ok := num.(int64); ok {
		return float64(tmp), ok
	}
	if tmp, ok := num.(float64); ok {
		return tmp, ok
	}
	return 0, false
}

// returns, as an int, a non-negative pseudo-random number in [0,n)
func (this *InternalFunctionSet) Random(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

// returns, as an int, a non-negative pseudo-random number in [n, m)
func (this *InternalFunctionSet) RandomRange(n, m int) int {
	interval := m - n
	rand.Seed(time.Now().UnixNano())
	return n + rand.Intn(interval)
}
