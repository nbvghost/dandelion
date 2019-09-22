package service

import (
	"github.com/nbvghost/dandelion/app/service/dao"

	"github.com/nbvghost/glog"
)

type ManagerService struct {
	dao.BaseDao
}

func (this ManagerService) FindManagerByAccount(Account string) *dao.Manager {
	Orm := dao.Orm()
	manager := &dao.Manager{}
	err := Orm.Where("Account=?", Account).First(manager).Error //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)
	return manager
}
