package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/nbvghost/captcha"
	"github.com/nbvghost/tool/object"
)

type Hashids struct {
}

func (Hashids) EncodeShareKey(UserID uint) string {
	return "UserID:" + strconv.Itoa(int(UserID))
}
func (Hashids) DecodeShareKey(ShareKey string) uint {
	_ShareKey, _ := url.QueryUnescape(ShareKey)
	SuperiorID := object.ParseUint(strings.Split(_ShareKey, ":")[1])
	return SuperiorID
}

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

type IJSON interface {

}
func JSONToStruct[T IJSON](j string) (T,error) {
	st:=reflect.TypeFor[T]()
	if st.Kind() == reflect.Ptr{
		st =st.Elem()
	}
	t:=reflect.New(st)
	err:=json.Unmarshal([]byte(j),t.Interface())

	if st.Kind()==reflect.Slice{
		return t.Elem().Interface().(T), err
	}else{
		return t.Interface().(T), err
	}
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
		log.Println(err)
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

func GetScheme(request *http.Request) string {
	// Can't use `r.Request.URL.Scheme`
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if request.TLS != nil {
		return "https"
	}
	if scheme := request.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if scheme := request.Header.Get("X-Forwarded-Protocol"); scheme != "" {
		return scheme
	}
	if ssl := request.Header.Get("X-Forwarded-Ssl"); ssl == "on" {
		return "https"
	}
	if scheme := request.Header.Get("X-Url-Scheme"); scheme != "" {
		return scheme
	}
	return "http"
}

func GetHost(request *http.Request) string {
	return fmt.Sprintf("%s://%s", GetScheme(request), request.Host)
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
func GetFullUrl(request *http.Request) string {
	return fmt.Sprintf("%s://%s%s", GetScheme(request), request.Host, request.RequestURI)
}
func IsMobile(request *http.Request) bool {
	UserAgent := request.Header.Get("User-Agent")
	//fmt.Println(UserAgent)
	if strings.Contains(strings.ToLower(UserAgent), "mobile") {
		//FrameworkHttp.OutHtmlFileWithPath(context,"game/web/ssc.html")
		return true
	} else {
		return false
	}
}
func GetIPLocation(ip string) (string, error) {
	response, err := http.Get(fmt.Sprintf("https://opendata.baidu.com/api.php?query=%s&co=&resource_id=6006&oe=utf8", ip))
	if err != nil {
		return "", err
	}
	type IpInfo struct {
		Status       string `json:"status"`
		T            string `json:"t"`
		SetCacheTime string `json:"set_cache_time"`
		Data         []struct {
			ExtendedLocation string `json:"ExtendedLocation"`
			OriginQuery      string `json:"OriginQuery"`
			Appinfo          string `json:"appinfo"`
			DispType         int    `json:"disp_type"`
			Fetchkey         string `json:"fetchkey"`
			Location         string `json:"location"`
			Origip           string `json:"origip"`
			Origipquery      string `json:"origipquery"`
			Resourceid       string `json:"resourceid"`
			RoleId           int    `json:"role_id"`
			ShareImage       int    `json:"shareImage"`
			ShowLikeShare    int    `json:"showLikeShare"`
			Showlamp         string `json:"showlamp"`
			Titlecont        string `json:"titlecont"`
			Tplt             string `json:"tplt"`
		} `json:"data"`
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var ipInfo = IpInfo{}
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return "", err
	}
	if len(ipInfo.Data) == 0 {
		return "", errors.New("没有找到ip的地理位置信息")
	}
	return ipInfo.Data[0].Location, nil

}
func GetIP(request *http.Request) string {
	//fmt.Println(context.Request)
	//fmt.Println(context.Request.Header.Get("X-Forwarded-For"))
	//fmt.Println(context.Request.RemoteAddr)
	//Ali-Cdn-Real-Ip
	IP := request.Header.Get("Ali-Cdn-Real-Ip")
	if strings.EqualFold(IP, "") {
		//_IP := context.Request.Header.Get("X-Forwarded-For")

		IP = strings.Split(request.Header.Get("X-Forwarded-For"), ",")[0]
		if strings.EqualFold(IP, "") {
			text := request.RemoteAddr
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
	log.Println(err)
	return f
	//strconv.ParseFloat(strconv.FormatFloat(float64(45454)/float64(100),'f',5,64),64)
	//return math.Floor(rounding+0.5)
}

/*func EncodeShareKey(UserID, ProductID uint) string {

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
	return hashData
}*/
/*func DecodeShareKey(ShareKey string) (UserID, ProductID uint) {

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
	log.Println(err)

	buffer := bytes.NewBuffer(b)
	binary.Read(buffer, binary.LittleEndian, &UserID)
	binary.Read(buffer, binary.LittleEndian, &ProductID)

	//_ShareKey, _ := url.QueryUnescape(ShareKey)
	//arrs := strings.Split(_ShareKey, ",")
	//UserID, _ = strconv.ParseUint(arrs[0], 10, 64)
	//ProductID, _ = strconv.ParseUint(arrs[1], 10, 64)
	return
}*/
func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile(`\<[\S\s]+?\>`)
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile(`\<style[\S\s]+?\</style\>`)
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile(`\<script[\S\s]+?\</script\>`)
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile(`\<[\S\s]+?\>`)
	src = re.ReplaceAllString(src, "")
	//去除连续的换行符
	re, _ = regexp.Compile(`\s{2,}`)
	src = re.ReplaceAllString(src, "")
	return strings.TrimSpace(src)
}
