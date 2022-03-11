package constrain

type IInterceptor interface {
	Execute(context IContext) (broken bool, err error)
}
