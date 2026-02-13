package account

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type GetLoginUser struct {
	User *model.User `mapping:""`
}

func (g *GetLoginUser) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	results := make(map[string]interface{})
	results["User"] = g.User
	results["MyShareKey"] = service.Wechat.WXQRCodeParams.EncodeShareKey(ctx, ctx.UID(), 0) //tool.Hashids{}.Encode(user.ID) //tool.CipherEncrypterData(strconv.Itoa(int(user.ID)))
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: results}}, nil

}
