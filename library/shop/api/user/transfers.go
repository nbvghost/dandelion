package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service"
)

type Transfers struct {
	Post struct {
		ReUserName string `form:"ReUserName"`
	} `method:"Post"`
	User         *model.User         `mapping:""`
	WechatConfig *model.WechatConfig `mapping:""`
}

func (m *Transfers) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	contextValue := contexext.FromContext(ctx)
	IP := util.GetIP(contextValue.Request)
	err := service.Order.Transfers.UserTransfers(ctx, m.User.ID, m.Post.ReUserName, IP, m.WechatConfig)

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "提现成功，请查看到账通知结果", nil)}, nil
}
func (m *Transfers) Handle(context constrain.IContext) (r constrain.IResult, err error) {

	//TODO implement me
	panic("implement me")
}
