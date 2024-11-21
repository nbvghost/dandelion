package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/server/redis"
)

type SignOut struct {
	Organization  *model.Organization  `mapping:""`
	ContentConfig *model.ContentConfig `mapping:""`
}

type SignOutReply struct {
	extends.ViewBase
}

func (m *SignOutReply) GetResult(context constrain.IContext, viewHandler constrain.IViewHandler) constrain.IResult {
	contextValue := contexext.FromContext(context)
	return &result.RedirectToUrlResult{Url: contextValue.Request.Header.Get("Referer")}
}
func (m *SignOut) Render(context constrain.IContext) (r constrain.IViewResult, err error) {
	reply := &SignOutReply{}
	_, err = context.Redis().Del(context, redis.NewTokenKey(context.Token()))
	if err != nil {
		return nil, err
	}
	return reply, nil
}
