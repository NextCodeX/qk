package core

import (
	"testing"
	"fmt"
)

func TestNum(t *testing.T)  {
	fmt.Println(1%2)
	fmt.Println(2%2)
	fmt.Println(3%2)
	fmt.Println(4%2)
	fmt.Println(5%2)
}

func TestRuntimeException(t *testing.T)  {
	fmt.Println("step1")
	runtimeExcption("unknow operation:", "false", "-", "true")
	fmt.Println("final")
}
