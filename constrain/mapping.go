package constrain

type IMapping interface {
	Call(context IContext) (instance interface{})
	Name() string
	Instance() interface{}
}
