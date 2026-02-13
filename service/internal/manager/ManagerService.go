package manager

import (
	"context"
	"log"

	"github.com/nbvghost/dandelion/library/db"

	"github.com/nbvghost/dandelion/entity/model"
)

type ManagerService struct {
	model.BaseDao
}

func (this ManagerService) FindManagerByAccount(ctx context.Context, Account string) *model.Manager {
	Orm := db.GetDB(ctx)
	manager := &model.Manager{}
	err := Orm.Where(`"Account"=?`, Account).First(manager).Error //SelectOne(user, "select * from User where Email=?", Email)
	log.Println(err)
	return manager
}
