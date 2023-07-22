package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/library/viewmodel"
	"github.com/nbvghost/dandelion/service/activity"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/tool"
	"github.com/pkg/errors"
)

type CreateOrders struct {
	OrdersService  order.OrdersService
	CollageService activity.CollageService
	User           *model.User `mapping:""`
	Post           struct {
		TotalAmount uint //总金额,跟订单进行匹对
		//PostType   int
		AddressID dao.PrimaryKey
		List      []viewmodel.GoodsSpecification
		//Type      string //拼团参数
		//No        string //拼团参数
	} `method:"post"`
}

func (m *CreateOrders) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	list := make([]*extends.OrdersGoods, 0)
	for _, goodsSpecification := range m.Post.List {
		goods, err := m.OrdersService.CreateOrdersGoods(ctx, m.User.ID, goodsSpecification.GoodsID, goodsSpecification.SpecificationID, goodsSpecification.Quantity)
		if err != nil {
			return nil, err
		}
		list = append(list, goods...)
	}

	/*address := model.Address{}
	err := util.JSONToStruct(m.Post.Address, &address)
	if err != nil {
		return nil, err
	}*/

	address := dao.GetByPrimaryKey(db.Orm(), &model.Address{}, m.Post.AddressID).(*model.Address)
	if address.ID == 0 {
		return nil, errors.New("the address cannot be empty")
	}

	confirmOrdersGoods, err := m.OrdersService.AnalyseOrdersGoodsList(m.User.OID, address, list)
	//如果 organizationOrders 存在着多个商家的订单，无法进入合拼支付，只能分开支付
	/*if len(organizationOrders) == 0 {
		return nil, errors.New("找不到订单")
	}
	if len(organizationOrders) > 1 {
		return nil, errors.New("多个商家的订单，无法进入合拼支付")
	}*/
	if err != nil {
		return nil, err
	}
	if m.Post.TotalAmount == confirmOrdersGoods.TotalAmount && err == nil {
		orderList := make([]model.Orders, 0)
		OutResult := make(map[string]interface{})
		//OrdersGoodsLen := float64(0)
		//OrdersGoodsNo := ""

		tx := db.Orm().Begin()

		/*op, err := m.OrdersService.AddOrdersPackage(tx, TotalPrice, m.User.ID)
		if err != nil {
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", nil)}, nil
		}*/
		//for _, value := range organizationOrders {
		//value := organizationOrders[0]
		oggs := confirmOrdersGoods.OrdersGoodsInfos
		//OrdersGoodsLen = math.Max(OrdersGoodsLen, float64(len(oggs)))

		//result["OrdersGoodsInfos"]=oggs
		FavouredPrice := confirmOrdersGoods.FavouredPrice
		FullCutAll := confirmOrdersGoods.FullCutAll
		GoodsPrice := confirmOrdersGoods.GoodsPrice
		ExpressPrice := confirmOrdersGoods.ExpressPrice
		//organization := confirmOrdersGoods.Organization

		PayMoney := GoodsPrice - FullCutAll + ExpressPrice //支付价格已经包含了 满减，限时抢购的扣去的部分  - _FullCutPrice-FavouredPrice

		orders := model.Orders{}
		orders.OrderNo = tool.UUID()
		orders.UserID = m.User.ID
		orders.OID = m.User.OID //organization.ID
		//orders.OrdersPackageNo = op.OrderNo
		//PayMoney = 100

		orders.PayMoney = PayMoney
		//orders.PostType = sqltype.OrdersPostType(m.Post.PostType)
		orders.Status = model.OrdersStatusOrder
		orders.Address = util.StructToJSON(address)
		orders.DiscountMoney = uint(FullCutAll + FavouredPrice)
		orders.GoodsMoney = uint(GoodsPrice)
		orders.ExpressMoney = uint(ExpressPrice)

		err := m.OrdersService.AddOrders(tx, &orders, oggs)
		if err != nil {
			tx.Rollback()
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", nil)}, nil
		}
		orderList = append(orderList, orders)
		//}
		tx.Commit()

		OutResult["OrderNo"] = orders.OrderNo
		OutResult["OrderCount"] = len(orderList)

		//拼团
		//todo 拼团要优化，不要入侵到订单里，通过统计去做。
		//Buy，Collage，Add
		/*if strings.EqualFold(m.Post.Type, "Collage") {
			if OrdersGoodsLen != 1 || len(orderList) != 1 {
				return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("订单数据有误，无法拼团"), "OK", nil)}, nil
			} else {
				OrderNo := OutResult["OrderNo"].(string)
				err := m.CollageService.AddCollageRecord(OrderNo, OrdersGoodsNo, m.Post.No, m.User.ID)
				return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", nil)}, nil
				//return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", nil)}
			}

		}*/

		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", OutResult)}, nil

	} else {
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: err.Error(), Data: nil}}, nil
	}
}

func (m *CreateOrders) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
