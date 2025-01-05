package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type Delete struct {
	POST struct {
		ID uint `form:"ID"`
	} `method:"POST"`
}

func (m *Delete) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Delete) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//context.Request.ParseForm()
	//fmt.Println(context.Request.FormValue("ID"))
	//ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
	//ID := object.ParseUint(context.Request.FormValue("ID"))
	err = dao.DeleteByPrimaryKey(db.Orm(), entity.Content, dao.PrimaryKey(m.POST.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}, err
}
