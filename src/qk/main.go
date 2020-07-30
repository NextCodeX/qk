package main

import (
    "core"
    "fmt"
    "reflect"
)

type Person map[string]interface{}

func main()  {
    core.Run()

    var m Person
    mtype := reflect.TypeOf(m)
    fmt.Println("map assert: ", mtype, mtype.Kind() == reflect.Map)
}


