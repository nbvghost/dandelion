package gobext

import (
	"encoding/gob"
	"fmt"
	"reflect"
)

var structMap = make(map[string]reflect.Type)

// GetStructName from encoding/gob/type.go:836
func GetStructName(value interface{}) string {
	rt := reflect.TypeOf(value)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	var name string

	if rt.PkgPath() == "" {
		name = rt.Name()
	} else {
		name = rt.PkgPath() + "." + rt.Name()
	}

	return name
}
func NewGob(name string) interface{} {
	if v, ok := structMap[name]; !ok {
		panic(fmt.Errorf("不存在的gob对象:%s", name))
	} else {
		return reflect.New(v).Interface()
	}
}
func Register(value interface{}) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	name := GetStructName(value)

	if _, ok := structMap[name]; ok {
		panic(fmt.Errorf("已经存在的gob对象:%s", name))
	}

	structMap[name] = v.Type()
	gob.RegisterName(name, value)
}
