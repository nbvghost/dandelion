package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"dandelion/app/play"

	"github.com/nbvghost/captcha"
	"github.com/nbvghost/gweb"

	"crypto/sha1"
)

func RequestBodyToJSON(body io.ReadCloser, target interface{}) error {
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, target)
	return err
}
func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
func CreateCaptchaCodeBytes() []byte {
	//d := captcha.RandomDigits(5)
	//w := captcha.NewImage(data.Session_Captcha, d, captcha.StdWidth, captcha.StdHeight)
	//context.Session.Attributes.Put(data.Session_Captcha, d)
	buf := bytes.NewBuffer(make([]byte, 0))
	//w.WriteTo(context.Response)
	//w.WriteTo(buf)

	captcha.NewLenByID(5, play.SessionCaptcha)

	captcha.WriteImage(buf, play.SessionCaptcha, captcha.StdWidth, captcha.StdHeight)

	return buf.Bytes()
}
func GetHost(context *gweb.Context) string {

	if context.Request.TLS == nil {
		return "http://" + context.Request.Host
	} else {
		return "https://" + context.Request.Host
	}
}
func GetFullUrl(context *gweb.Context) string {

	if context.Request.TLS == nil {
		return "http://" + context.Request.Host + context.Request.RequestURI
	} else {
		return "https://" + context.Request.Host + context.Request.RequestURI
	}

}
func IsMobile(context *gweb.Context) bool {
	UserAgent := context.Request.Header.Get("User-Agent")
	//fmt.Println(UserAgent)
	if strings.Contains(strings.ToLower(UserAgent), "mobile") {
		//FrameworkHttp.OutHtmlFileWithPath(context,"game/web/ssc.html")
		return true
	} else {
		return false
	}
}
func GetIP(context *gweb.Context) string {
	IP := context.Request.Header.Get("X-Forwarded-For")
	if strings.EqualFold(IP, "") {
		text := context.Request.RemoteAddr
		if strings.Contains(text, "::") {
			IP = "0.0.0.0"
		} else {
			IP = strings.Split(text, ":")[0]
		}

	}
	return IP
}
func Md5ByString(valeu string) string {
	ddf := md5.New()
	ddf.Write([]byte(valeu))
	md5Str := hex.EncodeToString(ddf.Sum(nil))
	return strings.ToUpper(md5Str)
}
func Md5ByBytes(valeu []byte) string {
	ddf := md5.New()
	ddf.Write(valeu)
	md5Str := hex.EncodeToString(ddf.Sum(nil))
	return strings.ToUpper(md5Str)
}
func SignSha1(m string) string {
	h := sha1.New()
	h.Write([]byte(m))
	//echostr=8588369340523248596&nonce=3174274245&signature=95366bb0dadb2ba70e0027135f2013a771163589&timestamp=1505651604
	//h.Write([]byte(""))
	//io.WriteString(h, "His money is twice tainted:")
	//io.WriteString(h, " 'taint yours and 'taint mine.")
	bs := h.Sum(nil)
	//fmt.Printf("% x", bs)
	return hex.EncodeToString(bs)
}
