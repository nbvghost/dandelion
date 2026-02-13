package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type MultiGetID struct {
	GET struct {
		ID dao.PrimaryKey `uri:"ID"`
	} `method:"Get"`
}

func (m *MultiGetID) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	//var article model.Content
	article := dao.GetByPrimaryKey(db.GetDB(ctx), entity.Content, dao.PrimaryKey(m.GET.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", article)}, err
}
