package core

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func (fns *InternalFunctionSet) UuidRaw() string {
	bs := make([]byte, 16)
	_, err := rand.Read(bs)
	if err != nil {
		log.Fatal(err)
	}
	rawuuid := fmt.Sprintf("%x-%x-%x-%x-%x", bs[0:4], bs[4:6], bs[6:8], bs[8:10], bs[10:])
	return rawuuid
}

func (fns *InternalFunctionSet) Uuid() string {
	bs := make([]byte, 16)
	_, err := rand.Read(bs)
	if err != nil {
		runtimeExcption(err)
	}
	uuid := fmt.Sprintf("%x", bs)
	return uuid
}

func (fns *InternalFunctionSet) Key16() string {
	bs := make([]byte, 16)
	_, err := rand.Read(bs)
	if err != nil {
		runtimeExcption(err)
	}
	return base64.StdEncoding.EncodeToString(bs)
}
func (fns *InternalFunctionSet) Key24() string {
	bs := make([]byte, 24)
	_, err := rand.Read(bs)
	if err != nil {
		runtimeExcption(err)
	}
	return base64.StdEncoding.EncodeToString(bs)
}
func (fns *InternalFunctionSet) Key32() string {
	bs := make([]byte, 32)
	_, err := rand.Read(bs)
	if err != nil {
		runtimeExcption(err)
	}
	return base64.StdEncoding.EncodeToString(bs)
}

func (fns *InternalFunctionSet) Fmt(args []interface{}) string {
	assert(len(args) < 2, "function str_format(format, any...) must has two parameters.")
	format, ok := args[0].(string)
	assert(!ok, "function str_format(format, any...), parameter format must be string type.")
	return fmt.Sprintf(format, args[1:]...)
}
