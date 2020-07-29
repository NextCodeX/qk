package module

import (
	"reflect"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"bufio"
	"encoding/json"
)

func init()  {
	f := &File{}
	v1 := reflect.ValueOf(&f).Elem()
	k1 := v1.Type()
	for i := 0; i < v1.NumMethod(); i++ {
		key := k1.Method(i)
		val := v1.Method(i)
		fmt.Println(key.Name, val.Type(), val.Interface())
	}
}

type File struct {

}

func (f *File) Bytes(filename string) []byte {
	bs, err := ioutil.ReadFile(filename)
	if err == nil {
		return bs
	}
	log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
	return nil
}

func (f *File) Content(filename string) string {
	bs, err := ioutil.ReadFile(filename)
	if err == nil {
		return string(bs)
	}
	log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
	return ""
}

func (f *File) Lines(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
		return nil
	}
	var res []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	return res
}

func (f *File) Json(filename string) map[string]interface{} {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to read %v file: %v", filename, err.Error()))
		return nil
	}
	var tmp interface{}
	err = json.Unmarshal(bs, &tmp)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to parse json[%v]: %v", string(bs), err.Error()))
		return nil
	}
	res := tmp.(map[string]interface{})
	return res
}