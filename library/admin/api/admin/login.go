package admin

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/server/httpext"
	"github.com/nbvghost/dandelion/server/redis"
	"github.com/nbvghost/dandelion/service"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/tool/encryption"
)

type Login struct {
	Post struct {
		//AppName  string `form:"AppName"`
		//AppKey   string `form:"AppKey"`
		//Referer  string `header:"Referer"`
		Account  string `form:"account"`
		Password string `form:"password"`
	} `method:"Post"`
}

func (m *Login) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Login) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {

	Orm := db.GetDB(ctx)

	account := strings.ToLower(strings.TrimSpace(m.Post.Account))                           //小写
	password := encryption.Md5ByString(strings.ToLower(strings.TrimSpace(m.Post.Password))) //小写

	admin := service.Admin.Service.FindAdminByAccountAndPassWord(Orm, account, password)

	if admin.IsZero() {
		return nil, result.NewErrorText("密码不正确")
	} else {
		admin.LastLoginAt = time.Now()
		err := dao.UpdateByPrimaryKey(Orm, &model.Admin{}, admin.ID, map[string]any{"LastLoginAt": admin.LastLoginAt})
		if err != nil {
			return nil, err
		}

		session, err := json.Marshal(httpext.Session{
			ID:    fmt.Sprintf("%d", admin.ID),
			Token: ctx.Token(),
		})
		if err != nil {
			return nil, err
		}
		err = ctx.Redis().Set(ctx, redis.NewTokenKey(ctx.Token()), string(session), time.Minute*10)
		if err != nil {
			return nil, err
		}
		return result.NewData(map[string]any{"Admin": admin}), nil
	}
}
