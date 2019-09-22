package service

import "github.com/nbvghost/dandelion/app/service/dao"

type QuestionService struct {
	dao.BaseDao
}

func (service QuestionService) ListQuestion() []dao.Question {
	var questions []dao.Question
	service.FindAll(dao.Orm(), &questions)
	return questions
}
