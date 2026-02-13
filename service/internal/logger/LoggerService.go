package logger

import (
	"context"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

// LoggerService
type LoggerService struct {
	model.BaseDao
}

func (s LoggerService) Error(ctx context.Context, Title, Data string) {
	Orm := db.GetDB(ctx)
	logger := &model.Logger{}
	logger.Title = Title
	logger.Data = Data
	logger.Key = 1
	dao.Create(Orm, logger)
}
func (s LoggerService) Warning(ctx context.Context, Title, Data string) {
	Orm := db.GetDB(ctx)
	logger := &model.Logger{}
	logger.Title = Title
	logger.Data = Data
	logger.Key = 2
	dao.Create(Orm, logger)
}
