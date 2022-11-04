package question

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/gpa/types"
)

type QuestionService struct {
	model.BaseDao
}

func (service QuestionService) ListQuestion() []types.IEntity {

	return dao.Find(singleton.Orm(), &model.Question{}).List()
}
