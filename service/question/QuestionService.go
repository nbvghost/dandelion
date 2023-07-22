package question

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type QuestionService struct {
	model.BaseDao
}

func (service QuestionService) ListQuestion() []dao.IEntity {

	return dao.Find(db.Orm(), &model.Question{}).List()
}
