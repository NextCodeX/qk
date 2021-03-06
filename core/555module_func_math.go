package core

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func (mr *ModuleRegister) MathModuleInit() {
	math := &QKMath{}
	fmath := collectFunctionInfo(&math)
	functionRegister("", fmath)
}


type QKMath struct {

}


// return absolute value of raw
func (mt *QKMath) Abs(raw interface{}) (res interface{}) {
	switch num := raw.(type) {
	case int:
		res = int(math.Abs(float64(num)))
	case float64:
		res = math.Abs(num)
	default:
		log.Fatal("abs(number) arg type error:", raw)
	}
	return
}

// Pow returns x**y, the base-x exponential of y.
func (mt *QKMath) Pow(x, y interface{}) interface{} {
	a, ok  := toFloat(x)
	assert(!ok, "pow(number, number) arg error:", x, y)
	b, ok  := toFloat(y)
	assert(!ok, "pow(number, number) arg error:", x, y)
	return math.Pow(a, b)
}

// Sqrt returns the square root of x.
func (mt *QKMath) Sqrt(x interface{}) interface{} {
	a, ok  := toFloat(x)
	assert(!ok, "Sqrt(number) arg error:", x)
	return math.Sqrt(a)
}

// 四舍五入取整
func (mt *QKMath) Round(raw float64) interface{} {
	return math.Round(raw)
}

//  float number format
func (mt *QKMath) Fix(raw float64, bitSize int) interface{} {
	dotFormat := "%." + strconv.Itoa(bitSize) + "f"
	tmp := fmt.Sprintf(dotFormat, raw)
	res, err := strconv.ParseFloat(tmp, 64)
	assert(err != nil, "float fix error", raw, bitSize, err)
	return res
}

// string type to number type
func (mt *QKMath) Number(raw string) interface{} {
	numI, errI := strconv.Atoi(raw)
	if errI == nil {
		return numI
	}
	numF, errF := strconv.ParseFloat(raw, 64)
	assert(errF != nil, "number(string) error!", raw)
	return numF
}

func toFloat(num interface{}) (float64, bool){
	if tmp, ok := num.(int); ok {
		return float64(tmp), ok
	}
	if tmp, ok := num.(float64); ok {
		return tmp, ok
	}
	return 0, false
}

// returns, as an int, a non-negative pseudo-random number in [0,n)
func (mt *QKMath) Random(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

// returns, as an int, a non-negative pseudo-random number in [n, m)
func (mt *QKMath) RandomRange(n, m int) int {
	interval := m - n
	rand.Seed(time.Now().UnixNano())
	return n + rand.Intn(interval)
}
