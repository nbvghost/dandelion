package constrain

import "github.com/gin-gonic/gin"

type IInterceptor interface {
	Execute(context IContext, info IRouteInfo, ginContext *gin.Context) (broken bool, err error)
}
