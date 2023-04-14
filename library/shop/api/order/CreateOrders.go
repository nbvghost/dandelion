package order

import (
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/library/viewmodel"
	"math"
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/library/util"
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
		TotalPrice uint
		PostType   int
		Address    string
		Type       string
		No         string
		List       []viewmodel.GoodsSpecification
	} `method:"post"`
}

func (m *CreateOrders) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)

	/*confirmOrdersJson, err := ctx.Redis().Get(ctx, redis.NewConfirmOrders(ctx.UID()))
	if err != nil {
		return nil, err
	}
	var ogs []model.OrdersGoods
	err = json.Unmarshal([]byte(confirmOrdersJson), &ogs)
	if err != nil {
		return nil, err
	}*/
	//ogs := context.Session.Attributes.Get(play.SessionConfirmOrders).(*[]entity.OrdersGoods)
	//context.Request.ParseForm()

	//_TotalPrice, _ := strconv.ParseUint(context.Request.FormValue("TotalPrice"), 10, 64)
	//_TotalPrice := object.ParseUint(context.Request.FormValue("TotalPrice"))
	//PostType, _ := strconv.ParseInt(context.Request.FormValue("PostType"), 10, 64)
	//AddressTxt := context.Request.FormValue("Address")

	//Type := context.Request.FormValue("Type") //Buy，Collage，Add
	//No := context.Request.FormValue("No")
	//fmt.Println(Type, No)

	list := make([]*extends.OrdersGoods, 0)
	for _, goodsSpecification := range m.Post.List {
		goods, err := m.OrdersService.CreateOrdersGoods(ctx, m.User.ID, goodsSpecification.GoodsID, goodsSpecification.SpecificationID, goodsSpecification.Quantity)
		if err != nil {
			return nil, err
		}
		list = append(list, goods...)
	}

	address := model.Address{}
	err := util.JSONToStruct(m.Post.Address, &address)
	if err != nil {
		return nil, err
	}

	organizationOrders, TotalPrice, Error := m.OrdersService.AnalyseOrdersGoodsList(m.User.ID, &address, int(m.Post.PostType), list)
	//如果 organizationOrders 存在着多个商家的订单，无法进入合拼支付，只能分开支付
	if len(organizationOrders) == 0 {
		return nil, errors.New("找不到订单")
	}
	if len(organizationOrders) > 1 {
		return nil, errors.New("多个商家的订单，无法进入合拼支付")
	}
	if m.Post.TotalPrice == TotalPrice && Error == nil {
		orderList := make([]model.Orders, 0)
		OutResult := make(map[string]interface{})
		OrdersGoodsLen := float64(0)
		OrdersGoodsNo := ""

		tx := singleton.Orm().Begin()

		/*op, err := m.OrdersService.AddOrdersPackage(tx, TotalPrice, m.User.ID)
		if err != nil {
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", nil)}, nil
		}*/
		//for _, value := range organizationOrders {
		value := organizationOrders[0]
		oggs := value.OrdersGoodsInfos
		OrdersGoodsLen = math.Max(OrdersGoodsLen, float64(len(oggs)))

		//result["OrdersGoodsInfos"]=oggs
		FavouredPrice := value.FavouredPrice
		FullCutAll := value.FullCutAll
		GoodsPrice := value.GoodsPrice
		ExpressPrice := value.ExpressPrice
		organization := value.Organization

		PayMoney := GoodsPrice - FullCutAll + ExpressPrice //支付价格已经包含了 满减，限时抢购的扣去的部分  - _FullCutPrice-FavouredPrice

		orders := model.Orders{}
		orders.OrderNo = tool.UUID()
		orders.UserID = m.User.ID
		orders.OID = organization.ID
		//orders.OrdersPackageNo = op.OrderNo
		//PayMoney = 100

		orders.PayMoney = PayMoney
		orders.PostType = sqltype.OrdersPostType(m.Post.PostType)
		orders.Status = model.OrdersStatusOrder
		orders.Address = util.StructToJSON(address)
		orders.DiscountMoney = uint(FullCutAll + FavouredPrice)
		orders.GoodsMoney = uint(GoodsPrice)
		orders.ExpressMoney = uint(ExpressPrice)

		err = m.OrdersService.AddOrders(tx, &orders, oggs)
		if err != nil {
			tx.Rollback()
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", nil)}, nil
		}
		orderList = append(orderList, orders)
		//}
		tx.Commit()

		OutResult["OrderNo"] = orders.OrderNo
		OutResult["OrderCount"] = len(orderList)

		//controller.Orders.AddOrdersPackage(tool.UUID(),)
		/*orderList := make([]model.Orders, 0)
		OutResult := make(map[string]interface{})
		OrdersGoodsLen := float64(0)
		OrdersGoodsNo := ""
		if len(results) > 1 {

			tx := singleton.Orm().Begin()
			op, err := m.OrdersService.AddOrdersPackage(tx, TotalPrice, m.User.ID)
			if err != nil {
				return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", nil)}, nil
			}
			for _, value := range results {

				oggs := value.OrdersGoodsInfos
				OrdersGoodsLen = math.Max(OrdersGoodsLen, float64(len(oggs)))

				//result["OrdersGoodsInfos"]=oggs
				FavouredPrice := value.FavouredPrice
				FullCutAll := value.FullCutAll
				GoodsPrice := value.GoodsPrice
				ExpressPrice := value.ExpressPrice
				organization := value.Organization

				PayMoney := GoodsPrice - FullCutAll + ExpressPrice //支付价格已经包含了 满减，限时抢购的扣去的部分  - _FullCutPrice-FavouredPrice

				orders := model.Orders{}
				orders.OrderNo = tool.UUID()
				orders.UserID = m.User.ID
				orders.OID = organization.ID
				orders.OrdersPackageNo = op.OrderNo
				//PayMoney = 100

				orders.PayMoney = PayMoney
				orders.PostType = sqltype.OrdersPostType(m.Post.PostType)
				orders.Status = model.OrdersStatusOrder
				orders.Address = util.StructToJSON(address)
				orders.DiscountMoney = uint(FullCutAll + FavouredPrice)
				orders.GoodsMoney = uint(GoodsPrice)
				orders.ExpressMoney = uint(ExpressPrice)

				err = m.OrdersService.AddOrders(tx, &orders, oggs)
				if err != nil {
					tx.Rollback()
					return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", nil)}, nil
				}
				orderList = append(orderList, orders)
			}
			tx.Commit()

			OutResult["OrderNo"] = op.OrderNo
			OutResult["OrderCount"] = len(orderList)

		} else {

			for _, value := range results {
				oggs := value.OrdersGoodsInfos
				OrdersGoodsNo = oggs[0].OrdersGoods.OrdersGoodsNo
				OrdersGoodsLen = math.Max(OrdersGoodsLen, float64(len(oggs)))
				//result["OrdersGoodsInfos"]=oggs
				FavouredPrice := value.FavouredPrice
				FullCutAll := value.FullCutAll
				GoodsPrice := value.GoodsPrice
				ExpressPrice := value.ExpressPrice
				organization := value.Organization

				PayMoney := GoodsPrice - FullCutAll + ExpressPrice //支付价格已经包含了 满减，限时抢购的扣去的部分  - _FullCutPrice-FavouredPrice

				orders := model.Orders{}
				orders.OrderNo = tool.UUID()
				orders.UserID = m.User.ID
				orders.OID = organization.ID
				//orders.OrdersPackageNo = op.OrderNo
				//PayMoney = 100

				orders.PayMoney = PayMoney
				orders.PostType = sqltype.OrdersPostType(m.Post.PostType)
				orders.Status = model.OrdersStatusOrder
				orders.Address = util.StructToJSON(address)
				orders.DiscountMoney = uint(FullCutAll + FavouredPrice)
				orders.GoodsMoney = uint(GoodsPrice)
				orders.ExpressMoney = uint(ExpressPrice)

				err = m.OrdersService.AddOrders(&orders, oggs)
				if err != nil {
					return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", nil)}, nil
				}
				orderList = append(orderList, orders)
			}

			OutResult["OrderNo"] = orderList[0].OrderNo
			OutResult["OrderCount"] = len(orderList)

		}*/
		//拼团
		//todo 拼团要优化，不要入侵到订单里，通过统计去做。
		//Buy，Collage，Add
		if strings.EqualFold(m.Post.Type, "Collage") {
			if OrdersGoodsLen != 1 || len(orderList) != 1 {
				return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("订单数据有误，无法拼团"), "OK", nil)}, nil
			} else {
				OrderNo := OutResult["OrderNo"].(string)
				err := m.CollageService.AddCollageRecord(OrderNo, OrdersGoodsNo, m.Post.No, m.User.ID)
				return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", nil)}, nil
				//return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", nil)}
			}

		}

		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", OutResult)}, nil

	} else {
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: Error.Error(), Data: nil}}, nil
	}
}

func (m *CreateOrders) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
