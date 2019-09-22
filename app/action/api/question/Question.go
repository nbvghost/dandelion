package question

import (
	"github.com/nbvghost/dandelion/app/service"

	"github.com/nbvghost/gweb"
)

type Controller struct {
	gweb.BaseController
	Question service.QuestionService
}

func (controller *Controller) Apply() {
	controller.AddHandler(gweb.ALLMethod("list_question", controller.listQuestionAction))

}
func (controller *Controller) listQuestionAction(context *gweb.Context) gweb.Result {
	//query := context.Request.URL.Query().Get("query")
	//controller.Question.ListQuestion(query)
	return &gweb.ImageBytesResult{Data: nil}
}
