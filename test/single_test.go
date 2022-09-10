package test

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"qk/core"
	"reflect"
	"testing"
)

func Test_demo(t *testing.T) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		errorLog(err)
	// 	}
	// }()
	demo, _ := filepath.Abs("../examples/type_str_enchance.qk")
	//demo, _ := filepath.Abs("../examples/demo.qk")
	// core.DEBUG = true
	core.TestFlag = true
	core.ExecScriptFile(demo)
}

func Test_api(t *testing.T) {
	coronaVirusJSON := `{
        "name" : "covid-11",
        "country" : "China",
        "city" : "Wuhan",
        "reason" : "Non vedge Food"
    }`

	// Declared an empty map interface
	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	err := json.Unmarshal([]byte(coronaVirusJSON), &result)
	if err != nil {
		return
	}

	// Print the data type of result variable
	fmt.Println(reflect.TypeOf(result), result)

	var arr []any
	arrStr := `[1, true, "check"]`
	err = json.Unmarshal([]byte(arrStr), &arr)
	if err != nil {
		return
	}
	fmt.Println(reflect.TypeOf(result), arr)

}
