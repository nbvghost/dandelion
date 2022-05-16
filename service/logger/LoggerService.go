package logger

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
)

//LoggerService
type LoggerService struct {
	model.BaseDao
}

func (s LoggerService) Error(Title, Data string) {
	Orm := singleton.Orm()
	logger := &model.Logger{}
	logger.Title = Title
	logger.Data = Data
	logger.Key = 1
	s.Add(Orm, logger)
}
func (s LoggerService) Warning(Title, Data string) {
	Orm := singleton.Orm()
	logger := &model.Logger{}
	logger.Title = Title
	logger.Data = Data
	logger.Key = 2
	s.Add(Orm, logger)
}
