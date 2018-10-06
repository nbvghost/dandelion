package service

import (
	"dandelion/app/service/dao"

	"github.com/nbvghost/gweb/tool"
)

type ManagerService struct {
	dao.BaseDao
}

func (this ManagerService) FindManagerByAccount(Account string) *dao.Manager {
	Orm := dao.Orm()
	manager := &dao.Manager{}
	err := Orm.Where("Account=?", Account).First(manager).Error //SelectOne(user, "select * from User where Email=?", Email)
	tool.CheckError(err)
	return manager
}
