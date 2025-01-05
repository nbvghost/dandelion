package file

import (
	"github.com/nbvghost/dandelion/constrain"
)

type Interceptor struct {
}

/*
	func (controller Interceptor) ActionService(context constrain.IContext) gweb.ServiceConfig {
		return gweb.ServiceConfig{}
	}
*/
func (controller Interceptor) ActionBefore(context constrain.IContext) (bool, constrain.IResult) {
	/*context.Response.Header().Set("Access-Control-Allow-Origin", context.Request.Header.Get("Origin"))
	context.Response.Header().Set("Access-Control-Allow-Credentials", "true")

	if context.Request.Method == http.MethodOptions {
		return false, &result.EmptyResult{}
	}*/
	return true, nil
}
func (controller Interceptor) ActionAfter(context constrain.IContext, result constrain.IResult) (r constrain.IResult, err error) {
	return nil, err
}
