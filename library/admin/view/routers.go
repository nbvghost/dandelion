package view

import (
	"github.com/nbvghost/dandelion/constrain"
)

func Register(route constrain.IRoute) {

	//adminController := &admin.Controller{}
	//adminController.Interceptors.Set(&admin.Interceptor{})
	//adminController := gweb.NewController("admin", "")
	//adminController.NewController("template").DefaultHandle(&admin.Index{})
	//adminController.AddInterceptor(&admin.Interceptor{})

	route.RegisterView("index", &Index{})
	route.RegisterView("", &Index{})
	route.RegisterView("favicon.ico", &Index{})

}
