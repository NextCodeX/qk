package core

import (
	"log"
	"fmt"
	"crypto/rand"
)

func stringModuleInit()  {
	s := &QkString{}
	collectFunctionInfo(&s, "str")
}

type QkString struct {
	
}

func (s *QkString) RawUUID() string {
	bs := make([]byte, 16)
	_, err := rand.Read(bs)
	if err != nil {
		log.Fatal(err)
	}
	rawuuid := fmt.Sprintf("%x-%x-%x-%x-%x", bs[0:4], bs[4:6], bs[6:8], bs[8:10], bs[10:])
	return rawuuid
}

func (s *QkString) Uuid() string {
	bs := make([]byte, 16)
	_, err := rand.Read(bs)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x", bs)
	return uuid
}

func (s *QkString) Format(args []interface{}) string {
	assert(len(args) < 2,"function str_format(format, any...) must has two parameters.")
	format, ok := args[0].(string)
	assert(!ok, "function str_format(format, any...), parameter format must be string type.")
	return fmt.Sprintf(format, args[1:]...)
}