package user

import (
	"io"
	"mime/multipart"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/pkg/errors"
)

type UploadAvatar struct {
	Post struct {
		File *multipart.FileHeader `form:"file"`
		//UserID dao.PrimaryKey        `form:"uid"`
	} `method:"Post"`
	User *model.User `mapping:""`
}

func (m *UploadAvatar) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return nil, nil
}
func (m *UploadAvatar) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	f, err := m.Post.File.Open()
	if err != nil {
		return nil, err
	}
	if m.User.ID == 0 {
		return nil, errors.New("数据错误")
	}
	fBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	changeMap := map[string]any{}
	avatar, err := oss.UploadAvatar(ctx, m.User.OID, m.User.ID, fBytes)
	if err != nil {
		return nil, err
	}
	if avatar.Code != 0 {
		return nil, errors.New(avatar.Message)
	}
	changeMap["Portrait"], err = oss.ReadUrl(ctx, avatar.Data.Path)

	if len(changeMap) > 0 {
		err := dao.UpdateByPrimaryKey(db.GetDB(ctx), entity.User, m.User.ID, changeMap)
		if err != nil {
			return &result.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: err.Error(), Data: nil}}, err
		}
	}
	user := dao.GetByPrimaryKey(db.GetDB(ctx), entity.User, m.User.ID)
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "OK", Data: map[string]any{"User": user}}}, nil
}
