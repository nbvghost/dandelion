package mapping

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/util"
	"reflect"
)

type Call func(context constrain.IContext) (instance interface{}, err error)

type IMapping interface {
	Call(context constrain.IContext) (instance interface{}, err error)
	Name() string
	Instance() interface{}
}

type mapping struct {
	poolList map[string]map[string]Call
}

func (m *mapping) Register(instance interface{}, call Call, mappingName string) error {
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
				instance, err := call(context)
				if err != nil {
					return err
				}
				fieldV.Set(reflect.ValueOf(instance))
			}
		}
	}
	return nil
}

func (m *mapping) ViewAfter(context constrain.IContext, r constrain.IViewResult) error {
	return m.Before(context, r)
}

func New(mappings ...IMapping) *mapping {
	v := &mapping{poolList: make(map[string]map[string]Call)}
	for index := range mappings {
		mapping := mappings[index]
		if err := v.Register(mapping.Instance(), mapping.Call, mapping.Name()); err != nil {
			panic(err)
		}
	}
	return v
}
