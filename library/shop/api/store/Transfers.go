package store

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service"
)

type Transfers struct {
	Store        *model.Store        `mapping:""`
	User         *model.User         `mapping:""`
	WechatConfig *model.WechatConfig `mapping:""`
	Post         struct {
		ReUserName string `form:"ReUserName"`
	} `method:"Post"`
}

func (m *Transfers) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	contextValue := contexext.FromContext(context)

	IP := util.GetIP(contextValue.Request)
	err := service.Order.Transfers.StoreTransfers(m.Store.ID, m.User.ID, m.Post.ReUserName, IP, m.WechatConfig)

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "提现申请成功，请查看到账通知结果", nil)}, nil
}

func (m *Transfers) Handle(context constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")

}
