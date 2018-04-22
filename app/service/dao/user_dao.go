package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb/tool"
)

type UserDao = BaseDao

func (UserDao) FindUserByTel(Orm *gorm.DB, Tel string) *User {
	user := &User{}
	err := Orm.Where("Tel=?", Tel).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	tool.CheckError(err)
	return user
}

func (UserDao) FindUserByOpenID(Orm *gorm.DB, CompanyID uint64, OpenID string) *User {

	user := &User{}
	//CompanyOpenID := user.GetCompanyOpenID(CompanyID, OpenID)
	err := Orm.Where("OpenID=?", OpenID).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	tool.CheckError(err)
	return user
}
