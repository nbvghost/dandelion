package service

import "dandelion/app/service/dao"

type ManagerService struct {
	dao dao.ManagerDao
}

func (self ManagerService) FindManagerByAccount(Account string) *dao.Manager {

	return self.dao.FindManagerByAccount(Orm, Account)
}
