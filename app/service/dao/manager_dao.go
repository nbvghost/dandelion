package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb/tool"
)

type ManagerDao struct{}

func (ManagerDao) FindManagerByAccount(Orm *gorm.DB, Account string) *Manager {
	manager := &Manager{}
	err := Orm.Where("Account=?", Account).First(manager).Error //SelectOne(user, "select * from User where Email=?", Email)
	tool.CheckError(err)
	return manager
}
