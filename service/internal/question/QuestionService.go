package question

import (
	"context"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type QuestionService struct {
	model.BaseDao
}

func (service QuestionService) ListQuestion(ctx context.Context) []dao.IEntity {

	return dao.Find(db.GetDB(ctx), &model.Question{}).List()
}
