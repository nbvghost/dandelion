package funcmap

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/pkg/errors"

	"github.com/nbvghost/tool/object"
)

var embeds FS

type FS interface {
	fs.ReadDirFS
	fs.ReadFileFS
}

func Embed(v FS) {
	embeds = v
}

type IFuncResult interface {
	Result() interface{}
}

type fakeFuncResult struct {
}

func (m *fakeFuncResult) Result() interface{} {
	return nil
}

type stringFuncResult struct {
	arg string
}

func (m *stringFuncResult) Result() interface{} {
	return m.arg
}

func NewStringFuncResult(arg string) IFuncResult {

	return &stringFuncResult{arg: arg}
}

type result struct {
	data interface{}
}

func (m *result) Result() interface{} {
	return m.data
}
func NewResult(data interface{}) IFuncResult {

	return &result{data: data}
}

type mapFuncResult struct {
	m map[string]interface{}
}

func (m *mapFuncResult) Result() interface{} {
	return m.m
}

func NewMapFuncResult(m map[string]interface{}) IFuncResult {

	return &mapFuncResult{m: m}
}

type stringArrayFuncResult struct {
	args []string
}

func (m *stringArrayFuncResult) Result() interface{} {
	if len(m.args) == 0 {
		return []string{}
	}
	return m.args
}

func NewStringArrayFuncResult(args []string) IFuncResult {

	return &stringArrayFuncResult{args: args}
}

type IFunc interface {
	Call(ctx constrain.IContext) IFuncResult
}

/*var FunctionMap = template.FuncMap{
	"IncludeHTML":     includeHTML,
	"Split":           splitFunc,
	"FromJSONToMap":   fromJSONToMap,
	"FromJSONToArray": fromJSONToArray,
	"ParseFloat":      parseFloat,
	"ParseInt":        parseInt,
	"ToJSON":          toJSON,
	"DateTimeFormat":  dateTimeFormat,
	"HTML":            html,
	"UrlQueryEncode":  urlQueryEncode,
	"DigitAdd":        digitAdd,
	"DigitSub":        digitSub,
	"DigitMul":        digitMul,
	"DigitDiv":        digitDiv,
	"MakeArray":       makeArray,
	"DigitMod":        digitMod,
	//"CipherDecrypter": cipherDecrypter,
	//"CipherEncrypter": cipherEncrypter,
}*/

/*func FuncMap() template.FuncMap {

	return FunctionMap
}*/

var regMap = make(map[string]interface{})

func RegisterFunction(funcName string, function IFunc) {
	if _, ok := regMap[funcName]; ok {
		log.Fatalln(errors.New(fmt.Sprintf("%v函数已经存在", funcName)))
	}

	regMap[funcName] = function

}
func RegisterWidget(funcName string, widget IWidget) {
	if _, ok := regMap[funcName]; ok {
		log.Fatalln(errors.New(fmt.Sprintf("%v函数已经存在", funcName)))
	}

	regMap[funcName] = widget

}

type ITemplateFunc interface {
	Build(context constrain.IContext) template.FuncMap
}
type templateFuncMap struct {
	funcMap template.FuncMap
}

