package question

import (
	"github.com/nbvghost/dandelion/app/service/question"

	"github.com/nbvghost/gweb"
)

type Controller struct {
	gweb.BaseController
	Question question.QuestionService
}

func (controller *Controller) Init() {
	controller.AddHandler(gweb.ALLMethod("list_question", controller.listQuestionAction))

}
func (controller *Controller) listQuestionAction(context *gweb.Context) gweb.Result {
	//query := context.Request.URL.Query().Get("query")
	//controller.Question.ListQuestion(query)
	return &gweb.ImageBytesResult{Data: nil}
}
