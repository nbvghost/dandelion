package account

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/file"
	"github.com/nbvghost/dandelion/service/wechat"
	"io/ioutil"
	"net/http"
	"strings"
)

type MiniprogramQRcode struct {
	WXQRCodeParamsService wechat.WXQRCodeParamsService
	WxService             wechat.WxService
	FileService           file.FileService
	//WechatConfig          *model.WechatConfig `mapping:""`
	Get struct {
		Page      string         `form:"Page"`
		UserID    dao.PrimaryKey `form:"UserID"`
		ProductID dao.PrimaryKey `form:"ProductID"`
	} `method:"Get"`
}

func (g *MiniprogramQRcode) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {

	return nil, nil
}

func (g *MiniprogramQRcode) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	user := dao.GetByPrimaryKey(db.Orm(), entity.User, g.Get.UserID).(*model.User)
	wechatConfig := g.WxService.MiniProgramByOID(db.Orm(), user.OID)

	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*entity.Organization)
	//context.Request.ParseForm()
	//Page := context.Request.FormValue("Page")
	//Page := context.Request.URL.Query().Get("Page")
	//UserID, _ := strconv.ParseUint(context.Request.URL.Query().Get("UserID"), 10, 64)
	//Page := object.ParseUint(context.Request.URL.Query().Get("Page"))
	//UserID := object.ParseUint(context.Request.URL.Query().Get("UserID"))

	MyShareKey := g.WXQRCodeParamsService.EncodeShareKey(g.Get.UserID, 0)

	//ProductID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ProductID"), 10, 64)
	//ProductID := object.ParseUint(context.Request.URL.Query().Get("ProductID"))
	if g.Get.ProductID != 0 {
		MyShareKey = g.WXQRCodeParamsService.EncodeShareKey(g.Get.UserID, uint(g.Get.ProductID))
	}

	accessToken := g.WxService.GetAccessToken(wechatConfig)

	postData := make(map[string]interface{})
	//results := make(map[string]interface{})

	postData["scene"] = MyShareKey
	postData["page"] = g.Get.Page
	postData["width"] = 600
	postData["auto_color"] = true

	body := strings.NewReader(util.StructToJSON(postData))
	//postData := url.Values{}
	//postData.Add("scene","sdfsd")
	resp, err := http.Post("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token="+accessToken, "application/json", body)
	if err != nil {
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: err.Error(), Data: nil}}, nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: err.Error(), Data: nil}}, nil
	}
	//fmt.Println(string(b))
	defer resp.Body.Close()

	path, err := g.FileService.WriteTempFile(b, "image/png")
	if err != nil {
		return nil, err
	}
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: path}}, nil
	//return &gweb.ImageBytesResult{Data:b,ContentType:"image/png"}
	//imageString := "data:image/png;base64," + base64.StdEncoding.EncodeToString(b)
	//results["QRCodeBase64"] = imageString
	//return &gweb.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: results}}

}
