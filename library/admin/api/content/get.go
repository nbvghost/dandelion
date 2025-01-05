package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/repository"
)

type Get struct {
	POST struct {
		ContentItemID int    `form:"ContentItemID"`
		Title         string `form:"Title"`
	} `method:"POST"`
}

func (m *Get) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Get) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//context.Request.ParseForm()
	//ContentItemID := object.ParseInt(context.Request.FormValue("ContentItemID"))
	//Title := context.Request.FormValue("Title")
	content := repository.ContentDao.GetContentByContentItemIDAndTitle(uint(m.POST.ContentItemID), m.POST.Title)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", content)}, err
}
