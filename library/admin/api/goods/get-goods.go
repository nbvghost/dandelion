package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/tool/object"
)

type GetGoods struct {
}

func (g *GetGoods) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	contextValue := contexext.FromContext(context)
	//ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
	ID := object.ParseUint(contextValue.Query.Get("ID"))
	goodsInfo, err := service.Goods.Goods.GetGoods(db.Orm(), context, dao.PrimaryKey(ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", goodsInfo)}, err
}
