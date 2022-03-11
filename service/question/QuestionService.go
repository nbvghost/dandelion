package question

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
)

type QuestionService struct {
	model.BaseDao
}

func (service QuestionService) ListQuestion() []model.Question {
	var questions []model.Question
	service.FindAll(singleton.Orm(), &questions)
	return questions
}
