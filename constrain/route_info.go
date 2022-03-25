package constrain

import "reflect"

type IRouteInfo interface {
	GetHandlerType() reflect.Type
	GetWithoutAuth() bool
}
