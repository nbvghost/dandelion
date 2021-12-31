package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"reflect"
)

type User struct {
	Nm string
	U  int
	F  float64
}

func main() {

	u := &User{
		Nm: "xx",
		U:  5455,
		F:  2.10,
	}

	buffer := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buffer).Encode(u)
	log.Println(err, string(buffer.Bytes()))

	gob.Register(map[string]interface{}{})

	var ms = reflect.New(reflect.TypeOf(struct{}{}))

	err = gob.NewDecoder(buffer).DecodeValue(ms)
	log.Println(err, ms.Interface())

}
