package constrain

type IResult interface {
	Apply(IContext)
}

type IResultError interface {
	Apply(IContext, error)
}

type IViewResult interface {
	GetName() string
	GetResult() IResult

	//SetName(name string)
	//GetPkgPath() string
	//SetPkgPath(path string)
	//SetGlobal(v interface{})
}
