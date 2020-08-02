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

func (s *QkString) Uuid() string {
	bs := make([]byte, 16)
	_, err := rand.Read(bs)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", bs[0:4], bs[4:6], bs[6:8], bs[8:10], bs[10:])
	return uuid
}

func (s *QkString) Rawuuid() string {
	bs := make([]byte, 16)
	_, err := rand.Read(bs)
	if err != nil {
		log.Fatal(err)
	}
	rawuuid := fmt.Sprintf("%x", bs)
	return rawuuid
}