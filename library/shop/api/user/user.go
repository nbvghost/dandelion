package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/tool/encryption"
	"github.com/pkg/errors"
	"strings"
)

type User struct {
	User *model.User `mapping:""`
	Get  struct {
	} `method:"Get"`
	Put struct {
		Email           string
		FirstName       string
		LastName        string
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
func (m *User) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {
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

	if len(m.Put.LastName) > 0 && len(m.Put.FirstName) > 0 {
		changeMap["Name"] = m.Put.LastName + " " + m.Put.FirstName
	}

	tx := singleton.Orm().Begin()

	if needValidPassword {
		hasUser := dao.GetByPrimaryKey(tx, &model.User{}, context.UID()).(*model.User)
		if !strings.EqualFold(hasUser.Password, encryption.Md5ByString(strings.TrimSpace(m.Put.CurrentPassword))) {
			tx.Rollback()
			return nil, errors.New("The password doesn't match this account. Verify the password and try again.")
		}
	}

	err = dao.UpdateByPrimaryKey(tx, &model.User{}, context.UID(), changeMap)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = dao.UpdateBy(tx, &model.UserInfo{}, map[string]any{"AllowAssistance": m.Put.AllowAssistance}, `"UserID"=?`, context.UID())
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	user := dao.GetByPrimaryKey(tx, &model.User{}, context.UID())
	userInfo := dao.GetBy(tx, &model.UserInfo{}, map[string]any{"UserID": context.UID()})
	tx.Commit()

	return result.NewData(map[string]any{"User": user, "UserInfo": userInfo}), nil
}
