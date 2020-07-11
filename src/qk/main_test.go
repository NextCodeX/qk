package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strings"
    "testing"
)

func Test_compile(t *testing.T) {
    var res []interface{}
    res = append(res, 3)
    res = append(res, false)
    res = append(res, "test")
    fmt.Println(res)
    fmt.Println(res...)
    fmt.Println("successfully!")
}



func Test_fhandle2(t *testing.T) {
    tableName := "table_info"
    fmt.Printf("DROP TABLE IF EXISTS `%v`;\n", tableName)
    fmt.Printf("CREATE TABLE `%v` (\n", tableName)
    tfs := getTfs()
    //fmt.Println(len(infos), len(comments))
    for _, tf := range tfs {
        fmt.Printf("%v    COMMENT '%v', \n", tf.f, tf.c)
    }

}

type tf struct {
    f string
    c string
}

func getTfs() []tf {
    var res []tf
    f, _ := os.Open("d:/origin.txt")
    scanner := bufio.NewScanner(f)
    var comment string
    var field string
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if comment == "" && strings.HasPrefix(line, "*") {
            line = strings.TrimSpace(line[1:])
            comment = line
        }
        if comment != "" && field == "" && strings.HasPrefix(line, "private") {
            re := regexp.MustCompile(`\w+\s+(\w+);`)
            arr := re.FindAllStringSubmatch(line, -1)
            //fmt.Println("arr:", arr[0][1])
            field = arr[0][1]

            //handleInfo(field, comment)
            res = append(res, tf{field, comment})
            field = ""
            comment = ""
        }
    }
    return res
}

func Test_fhandle1(t *testing.T) {
    tableName := "table_info"
    fmt.Printf("DROP TABLE IF EXISTS `%v`;\n", tableName)
    fmt.Printf("CREATE TABLE `%v` (\n", tableName)
    infos := fieldInfos()
    comments := getComments()
    //fmt.Println(len(infos), len(comments))
    for i, info := range infos {
      fmt.Printf("%v COMMENT '%v', \n", info, comments[i])
    }

}

func getComments() []string {
    var res []string
    f, _ := os.Open("d:/origin.txt")
    scanner := bufio.NewScanner(f)
    var comment string
    var field string
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if comment == "" && strings.HasPrefix(line, "*") {
            line = strings.TrimSpace(line[1:])
            comment = line
        }
        if comment != "" && field == "" && strings.HasPrefix(line, "private") {
            re := regexp.MustCompile(`\w+\s+(\w+);`)
            arr := re.FindAllStringSubmatch(line, -1)
            //fmt.Println("arr:", arr[0][1])
            field = arr[0][1]

            //handleInfo(field, comment)
            res = append(res, comment)
            field = ""
            comment = ""
        }
    }
    return res
}

func handleInfo(field, comment string) {

    //msg := fmt.Sprintf("alter table %v modify column %v comment '%v';", tableName, field, comment)
    //fmt.Println(field, " - ", comment)
    //fmt.Println(msg)
}

func fieldInfos() []string {
    tableCreate := `
id int primary key auto_increment,
tableName varchar(100) not null,
tableTag varchar(150) not null,
taskId int,
primaryKey varchar(200),
lastUpdatedColumnName varchar(200),
schemaInfo text,
description varchar(2000),
ext varchar(200),
createdDate datetime,
createdBy varchar(150),
lastUpdatedDate datetime,
lastUpdatedBy varchar(150),
`
var res []string
scanner := bufio.NewScanner(strings.NewReader(tableCreate))
for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    if line == "" {
        continue
    }
    res = append(res, line[:len(line)-1])
}
return res
}

