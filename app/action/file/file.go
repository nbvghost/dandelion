package file

import (
	"net/http"
	"net/url"

	"github.com/nbvghost/gweb"

	"dandelion/app/play"
	"dandelion/app/service"
	"dandelion/app/util"
)

type InterceptorFile struct {
}

func (this InterceptorFile) Execute(context *gweb.Context) bool {

	//util.Trace(context.Session,"context.Session")
	if context.Session.Attributes.Get(play.SessionManager) != nil || context.Session.Attributes.Get(play.SessionAdmin) != nil {

		return true
	} else {
		redirect := ""
		if len(context.Request.URL.Query().Encode()) == 0 {
			redirect = context.Request.URL.Path
		} else {
			redirect = context.Request.URL.Path + "?" + context.Request.URL.Query().Encode()
		}
		//fmt.Println(url.QueryEscape(redirect))
		http.Redirect(context.Response, context.Request, "/account/loginManagerPage?redirect="+url.QueryEscape(redirect), http.StatusFound)
		return false

	}
}

type Controller struct {
	gweb.BaseController
}

func (i *Controller) Apply() {

	//i.Interceptors.Add(&InterceptorFile{})
	i.AddHandler(gweb.ALLMethod("up", upFilePage))
	i.AddHandler(gweb.ALLMethod("load", loadFilePage))
	i.AddHandler(gweb.ALLMethod("captcha", captchaAction))

}
func captchaAction(context *gweb.Context) gweb.Result {

	buf := util.CreateCaptchaCodeBytes()

	return &gweb.ImageBytesResult{Data: buf}
}
func upFilePage(context *gweb.Context) gweb.Result {

	return service.File.UploadAction(context)
	/*if context.Session.Attributes.Get(play.SessionManager) != nil || context.Session.Attributes.Get(play.SessionAdmin) != nil {

		return service.File.UploadAction(context)
	} else {
		return &gweb.JsonResult{"没有登陆"}
	}*/
}
func loadFilePage(context *gweb.Context) gweb.Result {

	return service.File.LoadAction(context)

}
