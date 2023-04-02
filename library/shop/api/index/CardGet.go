package index

import (
	"encoding/base64"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/activity"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/gpa/types"
	"github.com/nbvghost/tool"
	"github.com/skip2/go-qrcode"
)

type CardGet struct {
	CardItemService     activity.CardItemService
	VerificationService order.VerificationService
	User                *model.User `mapping:""`
	Get                 struct {
		CardItemID types.PrimaryKey `uri:"CardItemID"`
	} `method:"get"`
}

func (m *CardGet) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	//CardItemID, _ := strconv.ParseUint(context.PathParams["CardItemID"], 10, 64)
	//CardItemID := object.ParseUint(context.PathParams["CardItemID"])
	//var cardItem model.CardItem
	cardItem := dao.GetByPrimaryKey(singleton.Orm(), entity.CardItem, types.PrimaryKey(m.Get.CardItemID)).(*model.CardItem)

	dao.DeleteBy(singleton.Orm(), &model.Verification{}, map[string]interface{}{
		"UserID":     m.User.ID,
		"CardItemID": cardItem.Primary(),
		"StoreID":    0,
		"Quantity":   0,
	})

	verification := model.Verification{}
	verification.CardItemID = cardItem.Primary()
	verification.UserID = m.User.ID
	verification.VerificationNo = tool.UUID()
	verification.Name, verification.Label = cardItem.GetNameLabel(singleton.Orm())

	dao.Create(singleton.Orm(), &verification)

	results := make(map[string]interface{})

	if false {
		//wxconfig := controller.Wx.MiniProgram(user.OID)
		//postData := make(map[string]interface{})

		/*access_token := controller.Wx.GetAccessToken(wxconfig.ID)
		postData["scene"] = verification.VerificationNo
		//postData["page"] = "pages/store_verification/store_verification"
		postData["width"] = 430
		postData["auto_color"] = true

		body := strings.NewReader(util.StructToJSON(postData))
		//postData := url.Values{}
		//postData.Add("scene","sdfsd")
		resp, err := http.Post("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token="+access_token, "application/json", body)
		if err != nil {
			return &gweb.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: err.Error(), Data: nil}}
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return &gweb.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: err.Error(), Data: nil}}
		}

		defer resp.Body.Close()

		imageString := "data:image/png;base64," + base64.StdEncoding.EncodeToString(b)

		results["QRCodeBase64"] = imageString*/
	} else {

		png, _ := qrcode.Encode(verification.VerificationNo, qrcode.Low, 256)
		imageString := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
		results["QRCodeBase64"] = imageString

	}

	results["Verification"] = verification
	results["CardItem"] = cardItem
	results["ExpireTime"] = cardItem.ExpireTime
	results["HasQuantity"] = cardItem.Quantity - cardItem.UseQuantity

	if (cardItem.Quantity - cardItem.UseQuantity) <= 0 {
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: "数量不足，无法核销", Data: nil}}, nil
	} else {
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: results}}, nil
	}

}
