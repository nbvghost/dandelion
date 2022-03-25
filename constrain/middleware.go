package constrain

import "net/http"

type IMiddleware interface {
	Path(w http.ResponseWriter, r *http.Request) (bool, error)
	Cookie(w http.ResponseWriter, r *http.Request) (bool, error)
	Handle(w http.ResponseWriter, r *http.Request) (bool, error)
	Defer(w http.ResponseWriter, r *http.Request, err error)
}
type IViewRender interface {
	Render(context IContext, request *http.Request, writer http.ResponseWriter, viewData IViewResult) error
}
