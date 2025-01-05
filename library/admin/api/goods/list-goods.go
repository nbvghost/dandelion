package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type ListGoods struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		Query    *serviceargument.ListQueryParam
		Order    dao.Sort
		PageNo   int
		PageSize int
	} `method:"post"`
}

func (g *ListGoods) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ListGoods) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//dts := &model.Datatables{}
	//util.RequestBodyToJSON(context.Request.Body, dts)

	r := service.Goods.Goods.GoodsList(g.Post.Query, g.Organization.ID, g.Post.Order.OrderByColumn(`"CreatedAt"`, true), g.Post.PageNo, g.Post.PageSize)

	//draw, recordsTotal, recordsFiltered, list := g.Goods.DatatablesListOrder(db.Orm(), g.Post.Datatables, &[]model.Goods{}, g.Organization.ID, "")
	return result.NewData(map[string]any{
		"Pagination": r,
	}), nil //&result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, err
}
