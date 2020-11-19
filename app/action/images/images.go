package images

import (
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service/file"
	"github.com/nbvghost/dandelion/app/service/wechat"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/util"

	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type Controller struct {
	gweb.BaseController
	File file.FileService
	Wx   wechat.WxService
}

func (controller *Controller) Init() {
	//controller.Interceptors.DisableManagerSession = true
	//i.Interceptors.Add(&InterceptorFile{})
	controller.AddHandler(gweb.ALLMethod("captcha", controller.captchaAction))
	controller.AddHandler(gweb.GETMethod("/miniprogram/qrcode", controller.miniprogramQRcodeAction))
}
func (controller *Controller) captchaAction(context *gweb.Context) gweb.Result {
	buf := util.CreateCaptchaCodeBytes(play.SessionCaptcha)
	return &gweb.ImageBytesResult{Data: buf}
}
func (controller *Controller) miniprogramQRcodeAction(context *gweb.Context) gweb.Result {
	//user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	//context.Request.ParseForm()
	//Page := context.Request.FormValue("Page")
	Page := context.Request.URL.Query().Get("Page")
	UserID, _ := strconv.ParseUint(context.Request.URL.Query().Get("UserID"), 10, 64)
	MyShareKey := util.EncodeShareKey(UserID, 0)

	ProductID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ProductID"), 10, 64)
	if ProductID != 0 {
		MyShareKey = util.EncodeShareKey(UserID, ProductID)
	}

	wxconfig := controller.Wx.MiniProgram()

	access_token := controller.Wx.GetAccessToken(wxconfig)

	postData := make(map[string]interface{})
	//results := make(map[string]interface{})

	postData["scene"] = MyShareKey
	postData["page"] = Page
	postData["width"] = 600
	postData["auto_color"] = true

	body := strings.NewReader(util.StructToJSON(postData))
	//postData := url.Values{}
	//postData.Add("scene","sdfsd")
	resp, err := http.Post("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token="+access_token, "application/json", body)
	if err != nil {
		return &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionFail, Message: err.Error(), Data: nil}}
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionFail, Message: err.Error(), Data: nil}}
	}
	//fmt.Println(string(b))
	defer resp.Body.Close()

	path := tool.WriteTempFile(b, "image/png")
	return &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionOK, Message: "", Data: path}}
	//return &gweb.ImageBytesResult{Data:b,ContentType:"image/png"}
	//imageString := "data:image/png;base64," + base64.StdEncoding.EncodeToString(b)
	//results["QRCodeBase64"] = imageString
	//return &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionOK, Message: "", Data: results}}

}
