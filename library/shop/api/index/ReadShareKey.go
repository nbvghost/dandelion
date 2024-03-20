package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type ReadShareKey struct {
	Get struct {
		ShareKey string `form:"ShareKey"`
	} `method:"get"`
}

func (m *ReadShareKey) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*entity.Organization)
	//ShareKey := context.Request.URL.Query().Get("ShareKey")
	UserID, ProductID := service.Wechat.WXQRCodeParams.DecodeShareKey(m.Get.ShareKey)

	Result := make(map[string]interface{})
	Result["UserID"] = UserID
	Result["ProductID"] = ProductID
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: Result}}, nil
}
