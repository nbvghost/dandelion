package mapping

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/pkg/errors"
	"log"

	"github.com/nbvghost/gpa/types"

	"reflect"
)

type IMapping interface {
	Call(context constrain.IContext) (instance types.IEntity)
	Name() string
	Instance() types.IEntity
}

type mapping struct {
	poolList map[string]map[string]IMapping
}

func (m *mapping) register(mapping IMapping) error {
	path := util.GetPkgPath(mapping.Instance())
	if _, ok := m.poolList[path]; !ok {
		m.poolList[path] = make(map[string]IMapping)
	}

	if _, ok := m.poolList[path][mapping.Name()]; ok {
		return errors.Errorf("已经注册的instance(%s)", path)
	}
	m.poolList[path][mapping.Name()] = mapping
	return nil
}
func (m *mapping) Before(context constrain.IContext, handler interface{}) {
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
			if mp, has := m.poolList[name][tag]; has {
				instance := mp.Call(context)
				if instance == nil {
					//return errors.New("mapping的Call不能返回空的实例")
					log.Printf("%v的Call返回空的数据", mp)
					instance = mp.Instance()
				}
				mappingValue := reflect.ValueOf(instance)
				if fieldV.Kind() == reflect.Ptr {
					if mappingValue.Kind() == reflect.Ptr {
						fieldV.Set(mappingValue)
					} else {
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
	return
}

func (m *mapping) ViewAfter(context constrain.IContext, r constrain.IViewResult) error {
	m.Before(context, r)
	return nil
}
func (m *mapping) AddMapping(mapping IMapping) {
	if err := m.register(mapping); err != nil {
		panic(err)
	}
}

func New(mappings ...IMapping) constrain.IMappingCallback {
	v := &mapping{poolList: make(map[string]map[string]IMapping)}
	for index := range mappings {
		mapping := mappings[index]
		if err := v.register(mapping); err != nil {
			panic(err)
		}
	}
	return v
}
