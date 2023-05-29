package user

import (
	"encoding/base64"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/library/db"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/user"
	"github.com/pkg/errors"
)

type Update struct {
	UserService user.UserService
	Post        struct {
		Name           string `form:"Name"`
		Portrait       string `form:"Portrait"`
		PortraitBase64 string `form:"PortraitBase64"`
		Gender         int    `form:"Gender"`
	} `method:"Post"`
	User *model.User `mapping:""`
}

func (m *Update) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	changeMap := map[string]any{}
	if len(m.Post.Name) > 0 {
		changeMap["Name"] = m.Post.Name
	}
	if len(m.Post.Portrait) > 0 {
		changeMap["Portrait"] = m.Post.Portrait
	}
	if m.Post.Gender > 0 {
		changeMap["Gender"] = m.Post.Gender
	}
	if len(m.Post.PortraitBase64) > 0 {
		decodeString, err := base64.StdEncoding.DecodeString(m.Post.PortraitBase64)
		if err != nil {
			return nil, err
		}
		avatar, err := oss.UploadAvatar(ctx, m.User.ID, decodeString)
		if err != nil {
			return nil, err
		}
		if avatar.Code != 0 {
			return nil, errors.New(avatar.Message)
		}
		changeMap["Portrait"], err = oss.ReadUrl(ctx, avatar.Data.Path)
		if err != nil {
			return nil, err
		}
	}

	if len(changeMap) > 0 {
		err := dao.UpdateByPrimaryKey(db.Orm(), entity.User, m.User.ID, changeMap)
		if err != nil {
			return &result.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: err.Error(), Data: nil}}, err
		}
	}

	user := dao.GetByPrimaryKey(db.Orm(), entity.User, m.User.ID)
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: user}}, nil

}
func (m *Update) Handle(context constrain.IContext) (r constrain.IResult, err error) {

	//TODO implement me
	panic("implement me")
}
