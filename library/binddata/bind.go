package binddata

import (
	"github.com/go-playground/validator/v10"
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

var validate = validator.New()

/*func Bind(apiHandler interface{}, desc *serviceobject.GrpcRequest) (interface{}, error) {
	var err error

	apiHandlerValue := reflect.ValueOf(apiHandler)
	handlerType := apiHandlerValue.Elem().Type()

	{
		numField := handlerType.NumField()
		for i := 0; i < numField; i++ {
			if tag, ok := handlerType.Field(i).Tag.Lookup("method"); ok {

				if strings.EqualFold(desc.HttpMethod, tag) {
					objValuePtr := reflect.New(apiHandlerValue.Elem().Field(i).Type())
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

								log.Println(jsonValuePtr.Interface())

								vh := apiHandlerValue.Elem().Field(i).Field(ii)
								if vh.Kind() == reflect.Ptr {
									objValuePtr.Elem().Field(ii).Set(jsonValuePtr)
								} else {
									objValuePtr.Elem().Field(ii).Set(jsonValuePtr.Elem())
								}
								log.Println(objValuePtr.Elem().Field(ii).Interface())
								//log.Println(jsonValuePtr.Interface())
							}
						}
					}
					if err = validate.Struct(obj); err != nil {
						return nil, action.NewCodeWithError(action.ValidateError, err)
					}

					log.Println(apiHandlerValue.Interface())
					log.Println(objValuePtr.Interface())
					log.Println(obj)
					//log.Println(obj)
					//log.Println(handlerValue.Interface())

					vh := apiHandlerValue.Elem().Field(i)
					if vh.Kind() == reflect.Ptr {
						apiHandlerValue.Elem().Field(i).Set(objValuePtr)
					} else {
						apiHandlerValue.Elem().Field(i).Set(objValuePtr.Elem())
					}
					log.Println(apiHandlerValue.Interface())
					return apiHandlerValue.Interface(), nil

				}

			}
		}
	}

	return apiHandlerValue.Interface(), nil
}
*/
