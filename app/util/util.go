package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/nbvghost/captcha"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
)

type Map map[string]string
type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m Map) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}
func (m *Map) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = Map{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

func RequestBodyToJSON(body io.ReadCloser, target interface{}) error {
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, target)
	return err
}
func JSONToStruct(body string, target interface{}) error {
	err := json.Unmarshal([]byte(body), target)
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
func StructToJSON(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(b)
}
func CreateCaptchaCodeBytes(SessionCaptcha string) []byte {
	//d := captcha.RandomDigits(5)
	//w := captcha.NewImage(data.Session_Captcha, d, captcha.StdWidth, captcha.StdHeight)
	//context.Session.Attributes.Put(data.Session_Captcha, d)
	buf := bytes.NewBuffer(make([]byte, 0))
	//w.WriteTo(context.Response)
	//w.WriteTo(buf)

	captcha.NewLenByID(5, SessionCaptcha)

	captcha.WriteImage(buf, SessionCaptcha, captcha.StdWidth, captcha.StdHeight)

	return buf.Bytes()
}
func GetHost(Request *http.Request) string {

	if Request.TLS == nil {
		return "http://" + Request.Host
	} else {
		return "https://" + Request.Host
	}
}
func GetFullPath(Request *http.Request) string {

	redirect := ""
	if len(Request.URL.Query().Encode()) == 0 {
		redirect = Request.URL.Path
	} else {
		redirect = Request.URL.Path + "?" + Request.URL.Query().Encode()
	}
	return url.QueryEscape(redirect)
}
func GetFullUrl(Request *http.Request) string {

	if Request.TLS == nil {
		return "http://" + Request.Host + Request.RequestURI
	} else {
		return "https://" + Request.Host + Request.RequestURI
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
	//fmt.Println(context.Request)
	//fmt.Println(context.Request.Header.Get("X-Forwarded-For"))
	//fmt.Println(context.Request.RemoteAddr)
	//Ali-Cdn-Real-Ip
	IP := context.Request.Header.Get("Ali-Cdn-Real-Ip")
	if strings.EqualFold(IP, "") {
		//_IP := context.Request.Header.Get("X-Forwarded-For")

		IP = strings.Split(context.Request.Header.Get("X-Forwarded-For"), ",")[0]
		if strings.EqualFold(IP, "") {
			text := context.Request.RemoteAddr
			if strings.Contains(text, "::") {
				IP = "0.0.0.0"
			} else {
				IP = strings.Split(text, ":")[0]
			}
		}
	}
	return IP
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
func Rounding45(rounding float64, prec int) float64 {

	f, err := strconv.ParseFloat(strconv.FormatFloat(rounding, 'f', prec, 64), 64)
	glog.Error(err)
	return f
	//strconv.ParseFloat(strconv.FormatFloat(float64(45454)/float64(100),'f',5,64),64)
	//return math.Floor(rounding+0.5)
}
func EncodeShareKey(UserID, ProductID uint64) string {

	buffer := bytes.NewBuffer(make([]byte, 0))
	binary.Write(buffer, binary.LittleEndian, &UserID)
	binary.Write(buffer, binary.LittleEndian, &ProductID)

	da := hex.EncodeToString(buffer.Bytes())
	//0123456789abcdef
	var hashkey = "10a29f38b45e7c6d"
	var hashData = ""
	for _, value := range da {
		switch string(value) {
		case "a":
			hashData += string(hashkey[10])
		case "b":
			hashData += string(hashkey[11])
		case "c":
			hashData += string(hashkey[12])
		case "d":
			hashData += string(hashkey[13])
		case "e":
			hashData += string(hashkey[14])
		case "f":
			hashData += string(hashkey[15])
		default:
			index, _ := strconv.Atoi(string(value))
			hashData += string(hashkey[index])

		}
	}
	//fmt.Println(hashData)

	//base32.NewEncoding("")
	return hashData
}
func DecodeShareKey(ShareKey string) (UserID, ProductID uint64) {

	var hashkey = "10a29f38b45e7c6d"
	var baseStr = "0123456789abcdef"

	//65c3421a11111111aa391b1911111111

	var readyData = ""
	for _, value := range ShareKey {

		for index, hashkeyValue := range hashkey {

			if strings.EqualFold(string(value), string(hashkeyValue)) {

				readyData += string(baseStr[index])

				break
			}

		}

	}

	b, err := hex.DecodeString(readyData)
	glog.Error(err)

	buffer := bytes.NewBuffer(b)
	binary.Read(buffer, binary.LittleEndian, &UserID)
	binary.Read(buffer, binary.LittleEndian, &ProductID)

	//_ShareKey, _ := url.QueryUnescape(ShareKey)
	//arrs := strings.Split(_ShareKey, ",")
	//UserID, _ = strconv.ParseUint(arrs[0], 10, 64)
	//ProductID, _ = strconv.ParseUint(arrs[1], 10, 64)
	return
}