func (fo *templateFuncMap) Build(context constrain.IContext) template.FuncMap {
	createFunc := func(funcName string) {
		function := regMap[funcName]

		v := reflect.ValueOf(function).Elem()
		functionType := v.Type()

		argsIn := make([]reflect.Type, 0)

		argsIndex := make([]int, 0)

		variadicArgIndex := -1
		numField := functionType.NumField()
		for i := 0; i < numField; i++ {
			if tagV, ok := functionType.Field(i).Tag.Lookup("arg"); ok {
				argType := functionType.Field(i).Type
				if strings.Contains(tagV, "...") {
					variadicArgIndex = i
					if argType.Kind() != reflect.Slice {
						panic(errors.Errorf("不定参数%s，必须是Slice类型", argType.Name()))
					}
				}
				argsIn = append(argsIn, argType)
				argsIndex = append(argsIndex, i)
			}
		}
		if variadicArgIndex > -1 && len(argsIn) > 0 && variadicArgIndex != len(argsIn)-1 {
			panic(errors.Errorf("不定参数%s，必须是最后一个参数", argsIn[variadicArgIndex].Name()))
		}
		var variadic bool
		if variadicArgIndex > -1 {
			variadic = true
		}

		var makeFuncType reflect.Type
		switch function.(type) {
		case IFunc:
			makeFuncType = reflect.FuncOf(argsIn, []reflect.Type{reflect.TypeOf(new(interface{})).Elem()}, variadic)
		case IWidget:
			makeFuncType = reflect.FuncOf(argsIn, []reflect.Type{reflect.TypeOf(new(interface{})).Elem()}, variadic)
		}

		contextValue := contexext.FromContext(context)
		contextValue.Mapping.Before(context, function)

		backCallFunc := reflect.MakeFunc(makeFuncType, func(args []reflect.Value) (results []reflect.Value) {
			for i := 0; i < len(args); i++ {
				v.Field(argsIndex[i]).Set(args[i])
			}
			var result interface{}
			var err error

			switch function.(type) {
			case IWidget:
				var resultData map[string]interface{}
				resultData, err = function.(IWidget).Render(context)
				if err != nil {
					return []reflect.Value{reflect.ValueOf(err)}
				}
				//todo resultData["Query"] = fm.c.Query()
				fileName := filepath.Join("view", contextValue.DomainName, "template", "widget", fmt.Sprintf("%s.%s", funcName, "gohtml"))
				var b []byte
				b, err = ioutil.ReadFile(fileName)
				if err != nil {
					b, err = embeds.ReadFile(fmt.Sprintf("template/%s.gohtml", funcName))
					if err != nil {
						return []reflect.Value{reflect.ValueOf(err)}
					}
				}
				var t *template.Template
				t, err = template.New(funcName).Funcs(NewFuncMap().Build(context)).Parse(string(b))
				if err != nil {
					return []reflect.Value{reflect.ValueOf(err)}
				}
				buffer := bytes.NewBuffer(nil)
				if err = t.Execute(buffer, resultData); err != nil {
					return []reflect.Value{reflect.ValueOf(err)}
				}
				result = template.HTML(buffer.Bytes())
			case IFunc:
				result = function.(IFunc).Call(context).Result()
			}

			return []reflect.Value{reflect.ValueOf(result)}
		})
		fo.funcMap[funcName] = backCallFunc.Interface()
	}
	for funcName := range regMap {
		createFunc(funcName)
	}
	return fo.funcMap
}

