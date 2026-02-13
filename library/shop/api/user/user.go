package user

import (
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/tool/encryption"
	"github.com/pkg/errors"
)

type User struct {
	User *model.User `mapping:""`
	Get  struct {
	} `method:"Get"`
	Put struct {
		Email           string
		Name            string
		ChangeEmail     bool
		ChangePassword  bool
		AllowAssistance bool
		CurrentPassword string
		NewPassword     string
	} `method:"Put"`
}

func (m *User) Handle(context constrain.IContext) (r constrain.IResult, err error) {

	return nil, nil
}
func (m *User) HandlePut(ctx constrain.IContext) (r constrain.IResult, err error) {
	changeMap := make(map[string]any)
	needValidPassword := false
	if len(m.Put.Email) > 0 && m.Put.ChangeEmail {
		changeMap["Email"] = m.Put.Email
		needValidPassword = true
	}
	if len(m.Put.NewPassword) > 0 && m.Put.ChangePassword {
		changeMap["Password"] = encryption.Md5ByString(strings.TrimSpace(m.Put.NewPassword))
		needValidPassword = true
	}

	if len(m.Put.Name) > 0 {
		changeMap["Name"] = m.Put.Name
	}

	tx := db.GetDB(ctx).Begin()

	if needValidPassword {
		hasUser := dao.GetByPrimaryKey(tx, &model.User{}, ctx.UID()).(*model.User)
		if !strings.EqualFold(hasUser.Password, encryption.Md5ByString(strings.TrimSpace(m.Put.CurrentPassword))) {
			tx.Rollback()
			return nil, errors.New("The password doesn't match this account. Verify the password and try again.")
		}
	}

	err = dao.UpdateByPrimaryKey(tx, &model.User{}, ctx.UID(), changeMap)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	userInfo := service.User.GetUserInfo(ctx, ctx.UID())
	userInfo.SetAllowAssistance(m.Put.AllowAssistance)
	err = userInfo.Update(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	u := dao.GetByPrimaryKey(tx, &model.User{}, ctx.UID())

	tx.Commit()

	return result.NewData(map[string]any{"User": u, "UserInfo": userInfo.Data()}), nil
}
