package mapping

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/util"

	"github.com/nbvghost/gpa/types"
	"log"
	"reflect"
)

type Call func(context constrain.IContext) (instance types.IEntity)

type IMapping interface {
	Call(context constrain.IContext) (instance types.IEntity)
	Name() string
	Instance() types.IEntity
}

type mapping struct {
	poolList map[string]map[string]Call
}

func (m *mapping) register(instance types.IEntity, call Call, mappingName string) error {
	path := util.GetPkgPath(instance)
	if _, ok := m.poolList[path]; !ok {
		m.poolList[path] = make(map[string]Call)
	}

	if _, ok := m.poolList[path][mappingName]; ok {
		return fmt.Errorf("已经注册的instance(%s)", path)
	}
	m.poolList[path][mappingName] = call
	return nil
}
func (m *mapping) Before(context constrain.IContext, handler interface{}) error {
	t := reflect.TypeOf(handler)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.ValueOf(handler)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		fieldT := t.Field(i)
		if tag, ok := fieldT.Tag.Lookup("mapping"); ok {
			fieldV := v.Field(i)
			name := util.GetPkgPath(fieldV.Interface())
			if call, has := m.poolList[name][tag]; has {
				instance := call(context)
				if instance == nil {
					return errors.New("mapping的Call不能返回空的实例")
				}

				mappingValue := reflect.ValueOf(instance)
				if fieldV.Kind() == reflect.Ptr {
					if mappingValue.Kind() == reflect.Ptr {
						fieldV.Set(mappingValue)
					} else {
						log.Println(mappingValue)
						fieldV.Set(reflect.ValueOf(mappingValue.Interface()).Addr())
					}
				} else {
					if mappingValue.Kind() == reflect.Ptr {
						fieldV.Set(mappingValue.Elem())
					} else {
						fieldV.Set(mappingValue)
					}
				}

			}
		}
	}
	return nil
}

func (m *mapping) ViewAfter(context constrain.IContext, r constrain.IViewResult) error {
	return m.Before(context, r)
}
func (m *mapping) AddMapping(mapping IMapping) {
	if err := m.register(mapping.Instance(), mapping.Call, mapping.Name()); err != nil {
		panic(err)
	}
}

func New(mappings ...IMapping) constrain.IMappingCallback {
	v := &mapping{poolList: make(map[string]map[string]Call)}
	for index := range mappings {
		mapping := mappings[index]
		if err := v.register(mapping.Instance(), mapping.Call, mapping.Name()); err != nil {
			panic(err)
		}
	}
	return v
}
