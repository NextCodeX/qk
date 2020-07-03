package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "strings"
    "testing"
)


func Test_json1(t *testing.T) {
    str := `[1, true, "testing"]`
    var v interface{}
    err := json.Unmarshal([]byte(str), &v)
    fmt.Printf("%v -> %T \n", v, v)
    fmt.Println(err)
}


func Test_genTypeAssert(t *testing.T) {
    const input = `
    Identifier 
    Str
    Int
    Float
    Symbol
`
    typename := "Token"
    var buf bytes.Buffer
    scanner := bufio.NewScanner(strings.NewReader(input))
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        buf.WriteString(fmt.Sprintf("func (this *%v) is%v() bool {\n", typename, line))
        buf.WriteString(fmt.Sprintf("\treturn (this.t & %v) == %v \n", line, line))
        buf.WriteString("}\n\n")
    }
    fmt.Println(buf.String())
}

func Test_print(t *testing.T) {
    fmt.Printf("%v", "hi, \"joker\" \n")
    fmt.Println("++++++++++++++++++++++++++")
}
