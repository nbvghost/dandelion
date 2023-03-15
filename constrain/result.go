package constrain

type IResult interface {
	Apply(IContext)
}

type IResultError interface {
	Apply(IContext, error)
}

type IViewResult interface {
	GetName() string
	GetResult(context IContext, viewHandler IViewHandler) IResult
	GetContentType() string
	//SetName(name string)
	//GetPkgPath() string
	//SetPkgPath(path string)
	//SetGlobal(v interface{})
}
