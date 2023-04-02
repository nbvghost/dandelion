package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/viewmodel"
	"github.com/nbvghost/dandelion/service/order"
)

type Buy struct {
	OrdersService order.OrdersService
	User          *model.User `mapping:""`
	Post          struct {
		List     []viewmodel.GoodsSpecification
		PostType int
		Address  model.Address
	} `method:"post"`
}

func (m *Buy) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	//context.Request.ParseForm()
	//_GSIDs := context.Request.FormValue("GSIDs")
	//Type := context.Request.FormValue("Type")
	//GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	//SpecificationID, _ := strconv.ParseUint(context.Request.FormValue("SpecificationID"), 10, 64)
	//Quantity, _ := strconv.ParseUint(context.Request.FormValue("Quantity"), 10, 64)

	//GoodsID := object.ParseUint(context.Request.FormValue("GoodsID"))
	//SpecificationID := object.ParseUint(context.Request.FormValue("SpecificationID"))
	//Quantity := object.ParseUint(context.Request.FormValue("Quantity"))

	var list []model.OrdersGoods
	for _, goodsSpecification := range m.Post.List {
		goods, err := m.OrdersService.CreateOrdersGoods(ctx, m.User.ID, goodsSpecification.GoodsID, goodsSpecification.SpecificationID, goodsSpecification.Quantity)
		if err != nil {
			return nil, err
		}
		list = append(list, goods...)
	}

	results, totalPrice, err := m.OrdersService.AnalyseOrdersGoodsList(m.User.ID, m.Post.Address, int(m.Post.PostType), list)

	return result.NewData(map[string]any{"List": results, "TotalPrice": totalPrice}), err

	/*if !strings.EqualFold(m.Post.GSIDs, "") && m.Post.GoodsID == 0 && m.Post.SpecificationID == 0 && m.Post.Quantity == 0 {
		GSIDs := strings.Split(m.Post.GSIDs, ",")
		if len(GSIDs) > 0 {
			GSIDsList := make([]string, 0)
			for _, value := range GSIDs {
				//ID, _ := strconv.ParseUint(value, 10, 64)
				ID := object.ParseUint(value)
				GSIDsList = append(GSIDsList, fmt.Sprintf("%d", ID))
			}
			err := m.OrdersService.AddCartOrdersByShoppingCartIDs(ctx, m.User.ID, GSIDsList)
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, nil
		} else {
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("没有相关ID"), "", nil)}, nil
		}
	} else {
		if m.Post.GoodsID != 0 && m.Post.SpecificationID != 0 && m.Post.Quantity != 0 {
			err := m.OrdersService.BuyOrders(ctx, m.User.ID, types.PrimaryKey(m.Post.GoodsID), types.PrimaryKey(m.Post.SpecificationID), uint(m.Post.Quantity))
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "立即购买", nil)}, nil
		} else {
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("订单数据出错"), "", nil)}, nil
		}
	}*/
}

func (m *Buy) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
