package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/goods"
)

type GoodsTypeList struct {
	GoodsTypeService goods.GoodsTypeService
	Organization     *model.Organization `mapping:""`
}

func (m *GoodsTypeList) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*entity.Organization)
	list := m.GoodsTypeService.ListGoodsType(m.Organization.ID)
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: list}}, nil

}
