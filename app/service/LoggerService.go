package service

import "github.com/nbvghost/dandelion/app/service/dao"

//LoggerService
type LoggerService struct {
	dao.BaseDao
}

func (s LoggerService) Error(Title, Data string) {
	Orm := dao.Orm()
	logger := &dao.Logger{}
	logger.Title = Title
	logger.Data = Data
	logger.Key = 1
	s.Add(Orm, logger)
}
func (s LoggerService) Warning(Title, Data string) {
	Orm := dao.Orm()
	logger := &dao.Logger{}
	logger.Title = Title
	logger.Data = Data
	logger.Key = 2
	s.Add(Orm, logger)
}
