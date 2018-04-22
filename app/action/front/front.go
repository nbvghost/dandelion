package front

import (
	"fmt"

	"github.com/nbvghost/gweb"

	"net/http"
	"net/url"

	"dandelion/app/action/front/appointment"
	"dandelion/app/play"
)

type Controller struct {
	gweb.BaseController
}
type InterceptorData struct {
}

func (this InterceptorData) Execute(context *gweb.Context) bool {
	if context.Session.Attributes.Get(play.SessionUser) == nil {
		//http.SetCookie(context.Response, &http.Cookie{Name: "UID", MaxAge:-1, Path: "/"})
		//fmt.Println(context.Request.URL.Path)
		//fmt.Println(context.Request.URL.Query().Encode())
		redirect := ""
		if len(context.Request.URL.Query().Encode()) == 0 {
			redirect = context.Request.URL.Path
		} else {
			redirect = context.Request.URL.Path + "?" + context.Request.URL.Query().Encode()
		}
		//fmt.Println(url.QueryEscape(redirect))
		http.Redirect(context.Response, context.Request, "/account/open.do?redirect="+url.QueryEscape(redirect), http.StatusFound)
		return false
	} else {
		return true
	}

}
func (i *Controller) Apply() {
	fmt.Println(i)
	i.Interceptors.Add(&InterceptorData{})
	i.AddSubController("/appointment/", &appointment.Controller{})
	//i.Interceptors.Add(&InterceptorFile{})
	//i.AddHandler((&appointment.ControllerAppointment{}).Init())

}
