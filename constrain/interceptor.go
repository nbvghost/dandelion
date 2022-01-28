package constrain

import (
	"net/http"
)

type IInterceptor interface {
	Execute(context IContext, info IRouteInfo, writer http.ResponseWriter, request *http.Request) (broken bool, err error)
}