func NewFuncMap() ITemplateFunc {
	fm := &templateFuncMap{}
	fm.funcMap = make(template.FuncMap)
	fm.funcMap["IncludeHTML"] = fm.includeHTML
	fm.funcMap["Split"] = fm.splitFunc
	fm.funcMap["FromJSONToMap"] = fm.fromJSONToMap
	fm.funcMap["FromJSONToArray"] = fm.fromJSONToArray
	fm.funcMap["ParseFloat"] = fm.parseFloat
	fm.funcMap["ParseInt"] = fm.parseInt
	fm.funcMap["ToJSON"] = fm.toJSON
	fm.funcMap["DateTimeFormat"] = fm.dateTimeFormat
	fm.funcMap["HTML"] = fm.html
	fm.funcMap["UrlQueryEncode"] = fm.urlQueryEncode
	fm.funcMap["DigitAdd"] = fm.digitAdd
	fm.funcMap["DigitSub"] = fm.digitSub
	fm.funcMap["DigitMul"] = fm.digitMul
	fm.funcMap["DigitDiv"] = fm.digitDiv
	fm.funcMap["MakeArray"] = fm.makeArray
	fm.funcMap["DigitMod"] = fm.digitMod
	fm.funcMap["Map"] = fm.mapFunc
	fm.funcMap["Index"] = fm.Index
	fm.funcMap["Empty"] = fm.empty
	return fm
}
func indirectInterface(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Interface {
		return v
	}
	if v.IsNil() {
		return reflect.Value{}
	}
	return v.Elem()
}
func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
	}
	return v, false
}
func indexArg(index reflect.Value, cap int) (int, error) {
	var x int64
	switch index.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x = index.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x = int64(index.Uint())
	case reflect.Invalid:
		return 0, fmt.Errorf("cannot index slice/array with nil")
	default:
		return 0, fmt.Errorf("cannot index slice/array with type %s", index.Type())
	}
	if x < 0 || int(x) < 0 || int(x) > cap {
		return 0, fmt.Errorf("index out of range: %d", x)
	}
	return int(x), nil
}
func (fo *templateFuncMap) empty(v any) bool {
	if v == nil {
		return true
	}
	item := reflect.ValueOf(v)
	if item.IsZero() {
		return true
	}
	if item.Kind() == reflect.Ptr {
		if item.IsNil() {
			return true
		}
	}
	if !item.IsValid() {
		return true
	}
	return false
}
func (fo *templateFuncMap) Index(item reflect.Value, index reflect.Value) (reflect.Value, error) {
	index = indirectInterface(index)
	if index.Int() < 0 {
		index = reflect.ValueOf(item.Len() + int(index.Int()))
	}

	if index.Int() < 0 {
		return reflect.Value{}, nil
	}

	var isNil bool
	if item, isNil = indirect(item); isNil {
		return reflect.Value{}, fmt.Errorf("index of nil pointer")
	}
	switch item.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		x, err := indexArg(index, item.Len())
		if err != nil {
			return reflect.Value{}, err
		}
		item = item.Index(x)
	default:
		return reflect.Value{}, fmt.Errorf("can't index item of type %s", item.Type())
	}
	return item, nil
}
func (fo *templateFuncMap) digitAdd(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a+_b, 'f', prec, 64), 64)
	return f
}
func (fo *templateFuncMap) digitSub(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a-_b, 'f', prec, 64), 64)
	return f
}
func (fo *templateFuncMap) digitMul(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a*_b, 'f', prec, 64), 64)
	return f
}
func (fo *templateFuncMap) digitDiv(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	//f, _ := strconv.ParseFloat(strconv.FormatFloat(_a/_b, 'f', prec, 64), 64)
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a/_b, 'f', prec, 64), 64)
	return f
}
func (fo *templateFuncMap) mapFunc(m interface{}, key interface{}) interface{} {
	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Map {
		if !v.MapIndex(reflect.ValueOf(key)).IsValid() {
			return reflect.New(v.Type().Elem()).Elem().Interface()
		}
		return v.MapIndex(reflect.ValueOf(key)).Interface()
	}
	panic(errors.Errorf("Map不能处理%v数据", v.Kind()))
}
func (fo *templateFuncMap) digitMod(a, b interface{}) uint64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()

	///f, _ := strconv.ParseFloat(strconv.FormatFloat(_a%_b, 'f', prec, 64), 64)
	return uint64(_a) % uint64(_b)

}
func (fo *templateFuncMap) makeArray(len int) []int {

	return make([]int, len)
}
func (fo *templateFuncMap) urlQueryEncode(source map[string]string) template.URL {
	//fmt.Println(source)
	v := &url.Values{}
	for key := range source {
		v.Set(key, source[key])
	}
	return template.URL(v.Encode())
}
func (fo *templateFuncMap) html(source string) template.HTML {
	//fmt.Println(source)
	return template.HTML(source)
}
func (fo *templateFuncMap) dateTimeFormat(source time.Time, format string) string {
	//fmt.Println(source)
	//fmt.Println(format)
	return source.Format(format)
}
func (fo *templateFuncMap) toJSON(source interface{}) template.JS {
	b, err := json.Marshal(source)
	if err != nil {
		log.Println(err)
	}
	return template.JS(b)
}
func (fo *templateFuncMap) parseInt(source interface{}) int {

	return object.ParseInt(source)
}

func (fo *templateFuncMap) parseFloat(source interface{}) float64 {
	return object.ParseFloat(source)
}

/*
func cipherDecrypter(source string) string {

		str := encryption.CipherDecrypter(encryption.public_PassWord, source)
		return str
	}

	func cipherEncrypter(source string) string {
		str := encryption.CipherEncrypter(encryption.public_PassWord, source)
		return str
	}
*/
func (fo *templateFuncMap) fromJSONToMap(source string) map[string]interface{} {
	d := make(map[string]interface{})
	err := json.Unmarshal([]byte(source), &d)
	if err != nil {
		log.Println(err)
	}
	return d
}
func (fo *templateFuncMap) fromJSONToArray(source string) []interface{} {
	d := make([]interface{}, 0)
	err := json.Unmarshal([]byte(source), &d)
	if err != nil {
		log.Println(err)
	}
	return d
}
func (fo *templateFuncMap) splitFunc(source string, sep string) []string {

	return strings.Split(source, sep)
}
func (fo *templateFuncMap) includeHTML(url string, params interface{}) template.HTML {
	//util.Trace(params)
	//paramsMap := make(map[string]interface{})

	b := bytes.NewBuffer(make([]byte, 0))
	ww := bufio.NewWriter(b)

	t, err := template.ParseFiles("view/" + url)
	if os.IsNotExist(err) {
		ww.WriteString("IncludeHTML:not found path in:" + url)
		t = template.New("static")
	} else {
		t.Execute(ww, params)
	}

	//checkError(err, "read from file template")

	ww.Flush()
	//template.JSEscape()
	//template.HTMLEscapeString()

	//	util.Trace(string(b.Bytes()))
	///return string(b.Bytes());
	return template.HTML(string(b.Bytes()))
}
