package route

import (
	"bytes"
	"encoding/json"
	"github.com/nbvghost/dandelion/library/action"
	"github.com/nbvghost/dandelion/service/serviceobject"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// interface{} as a Number instead of as a float64.
var EnableDecoderUseNumber = false

// EnableDecoderDisallowUnknownFields is used to call the DisallowUnknownFields method
// on the JSON Decoder instance. DisallowUnknownFields causes the Decoder to
// return an error when the destination is a struct and the input contains object
// keys which do not match any non-ignored, exported fields in the destination.
var EnableDecoderDisallowUnknownFields = false

func Bind(handlerType reflect.Type, desc *serviceobject.GrpcRequest) (interface{}, error) {
	var err error

	handlerValue := reflect.New(handlerType)

	{
		numField := handlerType.NumField()
		for i := 0; i < numField; i++ {
			if tag, ok := handlerType.Field(i).Tag.Lookup("method"); ok {

				if strings.EqualFold(desc.HttpMethod, tag) {
					objValuePtr := reflect.New(handlerValue.Elem().Field(i).Type())
					obj := objValuePtr.Interface()
					{
						//path params uri
						m := make(map[string][]string)
						for k := range desc.Uri {
							m[k] = []string{desc.Uri[k]}
						}
						if err = mapUri(obj, m); err != nil {
							return nil, err
						}
					}
					{
						//query
						var values url.Values
						values, err = url.ParseQuery(desc.Query)
						if err != nil {
							return nil, err
						}
						if err = mapForm(obj, values); err != nil {
							return nil, err
						}
					}
					{
						//header
						header := http.Header{}

						var values url.Values
						values, err = url.ParseQuery(desc.Header)
						if err != nil {
							return nil, err
						}

						for k := range values {
							l := len(values[k])
							if l > 1 {
								for p := 0; p < l; p++ {
									header.Add(k, values[k][p])
								}
							} else {
								header.Add(k, values.Get(k))
							}
						}
						if err = mapHeader(obj, header); err != nil {
							return nil, err
						}
					}
					{
						//form
						var values url.Values
						values, err = url.ParseQuery(desc.Form)
						if err != nil {
							return nil, err
						}
						if err = mapForm(obj, values); err != nil {
							return nil, err
						}
					}
					{
						fieldNum := handlerType.Field(i).Type.NumField()
						for ii := 0; ii < fieldNum; ii++ {
							if _, ok := handlerType.Field(i).Type.Field(ii).Tag.Lookup("body"); ok {
								//json

								jsonType := handlerType.Field(i).Type.Field(ii).Type
								if jsonType.Kind() == reflect.Ptr {
									jsonType = jsonType.Elem()
								}
								jsonValuePtr := reflect.New(jsonType)

								decoder := json.NewDecoder(bytes.NewBuffer(desc.Body))
								if EnableDecoderUseNumber {
									decoder.UseNumber()
								}
								if EnableDecoderDisallowUnknownFields {
									decoder.DisallowUnknownFields()
								}
								if err = decoder.Decode(jsonValuePtr.Interface()); err != nil {
									return nil, err
								}

								if err = setStructFieldValue(jsonValuePtr, handlerValue.Elem().Field(i).Field(ii)); err != nil {
									return nil, err
								}
							}
						}
					}
					if err = validate.Struct(obj); err != nil {
						return nil, action.NewCodeWithError(action.ValidateError, err)
					}

					if err = setStructFieldValue(objValuePtr, handlerValue.Elem().Field(i)); err != nil {
						return nil, err
					}
					return handlerValue.Interface(), nil

				}

			}
		}
	}

	return handlerValue.Interface(), nil
}
func setStructFieldValue(from reflect.Value, target reflect.Value) error {
	if from.Kind() == reflect.Ptr && target.Kind() == reflect.Ptr {
		target.Set(from)
		return nil
	}
	if from.Kind() == reflect.Ptr {
		from = from.Elem()
	}
	if target.Kind() == reflect.Ptr {
		target = target.Elem()
	}

	for i := 0; i < from.NumField(); i++ {
		v := from.Field(i)
		if v.IsZero() {
			//TODO 判断其它条件？
			continue
		}
		if v.Kind() == reflect.Struct {
			if err := setStructFieldValue(v, target.Field(i)); err != nil {
				return err
			}
			continue
		}
		target.Field(i).Set(v)
	}
	return nil
}
