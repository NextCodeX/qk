package core

import (
    "fmt"
    "testing"
)

func Test_swi(t *testing.T) {
    a := 3
    a = a ^ 1
    a = a ^ 1
    switch a {
    case 3:
        fmt.Println("a == 3")
        if a > 0 {
            fmt.Println("a > 0")
            break
        }
        fmt.Println("a <= 0")

    }
}

func Test_obj(t *testing.T) {
    v1 := Value{
        t:           IntValue,
        int_value:   99,
        float_value: 0,
        bool_value:  false,
        str_value:   "",
        arr_value:   nil,
        obj_value:   nil,
    }
    v2 := v1
    v2.t = FloatValue
    v2.float_value = 3.14
    fmt.Println("v1:", v1)
    fmt.Println("v2:", v2)
}
