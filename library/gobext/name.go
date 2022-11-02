package gobext

import (
	"encoding/gob"
	"reflect"

	"github.com/nbvghost/dandelion/library/util"
	"github.com/pkg/errors"
)

var structMap = make(map[string]reflect.Type)

func NewGob(name string) interface{} {
	if v, ok := structMap[name]; !ok {
		panic(errors.Errorf("不存在的gob对象:%s", name))
	} else {
		return reflect.New(v).Interface()
	}
}
func Register(value interface{}) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	name := util.GetPkgPath(value)
	structMap[name] = v.Type()
	gob.RegisterName(name, value)
}
