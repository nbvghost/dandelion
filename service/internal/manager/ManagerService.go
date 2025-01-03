package manager

import (
	"github.com/nbvghost/dandelion/library/db"
	"log"

	"github.com/nbvghost/dandelion/entity/model"
)

type ManagerService struct {
	model.BaseDao
}

func (this ManagerService) FindManagerByAccount(Account string) *model.Manager {
	Orm := db.Orm()
	manager := &model.Manager{}
	err := Orm.Where(`"Account"=?`, Account).First(manager).Error //SelectOne(user, "select * from User where Email=?", Email)
	log.Println(err)
	return manager
}
