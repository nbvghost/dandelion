package images

import (
	"dandelion/app/play"
	"dandelion/app/service"
	"dandelion/app/util"

	"github.com/nbvghost/gweb"
)

type Controller struct {
	gweb.BaseController
	File service.FileService
}

func (controller *Controller) Apply() {
	//controller.Interceptors.DisableManagerSession = true
	//i.Interceptors.Add(&InterceptorFile{})
	controller.AddHandler(gweb.ALLMethod("captcha", controller.captchaAction))

}
func (controller *Controller) captchaAction(context *gweb.Context) gweb.Result {
	buf := util.CreateCaptchaCodeBytes(play.SessionCaptcha)
	return &gweb.ImageBytesResult{Data: buf}
}
