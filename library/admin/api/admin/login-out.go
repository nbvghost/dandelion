package admin

import (
	"net/http"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/server/redis"
)

type LoginOut struct {
	Admin *model.Admin `mapping:""`
}

func (m *LoginOut) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	context.Redis().Del(context, redis.NewTokenKey(context.Token()))
	contextValue := contexext.FromContext(context)

	http.SetCookie(contextValue.Response, &http.Cookie{Name: "token", Path: "/", MaxAge: -1})

	//return &gweb.RedirectToUrlResult{Url: "/admin/"}, err
	return result.NewSuccess("账号已经安全退出"), nil
}
