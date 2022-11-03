package manager

import (
	"log"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
)

type ManagerService struct {
	model.BaseDao
}

func (this ManagerService) FindManagerByAccount(Account string) *model.Manager {
	Orm := singleton.Orm()
	manager := &model.Manager{}
	err := Orm.Where(`"Account"=?`, Account).First(manager).Error //SelectOne(user, "select * from User where Email=?", Email)
	log.Println(err)
	return manager
}
