package timesell

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type ID struct {
	Delete struct {
		ID uint `uri:"ID"`
	} `method:"Delete"`
}

func (m *ID) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *ID) HandleDelete(context constrain.IContext) (r constrain.IResult, err error) {
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	err = service.Activity.TimeSell.DeleteTimeSell(dao.PrimaryKey(m.Delete.ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}, err
}
