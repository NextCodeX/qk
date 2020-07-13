package core

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "regexp"
    "strings"
    "testing"
)

func Test_totoken(t *testing.T) {
    token := newToken("true", Identifier)
    value := tokenToValue(&token)
    fmt.Println(value.bool_value)
}

func Test_map1(t *testing.T) {
    nv := newVariables()
    addVar(&nv)
    printVar(&nv)
}

func printVar(nv *Variables) {
    vr := nv.get("tk")
    fmt.Println(vr)
}

func addVar(nv *Variables) {
    nv.add(newVar("tk", "unique"))
}

type Obj struct {
    val string
}

func Test_map(t *testing.T) {

    m := make(map[string]*Obj)
    change(m)
    fetch(m)
}

func fetch(m map[string]*Obj) {
    fmt.Println(m["name"])
}

func change(m map[string]*Obj) {
    m["name"] = &Obj{"changlie"}
}

func Test_clone(t *testing.T) {
    f, _ := os.Open("d:/raw.txt")
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if !strings.HasPrefix(line, "private") {
            continue
        }
        re := regexp.MustCompile("private\\s+\\w+\\s+(\\w+);")
        res := re.FindAllStringSubmatch(line, -1)
        fname := res[0][1]
        mt := strings.ToUpper(fname[:1]) + fname[1:]
        fmt.Printf("obj.set%v(this.%v);\n", mt, fname)
    }
}

func Test_fields(t *testing.T) {
    f, _ := os.Open("d:/raw.txt")
    scanner := bufio.NewScanner(f)
    var buf bytes.Buffer
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        re := regexp.MustCompile("`(\\w+)`")
        res := re.FindAllStringSubmatch(line, -1)
        buf.WriteString(res[0][1])
        buf.WriteString(", ")
    }
    fmt.Println(buf.String())
}

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
