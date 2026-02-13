package account

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/server/httpext"
	"github.com/nbvghost/dandelion/server/redis"
	"github.com/nbvghost/dandelion/service"
	"strings"
	"sync"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/tool/encryption"
)

type Login struct {
	Post struct {
		//AppName  string `form:"AppName"`
		//AppKey   string `form:"AppKey"`
		Referer  string `header:"Referer"`
		Account  string `form:"account"`
		Password string `form:"password"`
	} `method:"Post"`
}

func (m *Login) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

var once = &sync.Once{}

func (m *Login) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {

	contextValue := contexext.FromContext(ctx)

	Orm := db.GetDB(ctx)

	account := strings.TrimSpace(m.Post.Account) //小写

	user := service.Admin.Service.FindAdminByAccount(Orm, account)

	if user.ID == 0 {
		return nil, result.NewErrorText("账号不存在或密码不正确，请重试")
	} else {
		md5Password := encryption.Md5ByString(strings.TrimSpace(m.Post.Password))
		if strings.EqualFold(user.PassWord, md5Password) {
			err := dao.UpdateByPrimaryKey(Orm, &model.Admin{}, user.ID, &model.Admin{LastLoginAt: time.Now()})
			if err != nil {
				return nil, err
			}

			Redirect := contextValue.Request.URL.Query().Get("Redirect")

			session, err := json.Marshal(httpext.Session{
				ID:    fmt.Sprintf("%d", user.ID),
				Token: ctx.Token(),
			})
			if err != nil {
				return nil, err
			}
			err = ctx.Redis().Set(ctx, redis.NewTokenKey(ctx.Token()), string(session), time.Minute*10)
			if err != nil {
				return nil, err
			}

			return &result.ActionResult{Data: map[string]any{
				"Redirect": Redirect,
			}, Code: 0, Message: "登陆成功"}, nil

		} else {
			return nil, result.NewErrorText("账号不存在或密码不正确，请重试")
		}
	}
}
