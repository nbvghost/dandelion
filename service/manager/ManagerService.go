package manager

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"

	"github.com/nbvghost/glog"
)

type ManagerService struct {
	model.BaseDao
}

func (this ManagerService) FindManagerByAccount(Account string) *model.Manager {
	Orm := singleton.Orm()
	manager := &model.Manager{}
	err := Orm.Where("Account=?", Account).First(manager).Error //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)
	return manager
}
