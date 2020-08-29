package api

import (
	"github.com/nbvghost/dandelion/app/action/api/question"

	"github.com/nbvghost/gweb"
)

type Interceptor struct {
}

func (interceptor Interceptor) Execute(context *gweb.Context) (bool, gweb.Result) {

	return true, nil

}

type Controller struct {
	gweb.BaseController
}

func (controller *Controller) Init() {

	controller.AddSubController("/question/", &question.Controller{})
}
