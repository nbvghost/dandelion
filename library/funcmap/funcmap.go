package funcmap

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"html/template"

	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb/conf"
	"github.com/nbvghost/tool/object"
)

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

type FuncObject struct {
	funcMap template.FuncMap
	c       constrain.IContext
}

func NewFuncMap(context constrain.IContext) template.FuncMap {
	fm := &FuncObject{}
	fm.c = context
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
	fm.funcMap["Test"] = fm.test

	for funcName := range regMap {
		func(funcName string) {
			//闭包
			function := regMap[funcName]

			v := reflect.ValueOf(function).Elem()
			functionType := v.Type()

			argsIn := make([]reflect.Type, 0)

			argsIndex := make([]int, 0)

			numField := functionType.NumField()
			for i := 0; i < numField; i++ {
				if _, ok := functionType.Field(i).Tag.Lookup("arg"); ok {
					argsIn = append(argsIn, functionType.Field(i).Type)
					argsIndex = append(argsIndex, i)
				}
			}

			var makeFuncType reflect.Type
			switch function.(type) {
			case IFunc:
				makeFuncType = reflect.FuncOf(argsIn, []reflect.Type{reflect.TypeOf(new(interface{})).Elem()}, false)
			case IWidget:

				makeFuncType = reflect.FuncOf(argsIn, []reflect.Type{reflect.TypeOf(new(interface{})).Elem()}, false)
			}

			backCallFunc := reflect.MakeFunc(makeFuncType, func(args []reflect.Value) (results []reflect.Value) {
				for i := 0; i < len(args); i++ {
					v.Field(argsIndex[i]).Set(args[i])
				}

				var result interface{}

				contextValue := contexext.FromContext(context)

				err := contextValue.Mapping.Before(context, function)
				if err != nil {
					return []reflect.Value{reflect.ValueOf(err)}
				}

				switch function.(type) {
				case IWidget:
					resultData, err := function.(IWidget).Render(fm.c)
					if err != nil {
						return []reflect.Value{reflect.ValueOf(err)}
					}
					//todo resultData["Query"] = fm.c.Query()
					fileName := filepath.Join("view", contextValue.DomainName, "template", "widget", fmt.Sprintf("%s.%s", funcName, "gohtml"))
					var b []byte
					b, err = ioutil.ReadFile(fileName)
					if err != nil {
						return []reflect.Value{reflect.ValueOf(err)}
					}
					var t *template.Template
					t, err = template.New(funcName).Funcs(NewFuncMap(fm.c)).Parse(string(b))
					if err != nil {
						return []reflect.Value{reflect.ValueOf(err)}
					}
					buffer := bytes.NewBuffer(nil)
					if err = t.Execute(buffer, resultData); err != nil {
						return []reflect.Value{reflect.ValueOf(err)}
					}
					result = template.HTML(buffer.Bytes())

				case IFunc:
					result = function.(IFunc).Call(fm.c).Result()
				}

				return []reflect.Value{reflect.ValueOf(result)}
			})
			fm.funcMap[funcName] = backCallFunc.Interface()
		}(funcName)
	}
	return fm.funcMap
}

func (fo *FuncObject) test() map[string]interface{} {

	return map[string]interface{}{"fdsfds": 4545}
}
func (fo *FuncObject) digitAdd(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a+_b, 'f', prec, 64), 64)
	return f
}
func (fo *FuncObject) digitSub(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a-_b, 'f', prec, 64), 64)
	return f
}
func (fo *FuncObject) digitMul(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a*_b, 'f', prec, 64), 64)
	return f
}
func (fo *FuncObject) digitDiv(a, b interface{}, prec int) float64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()
	//f, _ := strconv.ParseFloat(strconv.FormatFloat(_a/_b, 'f', prec, 64), 64)
	f, _ := strconv.ParseFloat(strconv.FormatFloat(_a/_b, 'f', prec, 64), 64)
	return f
}
func (fo *FuncObject) mapFunc(m interface{}, key interface{}) interface{} {
	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Map {
		return v.MapIndex(reflect.ValueOf(key)).Interface()
	}
	panic(fmt.Errorf("Map不能处理%v数据", v.Kind()))
}
func (fo *FuncObject) digitMod(a, b interface{}) uint64 {
	_a := reflect.ValueOf(a).Convert(reflect.TypeOf(float64(0))).Float()
	_b := reflect.ValueOf(b).Convert(reflect.TypeOf(float64(0))).Float()

	///f, _ := strconv.ParseFloat(strconv.FormatFloat(_a%_b, 'f', prec, 64), 64)
	return uint64(_a) % uint64(_b)

}
func (fo *FuncObject) makeArray(len int) []int {

	return make([]int, len)
}
func (fo *FuncObject) urlQueryEncode(source map[string]string) template.URL {
	//fmt.Println(source)
	v := &url.Values{}
	for key := range source {
		v.Set(key, source[key])
	}
	return template.URL(v.Encode())
}
func (fo *FuncObject) html(source string) template.HTML {
	//fmt.Println(source)
	return template.HTML(source)
}
func (fo *FuncObject) dateTimeFormat(source time.Time, format string) string {
	//fmt.Println(source)
	//fmt.Println(format)
	return source.Format(format)
}
func (fo *FuncObject) toJSON(source interface{}) template.JS {
	b, err := json.Marshal(source)
	glog.Error(err)
	return template.JS(b)
}
func (fo *FuncObject) parseInt(source interface{}) int {

	return object.ParseInt(source)
}

func (fo *FuncObject) parseFloat(source interface{}) float64 {
	return object.ParseFloat(source)
}

/*func cipherDecrypter(source string) string {

	str := encryption.CipherDecrypter(encryption.public_PassWord, source)
	return str
}
func cipherEncrypter(source string) string {
	str := encryption.CipherEncrypter(encryption.public_PassWord, source)
	return str
}*/
func (fo *FuncObject) fromJSONToMap(source string) map[string]interface{} {
	d := make(map[string]interface{})
	err := json.Unmarshal([]byte(source), &d)
	glog.Error(err)
	return d
}
func (fo *FuncObject) fromJSONToArray(source string) []interface{} {
	d := make([]interface{}, 0)
	err := json.Unmarshal([]byte(source), &d)
	glog.Error(err)
	return d
}
func (fo *FuncObject) splitFunc(source string, sep string) []string {

	return strings.Split(source, sep)
}
func (fo *FuncObject) includeHTML(url string, params interface{}) template.HTML {
	//util.Trace(params)
	//paramsMap := make(map[string]interface{})

	b := bytes.NewBuffer(make([]byte, 0))
	ww := bufio.NewWriter(b)

	t, err := template.ParseFiles(conf.Config.ViewDir + "/" + url)
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
