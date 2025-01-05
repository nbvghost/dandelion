package account

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"strings"
)

type Register struct {
	Post struct {
		Account  string `form:"account"`
		Password string `form:"password"`
	} `method:"Post"`
}

func (m *Register) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Register) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {

	tx := db.Orm().Begin()

	if haveAdmin := service.Admin.Service.FindAdminByAccount(tx, strings.TrimSpace(m.Post.Account)); !haveAdmin.IsZero() {
		tx.Rollback()
		return nil, result.NewErrorText("这个账号已经存在")
	}

	if _, err := service.Admin.Service.InitOrganizationInfo(m.Post.Account, m.Post.Password); err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &result.JsonResult{Data: result.NewSuccess("注册成功")}, nil
}
