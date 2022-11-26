package com

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func IndexOf(s []string, v string) int {
	for i, s2 := range s {
		if strings.EqualFold(s2, v) {
			return i
		}
	}
	return -1
}
func CreateStructByMap(m map[string]any, structType any) any {
	st := reflect.TypeOf(structType)
	if st.Kind() == reflect.Ptr {
		st = st.Elem()
	}
	newStruct := reflect.New(st)
	for k := range m {
		field := newStruct.Elem().FieldByName(k)
		if field.IsValid() {
			field.Set(reflect.ValueOf(m[k]))
		}
	}
	return newStruct.Interface()

}

// Diff 只比较struct 里元字段类型的字段，如果字段是struct则不进行比较，from,to 都是(*struct)
func Diff(fromPtr, toPtr any, ignoreField ...string) map[string]any {
	ma := make(map[string]any)
	vfrom := reflect.ValueOf(fromPtr).Elem()
	tfrom := vfrom.Type()
	vto := reflect.ValueOf(toPtr).Elem()
	tto := vfrom.Type()

	if tfrom.Kind() != tto.Kind() {
		//panic(fmt.Errorf("from kind %v to kind %v,error,kind is not match", tfrom.Kind(), tto.Kind()))
		log.Println(fmt.Errorf("from kind %v to kind %v,error,kind is not match", tfrom.Kind(), tto.Kind()))
		return ma
	}
	if tfrom.Kind() == reflect.Map {

		keys := vto.MapKeys()
		for i := 0; i < len(keys); i++ {
			tt := vto.MapIndex(keys[i])
			tv := tt.Interface()
			fv := vfrom.MapIndex(keys[i]).Interface()
			log.Println(tv, fv)

			k := keys[i].Interface().(string)
			if IndexOf(ignoreField, k) != -1 {
				continue
			}

			switch tt.Type().Kind() {
			case reflect.Slice:
				if !reflect.DeepEqual(fv, tv) {
					ma[k] = tv
				}
			default:
				if tt.Type().Kind() != reflect.Struct {
					if !reflect.DeepEqual(fv, tv) {
						ma[k] = tv
					}
				}
			}
		}

		return ma
	}

	for i := 0; i < tfrom.NumField(); i++ {
		field := tfrom.Field(i)
		fieldKind := field.Type.Kind()
		fieldName := field.Name
		if IndexOf(ignoreField, fieldName) != -1 {
			continue
		}
		fv := vfrom.Field(i).Interface()
		tv := vto.Field(i).Interface()
		switch fieldKind {
		case reflect.Slice:
			if !reflect.DeepEqual(fv, tv) {
				ma[fieldName] = tv
			}
		default:
			if fieldKind != reflect.Struct {
				if fv != tv {
					ma[fieldName] = tv
				}
			}
		}

	}
	return ma
}
