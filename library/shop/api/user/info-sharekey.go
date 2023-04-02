package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/user"
	"github.com/nbvghost/gpa/types"
	"github.com/nbvghost/tool/encryption"
	"github.com/nbvghost/tool/object"
)

type InfoSharekey struct {
	UserService user.UserService
	Post        struct {
		ShareKey string `form:"ShareKey"`
	} `method:"Post"`
}

func (m *InfoSharekey) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {

	//UserID, _ := strconv.ParseUint(encryption.CipherDecrypter(play.GWebSecretKey, ShareKey), 10, 64)
	UserID := object.ParseUint(encryption.CipherDecrypter(play.GWebSecretKey, m.Post.ShareKey))

	//var user model.User
	user := dao.GetByPrimaryKey(singleton.Orm(), entity.User, types.PrimaryKey(UserID))
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: user}}, nil
}
func (m *InfoSharekey) Handle(context constrain.IContext) (r constrain.IResult, err error) {

	//TODO implement me
	panic("implement me")
}
