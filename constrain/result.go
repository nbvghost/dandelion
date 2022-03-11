package constrain

type IResult interface {
	Apply(IContext)
}
type IViewResult interface {
	GetName() string
	SetName(name string)
	GetPkgPath() string
	SetPkgPath(path string)
	SetGlobal(v interface{})
}
