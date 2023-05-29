package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/user"
	"github.com/nbvghost/gpa/types"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
)

type UploadAvatar struct {
	UserService user.UserService
	Post        struct {
		File   *multipart.FileHeader `form:"file"`
		UserID types.PrimaryKey      `form:"uid"`
	} `method:"Post"`
}

func (m *UploadAvatar) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return nil, nil
}
func (m *UploadAvatar) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	f, err := m.Post.File.Open()
	if err != nil {
		return nil, err
	}
	if m.Post.UserID == 0 {
		return nil, errors.New("数据错误")
	}
	fBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	changeMap := map[string]any{}
	avatar, err := oss.UploadAvatar(context, m.Post.UserID, fBytes)
	if err != nil {
		return nil, err
	}
	if avatar.Code != 0 {
		return nil, errors.New(avatar.Message)
	}
	changeMap["Portrait"], err = oss.ReadUrl(context, avatar.Data.Path)

	if len(changeMap) > 0 {
		err := dao.UpdateByPrimaryKey(db.Orm(), entity.User, m.Post.UserID, changeMap)
		if err != nil {
			return &result.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: err.Error(), Data: nil}}, err
		}
	}
	user := dao.GetByPrimaryKey(db.Orm(), entity.User, m.Post.UserID)
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: map[string]any{"User": user}}}, nil
}
