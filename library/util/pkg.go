package util

import "reflect"

// GetPkgPath from encoding/gob/type.go:836
func GetPkgPath(value interface{}) string {
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
