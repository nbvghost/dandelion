package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type ListGoodsType struct {
	Organization *model.Organization `mapping:""`
	Get          struct {
		Query struct{}
		Order dao.Sort
	} `method:"Get"`
	Post struct {
		Query struct{}
		Order dao.Sort
	} `method:"Post"`
}

func (g *ListGoodsType) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//dts := &model.Datatables{}
	//util.RequestBodyToJSON(context.Request.Body, dts)
	//draw, recordsTotal, recordsFiltered, list := g.Goods.DatatablesListOrder(db.GetDB(ctx), g.Post.Datatables, &[]model.GoodsType{}, g.Organization.ID, "")

	var goodsList []model.GoodsType

	Orm := db.GetDB(ctx)

	orm := Orm.Model(model.GoodsType{})
	if g.Organization.ID > 0 {
		orm = orm.Where(`"OID"=?`, g.Organization.ID)
	}

	//var recordsTotal int64

	orm.Order(g.Post.Order.OrderByColumn(`"CreatedAt"`, true)).Find(&goodsList)

	return result.NewData(map[string]any{"List": goodsList}), nil //result.NewDa{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, err
}
func (g *ListGoodsType) Handle(context constrain.IContext) (constrain.IResult, error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//dts := &model.Datatables{}
	//util.RequestBodyToJSON(context.Request.Body, dts)
	//draw, recordsTotal, recordsFiltered, list := g.Goods.DatatablesListOrder(db.GetDB(ctx), g.Get.Datatables, &[]model.GoodsType{}, g.Organization.ID, "")
	//return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, err
	return nil, nil
}
