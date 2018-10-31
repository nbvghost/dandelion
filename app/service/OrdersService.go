package service

import (
	"dandelion/app/play"
	"dandelion/app/service/dao"

	"dandelion/app/util"

	"errors"

	"encoding/json"
	"strings"

	"strconv"

	"time"

	"math"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type OrdersService struct {
	dao.BaseDao
	Goods           GoodsService
	ShoppingCart    ShoppingCartService
	TimeSell        TimeSellService
	Collage         CollageService
	ExpressTemplate ExpressTemplateService
	FullCut         FullCutService
	OrdersGoods     OrdersGoodsService
	User            UserService
	Wx              WxService
	Journal         JournalService
	CardItem        CardItemService
	Organization    OrganizationService
}

func (service OrdersService) Situation(StartTime, EndTime int64) interface{} {

	st := time.Unix(StartTime/1000, 0)
	st = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
	et := time.Unix(EndTime/1000, 0).Add(24 * time.Hour)
	et = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())

	Orm := dao.Orm()

	type Result struct {
		TotalMoney uint64 `gorm:"column:TotalMoney"`
		TotalCount uint64 `gorm:"column:TotalCount"`
	}

	var result Result

	Orm.Table("Orders").Select("SUM(PayMoney) as TotalMoney,COUNT(ID) as TotalCount").Where("CreatedAt>=?", st).Where("CreatedAt<?", et).Where("IsPay=?", 1).Find(&result)
	//fmt.Println(result)
	return result
}
func (service OrdersService) RefundInfo(OrdersGoodsID uint64, ShipName, ShipNo string) (error, string) {
	Orm := dao.Orm()

	var ordersGoods dao.OrdersGoods
	service.Get(Orm, OrdersGoodsID, &ordersGoods)

	var RefundInfo dao.RefundInfo
	util.JSONToStruct(ordersGoods.RefundInfo, &RefundInfo)
	RefundInfo.ShipName = ShipName
	RefundInfo.ShipNo = ShipNo

	err := service.ChangeMap(Orm, OrdersGoodsID, &dao.OrdersGoods{}, map[string]interface{}{"RefundInfo": util.StructToJSON(&RefundInfo), "Status": play.OS_OGRefundInfo})
	if err != nil {

		return err, ""
	}
	return nil, "快递信息填写成功"
}
func (service OrdersService) RefundComplete(OrdersGoodsID, RefundType uint64) (error, string) {
	tx := dao.Orm().Begin()

	var ordersGoods dao.OrdersGoods
	service.Get(tx, OrdersGoodsID, &ordersGoods)

	var orders dao.Orders
	service.Get(tx, ordersGoods.OrdersID, &orders)

	RefundPrice := int64(ordersGoods.SellPrice) - int64(math.Floor(((float64(ordersGoods.SellPrice)*float64(ordersGoods.Quantity))/float64(orders.GoodsMoney)*float64(orders.DiscountMoney))+0.5))
	if RefundPrice < 0 {
		RefundPrice = 0
	}
	var RefundInfo dao.RefundInfo
	util.JSONToStruct(ordersGoods.RefundInfo, &RefundInfo)
	RefundInfo.RefundPrice = uint64(RefundPrice)

	err := service.ChangeMap(tx, OrdersGoodsID, &dao.OrdersGoods{}, map[string]interface{}{"RefundInfo": util.StructToJSON(&RefundInfo), "Status": play.OS_OGRefundComplete})
	if err != nil {
		tx.Rollback()
		return err, ""
	}

	Success, Message := service.Wx.Refund(orders, orders.PayMoney, RefundInfo.RefundPrice, "用户申请退款", RefundType)
	if !Success {
		tx.Rollback()
		return errors.New(Message), ""
	}

	ogs, err := service.OrdersGoods.FindByOrdersID(tx, ordersGoods.OrdersID)
	if err != nil {
		tx.Rollback()
		return err, ""
	}
	haveRefunc := false
	for _, value := range ogs {
		if !strings.EqualFold(value.Status, play.OS_OGRefundComplete) && !strings.EqualFold(value.Status, "") {
			haveRefunc = true
			break
		}
	}

	if haveRefunc == false {
		//orders 所有的子单品订单，已经全部退款成功。改orders为完成
		err := service.ChangeMap(tx, orders.ID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_OrderOk})
		if err != nil {
			tx.Rollback()
			return err, ""
		}

	}

	tx.Commit()

	//err := service.ChangeMap(Orm, OrdersGoodsID, &dao.OrdersGoods{}, map[string]interface{}{"Status": play.OS_OGRefundOk})
	return nil, "已经同意,并已退款"
}
func (service OrdersService) RefundOk(OrdersGoodsID uint64) (error, string) {
	Orm := dao.Orm()
	err := service.ChangeMap(Orm, OrdersGoodsID, &dao.OrdersGoods{}, map[string]interface{}{"Status": play.OS_OGRefundOk})
	return err, "已经同意"
}
func (service OrdersService) RefundNo(OrdersGoodsID uint64) (error, string) {
	Orm := dao.Orm()
	err := service.ChangeMap(Orm, OrdersGoodsID, &dao.OrdersGoods{}, map[string]interface{}{"Status": play.OS_OGRefundNo})
	return err, "已经拒绝"
}
func (service OrdersService) AskRefund(OrdersGoodsID uint64, RefundInfo dao.RefundInfo) (error, string) {
	tx := dao.Orm().Begin()

	var ordersGoods dao.OrdersGoods
	service.Get(tx, OrdersGoodsID, &ordersGoods)

	var orders dao.Orders
	service.Get(tx, ordersGoods.OrdersID, &orders)

	if ordersGoods.ID == 0 {

		return errors.New("订单不存在"), ""
	}
	if orders.ID == 0 {

		return errors.New("订单不存在"), ""
	}
	//下单状态,如果订单状态为，已经发货状态或正在退款中
	if strings.EqualFold(orders.Status, play.OS_Deliver) || strings.EqualFold(orders.Status, play.OS_Refund) {

		err := service.ChangeMap(tx, OrdersGoodsID, &dao.OrdersGoods{}, map[string]interface{}{"RefundInfo": util.StructToJSON(&RefundInfo), "Status": play.OS_OGAskRefund})
		if err != nil {
			tx.Rollback()
			return err, ""
		} else {
			var err error
			if strings.EqualFold(orders.Status, play.OS_Deliver) {
				err = service.ChangeMap(tx, orders.ID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_Refund, "RefundTime": time.Now()})
			} else {
				err = service.ChangeMap(tx, orders.ID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_Refund})
			}

			if err != nil {
				tx.Rollback()
				return err, ""
			}
			tx.Commit()
			return nil, "已经申请，等待商家确认"
		}

	}
	return errors.New("不允许申请退款"), ""
}

/*func (service OrdersService) AddOrderBrokerageTemp(UserID uint64, OrderNo string, Amount int64) error {
	var orderbrokerage dao.OrderBrokerageTemp
	orderbrokerage.OrderNo = OrderNo
	orderbrokerage.Brokerage = uint64(Amount)
	orderbrokerage.UserID = UserID
	err := service.Add(dao.Orm(), &orderbrokerage)
	return err
}*/
func (service OrdersService) AddOrdersPackage(TotalMoney uint64, UserID uint64) (error, dao.OrdersPackage) {

	//OrderNo       string    `gorm:"column:OrderNo;unique"` //订单号
	//OrderList string `gorm:"column:OrderList;type:LONGTEXT"`//json []
	//PayMoney      uint64    `gorm:"column:PayMoney"`      //支付价
	//IsPay         uint64    `gorm:"column:IsPay"`          //是否支付成功,0=未支付，1，支付成功，2过期
	//PrepayID      string    `gorm:"column:PrepayID"`
	//UserID        uint64    `gorm:"column:UserID"`         //用户ID

	var orderbrokerage dao.OrdersPackage
	orderbrokerage.OrderNo = tool.UUID()

	//orderbrokerage.OrderList = util.StructToJSON(OrderList)

	/*var totalMoney uint64 = 0
	for _,v := range OrderList{
		totalMoney=totalMoney+v.PayMoney
	}*/
	orderbrokerage.TotalPayMoney = TotalMoney
	orderbrokerage.IsPay = 0
	orderbrokerage.UserID = UserID
	err := service.Add(dao.Orm(), &orderbrokerage)
	return err, orderbrokerage
}

//确认收货
func (service OrdersService) TakeDeliver(OrdersID uint64) (error, string) {
	Orm := dao.Orm()

	var orders dao.Orders
	service.Get(Orm, OrdersID, &orders)
	if orders.ID == 0 {

		return errors.New("订单不存在"), ""
	}
	//下单状态,只有邮寄才能确认收货
	if strings.EqualFold(orders.Status, play.OS_Deliver) && orders.PostType == 1 {

		tx := Orm.Begin()

		err := service.ChangeMap(tx, orders.ID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_OrderOk, "ReceiptTime": time.Now()})
		if err != nil {
			tx.Rollback()
			return err, ""
		}

		ogs, err := service.OrdersGoods.FindByOrdersID(tx, orders.ID)
		if err != nil {
			tx.Rollback()
			return err, ""
		}

		var Brokerage uint64
		for _, value := range ogs {
			//var specification dao.Specification
			//util.JSONToStruct(value.Specification, &specification)
			Brokerage = Brokerage + value.TotalBrokerage
		}

		/*err = service.CardItem.AddOrdersGoodsCardItem(tx, orders, ogs)
		if err != nil {
			tx.Rollback()
			return err, ""
		}*/

		//err = service.CardItem.AddOrdersGoodsCardItem()

		//Orm *gorm.DB, UserID uint64, Brokerage uint64, TargetID uint64, PayMenoy uint64
		err = service.User.SettlementUser(tx, Brokerage, orders)
		if err != nil {
			tx.Rollback()
			return err, ""
		} else {

			tx.Commit()
			go func(ogs []dao.OrdersGoods) {
				for _, value := range ogs {
					var _goods dao.Goods
					//service.Goods.Get(dao.Orm(), value.GoodsID, &_goods)
					util.JSONToStruct(value.Goods, &_goods)
					if _goods.ID != 0 {
						service.Goods.ChangeModel(dao.Orm(), _goods.ID, &dao.Goods{CountSale: _goods.CountSale + uint64(value.Quantity)})
					}
				}

			}(ogs)
			return nil, "确认收货成功"
		}

	}
	return errors.New("不允许收货"), ""
}

//检查订单状态
func (service OrdersService) AnalysisOrdersStatus(OrdersID uint64) {

	Orm := dao.Orm()

	var orders dao.Orders
	service.Get(Orm, OrdersID, &orders)
	if orders.ID == 0 {
		//return errors.New("订单不存在"), ""
		return
	}
	if strings.EqualFold(orders.Status, play.OS_Order) {

		if time.Now().Unix() >= orders.CreatedAt.Add(3*time.Hour*24).Unix() {
			//一直处于下单状态超过3天，没有付款，自动关闭订单，并加回库存
			service.ChangeMap(Orm, orders.ID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_Closed})
			//管理商品库存
			service.Goods.OrdersStockManager(orders, false)
		}

	} else if strings.EqualFold(orders.Status, play.OS_Deliver) {
		if time.Now().Unix() >= orders.DeliverTime.Add(15*time.Hour*24).Unix() {
			//等待收货时间超过15天，自动订单完成
			//service.ChangeMap(Orm, orders.ID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_OrderOk, "ReceiptTime": time.Now()})
			//管理商品库存
			//service.Goods.OrdersStockManager(orders, false)
			service.TakeDeliver(OrdersID)
		}

	} else if strings.EqualFold(orders.Status, play.OS_Cancel) {
		if time.Now().Unix() >= orders.UpdatedAt.Add(5*time.Hour*24).Unix() {
			//订单已经支付，用户申请了取消订单，超过5天，自动取消
			err, _ := service.CancelOk(OrdersID, 0)
			if err != nil {
				service.CancelOk(OrdersID, 1)
			}
		}

	}

}
func (service OrdersService) CancelOk(OrdersID, Type uint64) (error, string) {
	Orm := dao.Orm()

	var orders dao.Orders
	service.Get(Orm, OrdersID, &orders)
	if orders.ID == 0 {

		return errors.New("订单不存在"), ""
	}
	//下单状态
	if strings.EqualFold(orders.Status, play.OS_Cancel) {
		if orders.IsPay == 1 {

			//邮寄
			if orders.PostType == 1 {
				Success, Message := service.Wx.Refund(orders, orders.PayMoney, orders.PayMoney, "用户取消", Type)
				if Success {
					err := service.ChangeMap(Orm, orders.ID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_CancelOk})
					//管理商品库存
					service.Goods.OrdersStockManager(orders, false)
					service.User.MinusSettlementUserBrokerage(Orm, orders)
					return err, Message
				} else {
					return errors.New(Message), ""
				}
			}
			if orders.PostType == 2 {
				Success, Message := service.Wx.Refund(orders, orders.PayMoney, orders.PayMoney, "用户取消", Type)
				if Success {
					tx := Orm.Begin()
					err := service.ChangeMap(tx, orders.ID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_CancelOk})
					if err != nil {
						tx.Rollback()
						return err, ""
					}
					err = service.CardItem.CancelOrdersGoodsCardItem(tx, orders.UserID, orders.ID)
					if err != nil {
						tx.Rollback()
						return err, ""
					}
					tx.Commit()

					//管理商品库存
					service.Goods.OrdersStockManager(orders, false)

					return err, Message
				} else {
					return errors.New(Message), ""
				}
			}

		}

	}
	return errors.New("不允许取消订单"), ""
}

//申请取消
func (service OrdersService) Cancel(OrdersID uint64) (error, string) {
	Orm := dao.Orm()

	var orders dao.Orders
	service.Get(Orm, OrdersID, &orders)
	if orders.ID == 0 {

		return errors.New("订单不存在"), ""
	}
	//下单状态
	if strings.EqualFold(orders.Status, play.OS_Order) {
		if orders.IsPay == 1 {
			err := service.ChangeMap(Orm, OrdersID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_Cancel})
			return err, "申请取消，等待客服确认"
		} else {
			Success, _ := service.Wx.OrderQuery(orders.OrderNo)
			if Success {
				//如果查询订单已经支付，由客服确认
				err := service.ChangeMap(Orm, OrdersID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_Cancel})
				return err, "申请取消，等待客服确认"
			} else {
				//没支付的订单
				Success, Message1 := service.Wx.Refund(orders, orders.PayMoney, orders.PayMoney, "用户取消", 0)
				if Success == false {
					Success, Message1 = service.Wx.Refund(orders, orders.PayMoney, orders.PayMoney, "用户取消", 1)
				}

				if Success {
					//管理商品库存
					service.Goods.OrdersStockManager(orders, false)
					err := service.ChangeMap(Orm, OrdersID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_CancelOk})
					return err, "取消成功"
				} else {

					//管理商品库存
					service.Goods.OrdersStockManager(orders, false)
					err := service.ChangeMap(Orm, OrdersID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_CancelOk})
					return err, Message1 + "，取消成功"

					//return errors.New(Message1), ""
					/*Success, Message2 := service.Wx.Refund(orders, orders.PayMoney, orders.PayMoney, "用户取消", 1)
					if Success {

					} else {

					}*/
				}
			}

		}

	} else if strings.EqualFold(orders.Status, play.OS_Pay) {
		err := service.ChangeMap(Orm, OrdersID, &dao.Orders{}, map[string]interface{}{"Status": play.OS_Cancel})
		return err, "申请取消，等待客服确认"
	} else {
		return errors.New("不允许取消订单"), ""
	}

}

//发货
func (service OrdersService) Deliver(ShipName, ShipNo string, OrdersID uint64) error {
	Orm := dao.Orm().Begin()

	var orders dao.Orders
	service.Get(Orm, OrdersID, &orders)
	if orders.ID == 0 {
		Orm.Rollback()
		return errors.New("订单不存在")
	}
	if orders.IsPay != 1 {
		Orm.Rollback()
		return errors.New("订单没有支付")
	}

	err := service.ChangeModel(Orm, OrdersID, &dao.Orders{ShipName: ShipName, ShipNo: ShipNo, DeliverTime: time.Now(), Status: play.OS_Deliver})
	if err != nil {
		Orm.Rollback()
		return err
	}
	orders.ShipName = ShipName
	orders.ShipNo = ShipNo
	orders.DeliverTime = time.Now()
	orders.Status = play.OS_Deliver

	as := service.Wx.OrderDeliveryNotify(orders)
	if as.Success == false {

		err = errors.New(as.Message)
	}
	Orm.Commit()
	return err
}
func (service OrdersService) GetOrdersPackageByOrderNo(OrderNo string) dao.OrdersPackage {
	Orm := dao.Orm()
	var orders dao.OrdersPackage
	Orm.Where(&dao.OrdersPackage{OrderNo: OrderNo}).First(&orders)
	return orders
}
func (service OrdersService) GetOrdersByOrderNo(OrderNo string) dao.Orders {
	Orm := dao.Orm()
	var orders dao.Orders
	Orm.Where(&dao.Orders{OrderNo: OrderNo}).First(&orders)
	return orders
}
func (service OrdersService) GetOrdersByOrdersPackageNo(OrdersPackageNo string) []dao.Orders {
	Orm := dao.Orm()
	var orders []dao.Orders
	Orm.Where(&dao.Orders{OrdersPackageNo: OrdersPackageNo}).Find(&orders)
	return orders
}
func (service OrdersService) GetSupplyOrdersByOrderNo(OrderNo string) dao.SupplyOrders {
	Orm := dao.Orm()
	var orders dao.SupplyOrders
	Orm.Where(&dao.SupplyOrders{OrderNo: OrderNo}).First(&orders)
	return orders
}
func (service OrdersService) GetOrdersByID(ID uint64) dao.Orders {
	Orm := dao.Orm()
	var orders dao.Orders
	Orm.First(&orders, ID)
	return orders
}
func (service OrdersService) ListOrdersStatusCount(UserID uint64, Status []string) (TotalRecords int) {
	Orm := dao.Orm()
	var orders []dao.Orders
	db := Orm.Model(dao.Orders{})

	now := time.Now()
	ts := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	te := ts.Add(24 * time.Hour)

	db = db.Where("UpdatedAt>=? and UpdatedAt<?", ts, te)
	db = db.Where("UserID=?", UserID)

	if len(Status) > 0 {
		db = db.Where("Status = ?", Status[0])
		for index, value := range Status {
			if index != 0 {
				db = db.Or("Status = ?", value)
			}
		}
	}

	db.Find(&orders).Count(&TotalRecords)
	return
}
func (service OrdersService) ListOrders(UserID, OID uint64, PostType int, Status []string, Limit int, Offset int) (List []interface{}, TotalRecords int) {
	Orm := dao.Orm()
	var orders []dao.Orders

	db := Orm.Model(dao.Orders{})

	if UserID != 0 {
		db = db.Where("UserID=?", UserID)

	}
	if PostType != 0 {
		db = db.Where("PostType=?", PostType)
	}

	if len(Status) > 0 {
		db = db.Where("Status in (?)", Status)
		/*for index, value := range Status {
			if index != 0 {
				db = db.Or("Status = ?", value)
			}
		}*/
	}

	if OID > 0 {
		db = db.Where("OID=?", OID)
	}

	var recordsTotal = 0
	if Limit > 0 {
		db = db.Limit(Limit).Offset(Offset).Order("CreatedAt desc").Find(&orders).Offset(0).Count(&recordsTotal)
	} else {
		db = db.Order("CreatedAt desc").Find(&orders).Count(&recordsTotal)
	}

	results := make([]interface{}, 0)
	for _, value := range orders {

		pack := struct {
			Orders          dao.Orders
			User            dao.User
			OrdersGoodsList []dao.OrdersGoods
		}{}

		pack.Orders = value

		service.User.Get(Orm, value.UserID, &pack.User)

		ogs, _ := service.OrdersGoods.FindByOrdersID(Orm, value.ID)
		pack.OrdersGoodsList = ogs
		results = append(results, pack)
	}
	return results, recordsTotal
}

//func (service OrdersService) OrderNotify(result util.Map) (Success bool, Message string) {
func (service OrdersService) OrderNotify(total_fee uint64, out_trade_no, pay_time, attach string) (Success bool, Message string) {

	//Orm := dao.Orm()

	//TotalFee, _ := strconv.ParseUint(result["total_fee"], 10, 64)
	//OrderNo := result["out_trade_no"]
	//TimeEnd := result["time_end"]
	//attach := result["attach"]

	if strings.EqualFold(attach, play.OrdersType_Supply) {
		//充值的，目前只涉及到门店自主核销的时候，才需要用到充值
		orders := service.GetSupplyOrdersByOrderNo(out_trade_no)
		if orders.IsPay == 0 {
			tx := dao.Orm().Begin()
			t, _ := time.ParseInLocation("20060102150405", pay_time, time.Local)
			err := service.ChangeModel(tx, orders.ID, &dao.SupplyOrders{PayTime: t, IsPay: 1, PayMoney: total_fee})
			if err != nil {
				tx.Rollback()
				return false, err.Error()
			} else {
				if strings.EqualFold(orders.Type, play.SupplyType_Store) {
					err := service.Journal.AddStoreJournal(tx, orders.StoreID, "门店", "充值", play.StoreJournal_Type_CZ, int64(total_fee), orders.ID)
					if err != nil {
						tx.Rollback()
						return false, err.Error()
					} else {
						tx.Commit()
						return true, "已经支付成功"
					}
				} else {
					tx.Commit()
					strings.EqualFold(orders.Type, play.SupplyType_User)
					return false, "未实现的数据类型" + orders.Type
				}

			}
		} else {
			return false, "订单已经处理或过期"
		}

	} else if strings.EqualFold(attach, play.OrdersType_GoodsPackage) { //合并商品订单
		tx := dao.Orm().Begin()
		ordersPackage := service.GetOrdersPackageByOrderNo(out_trade_no)
		if ordersPackage.TotalPayMoney == total_fee {
			//var OrderNoList []string
			//util.JSONToStruct(ordersPackage.OrderList, &OrderNoList)

			err := service.ChangeModel(tx, ordersPackage.ID, &dao.OrdersPackage{IsPay: 1})
			if err != nil {
				tx.Rollback()
				return false, err.Error()
			}

			OrderList := service.GetOrdersByOrdersPackageNo(ordersPackage.OrderNo)

			for index, _ := range OrderList {
				//orders := service.GetOrdersByOrderNo(value)
				df, msg := service.ProcessingOrders(tx, OrderList[index], pay_time)
				if df == false {
					tx.Rollback()
					return df, msg
				}
			}
			tx.Commit()
			return true, "已经支付成功"
		} else {
			tx.Commit()
			return false, "金额不正确或订单不允许"
		}

	} else if strings.EqualFold(attach, play.OrdersType_Goods) { //商品订单
		//orders.PayMoney == total_fee.
		tx := dao.Orm().Begin()
		orders := service.GetOrdersByOrderNo(out_trade_no)
		if orders.PayMoney == total_fee {
			su, msg := service.ProcessingOrders(tx, orders, pay_time)
			if su == false {
				tx.Rollback()
				return su, msg
			}
			tx.Commit()
			return su, msg
		} else {
			tx.Commit()
			return false, "金额不正确或订单不允许"
		}

	} else {
		return false, "未实现的订单类型" + attach
	}

}

func (service OrdersService) ProcessingOrders(tx *gorm.DB, orders dao.Orders, pay_time string) (Success bool, Message string) {

	//orders := service.GetOrdersByOrderNo(out_trade_no)
	if orders.IsPay == 0 {
		if strings.EqualFold(orders.Status, play.OS_Order) {

			t, _ := time.ParseInLocation("20060102150405", pay_time, time.Local)
			//var TotalBrokerage uint64
			var err error
			if orders.PostType == 1 {
				//邮寄
				err = service.ChangeModel(tx, orders.ID, &dao.Orders{PayTime: t, IsPay: 1, Status: play.OS_Pay})
				if err != nil {

					return false, err.Error()
				}
				/*ogs, err := service.OrdersGoods.FindByOrdersID(tx, orders.ID)
				if err != nil {

					return false, err.Error()
				}

				for _, value := range ogs {
					//var specification dao.Specification
					//util.JSONToStruct(value.Specification, &specification)
					TotalBrokerage = TotalBrokerage + value.TotalBrokerage
				}*/

			} else {
				//线下使用
				err = service.ChangeModel(tx, orders.ID, &dao.Orders{PayTime: t, IsPay: 1, Status: play.OS_Pay})
				if err != nil {

					return false, err.Error()
				}

				/*ogs, err := service.OrdersGoods.FindByOrdersID(tx, orders.ID)
				if err != nil {

					return false, err.Error()
				}

				for _, value := range ogs {
					//var specification dao.Specification
					//util.JSONToStruct(value.Specification, &specification)
					TotalBrokerage = TotalBrokerage + value.TotalBrokerage
				}

				err = service.CardItem.AddOrdersGoodsCardItem(tx, orders, ogs)
				if err != nil {

					return false, err.Error()
				}*/
			}

			if err != nil {

				return false, err.Error()

			} else {

				err := service.User.FirstSettlementUserBrokerage(tx, orders)
				if err != nil {

					return false, err.Error()
				}

				return true, "已经支付成功"
			}
		} else {

			return false, "金额不正确或订单不允许"
		}
	} else {

		return false, "订单已经处理或过期"
	}

}

//从商品外直接购买，生成OrdersGoods，添加到 play.SessionConfirmOrders
func (service OrdersService) BuyOrders(Session *gweb.Session, UserID, GoodsID, SpecificationID uint64, Quantity uint) error {
	Orm := dao.Orm()
	var goods dao.Goods
	var specification dao.Specification
	//var expresstemplate dao.ExpressTemplate

	err := service.Goods.Get(Orm, GoodsID, &goods)
	if err != nil {
		return err
	}
	err = service.Goods.Get(Orm, SpecificationID, &specification)
	if err != nil {
		return err
	}
	if specification.GoodsID != goods.ID {
		return errors.New("产品与规格不匹配")
	}

	shoppingCart := dao.ShoppingCart{}
	shoppingCart.Quantity = Quantity
	shoppingCart.Specification = util.StructToJSON(specification)
	shoppingCart.Goods = util.StructToJSON(goods)
	shoppingCart.UserID = UserID

	ordersGoods := service.createOrdersGoods(shoppingCart)

	ogs := make([]dao.OrdersGoods, 0)
	ogs = append(ogs, ordersGoods)
	Session.Attributes.Put(play.SessionConfirmOrders, &ogs)

	return nil

}

func (service OrdersService) createOrdersGoods(shoppingCart dao.ShoppingCart) dao.OrdersGoods {
	//Orm := Orm()

	ordersGoods := dao.OrdersGoods{}
	var goods dao.Goods
	var specification dao.Specification
	//var timesell dao.TimeSell

	util.JSONToStruct(shoppingCart.Goods, &goods)
	util.JSONToStruct(shoppingCart.Specification, &specification)

	/*err := service.Goods.Get(Orm, shoppingCart.GoodsID, &goods)
	if err != nil {
		ordersGoods.AddError(err.Error())
	}

	err = service.Goods.Get(Orm, shoppingCart.SpecificationID, &specification)
	if err != nil {
		ordersGoods.AddError(err.Error())
	}*/
	if specification.GoodsID != goods.ID {
		//return errors.New("产品规格不匹配")
		ordersGoods.AddError("产品规格不匹配")
	}
	if specification.ID == 0 {
		//return errors.New("找不到规格")
		ordersGoods.AddError("找不到规格")
	}
	if specification.Stock-shoppingCart.Quantity < 0 {
		//return errors.New(specification.Label + "库存不足")
		ordersGoods.AddError("库存不足")
		//shoppingCart.Quantity = specification.Stock
	}

	ordersGoods.Specification = util.StructToJSON(specification)
	ordersGoods.Goods = util.StructToJSON(goods)
	ordersGoods.OID = goods.OID
	//ordersGoods.GoodsID = goods.ID
	//ordersGoods.SpecificationID = specification.ID
	ordersGoods.Quantity = shoppingCart.Quantity
	ordersGoods.CostPrice = specification.MarketPrice
	ordersGoods.SellPrice = specification.MarketPrice
	ordersGoods.OrdersGoodsNo = tool.UUID()
	ordersGoods.TotalBrokerage = uint64(ordersGoods.Quantity) * specification.Brokerage

	/*ogi:=dao.OrdersGoodsInfo{}
	ogi.OrdersGoods = ordersGoods
	ogi.Favoureds = make([]dao.Favoured,0)


	service.TimeSell.Get(Orm, goods.TimeSellID, &timesell)
	//抢购活动可用
	tsEnable := timesell.IsEnable()
	if tsEnable {

		//ordersGoods.TimeSellID = timesell.ID
		//ordersGoods.TimeSell = util.StructToJSON(timesell)
		//ordersGoods.SellPrice = ordersGoods.SellPrice - (ordersGoods.SellPrice * (uint64(timesell.Discount) / 100)) //抢购销售价
		ogi.Favoureds=append(ogi.Favoureds,dao.Favoured{Name:"限时抢购",TypeName:"TimeSell",TargetID:timesell.ID,DiscountPrice:ordersGoods.SellPrice * (uint64(timesell.Discount) / 100)})
	}else{
		//ordersGoods.TimeSellID = 0
		//ordersGoods.TimeSell = util.StructToJSON(dao.TimeSell{})
		//ogi.Favoureds=append(ogi.Favoureds,dao.Favoured{Name:"限时抢购",TypeName:"TimeSell",TargetID:timesell.ID,DiscountPrice:ordersGoods.SellPrice * (uint64(timesell.Discount) / 100)})

		brokerageProvisoConf:=service.Configuration.GetConfiguration(play.ConfigurationKey_BrokerageProviso)
		brokerageProvisoConfV,_:=strconv.ParseUint(brokerageProvisoConf.V,10,64)
		if user.Growth>=brokerageProvisoConfV{
			vipdiscountConf:=service.Configuration.GetConfiguration(play.ConfigurationKey_VIPDiscount)
			vipres:=make(map[string]interface{})
			vipres["VIPDiscount"],_=strconv.ParseUint(vipdiscountConf.V,10,64)
			item["VIP"] = vipres
		}

	}*/

	return ordersGoods
}
func (service OrdersService) AddOrders(orders *dao.Orders, list []dao.OrdersGoodsInfo) error {
	Orm := dao.Orm()

	tx := Orm.Begin()

	err := service.Add(tx, orders)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			//减掉商品库存
			tx.Commit()
			service.Goods.OrdersStockManager(*orders, true)
		}
	}()
	for _, value := range list {
		(&value).OrdersGoods.OrdersID = orders.ID
		(&value).OrdersGoods.Favoured = util.StructToJSON((&value).Favoured)
		err = service.Add(tx, &((&value).OrdersGoods))
		if err != nil {
			return err
		}

		var goods dao.Goods
		var specification dao.Specification
		util.JSONToStruct(value.OrdersGoods.Goods, &goods)
		util.JSONToStruct(value.OrdersGoods.Specification, &specification)
		err = service.ShoppingCart.DeleteByUserIDAndGoodsIDAndSpecificationID(orders.UserID, goods.ID, specification.ID)
		if err != nil {
			return err
		}
	}

	return nil

}
func (service OrdersService) ChangeOrdersPayMoney(PayMoney float64, OrdersID uint64) (Success bool, Message string) {
	tx := dao.Orm().Begin()

	orders := service.GetOrdersByID(OrdersID)

	if strings.EqualFold(orders.PrepayID, "") == false {

		success, message := service.Wx.CloseOrder(orders.OrderNo, orders.OID)
		if success == false {
			tx.Rollback()
			return false, message
		}
	}

	err := service.ChangeMap(tx, OrdersID, &dao.Orders{}, map[string]interface{}{"PayMoney": uint64(PayMoney * 100), "PrepayID": "", "OrderNo": tool.UUID()})
	if err != nil {
		tx.Rollback()
		return false, err.Error()
	}

	tx.Commit()

	return true, "订单金额修改成功"

}

//订单分析，
func (service OrdersService) AnalyseOrdersGoodsList(UserID uint64, addressee dao.Address, PostType int, AllList []dao.OrdersGoods) (error, []map[string]interface{}, uint64) {

	oslist := make(map[uint64][]dao.OrdersGoods)
	for index, v := range AllList {
		items := oslist[v.OID]
		if items == nil {
			oslist[v.OID] = make([]dao.OrdersGoods, 0)
		}
		oslist[v.OID] = append(oslist[v.OID], AllList[index])
	}

	out_result := make([]map[string]interface{}, 0)

	var golErr error
	var TotalPrice uint64 = 0

	for key, _ := range oslist {
		result := make(map[string]interface{})

		var org dao.Organization
		service.Organization.Get(dao.Orm(), key, &org)
		result["Organization"] = org

		Error, fullcut, oggs, FavouredPrice, FullCutAll, GoodsPrice, ExpressPrice := service.analyseOne(UserID, org.ID, addressee, PostType, oslist[key])
		if Error != nil {
			golErr = Error
		}
		result["Error"] = Error
		result["OrdersGoodsInfos"] = oggs
		result["FavouredPrice"] = FavouredPrice
		result["FullCutAll"] = FullCutAll
		result["GoodsPrice"] = GoodsPrice
		result["ExpressPrice"] = ExpressPrice
		result["FullCut"] = fullcut

		TotalPrice = TotalPrice + (GoodsPrice - FullCutAll + ExpressPrice)
		out_result = append(out_result, result)
	}

	return golErr, out_result, TotalPrice
}

//订单分析，
func (service OrdersService) analyseOne(UserID, OID uint64, addressee dao.Address, PostType int, list []dao.OrdersGoods) (Error error, fullcut dao.FullCut, oggs []dao.OrdersGoodsInfo, FavouredPrice, FullCutAll uint64, GoodsPrice uint64, ExpressPrice uint64) {
	Orm := dao.Orm()

	fullcuts := service.FullCut.FindOrderByAmountDesc(Orm, OID)

	//可以使用满减的金额
	FullCutPrice := uint64(0)
	//FavouredPrice := uint64(0)

	oggs = make([]dao.OrdersGoodsInfo, 0)

	expresstemplateMap := make(map[uint64]dao.ExpressTemplateNMW)

	for index, _ := range list {
		value := &list[index]
		//value.ID = 5445
		var goods dao.Goods
		var specification dao.Specification

		util.JSONToStruct(value.Goods, &goods)
		util.JSONToStruct(value.Specification, &specification)

		if PostType == 1 {
			//邮寄时，才判断库存
			if int64(specification.Stock-value.Quantity) < 0 {
				Error = errors.New(specification.Label + "库存不足")
				value.AddError(Error.Error())
				return
			}
		}

		value.Goods = util.StructToJSON(goods)
		value.Specification = util.StructToJSON(specification)

		Price := specification.MarketPrice * uint64(value.Quantity)

		value.CostPrice = specification.MarketPrice
		value.SellPrice = specification.MarketPrice
		//value.TotalBrokerage =

		ogs := dao.OrdersGoodsInfo{}
		ogs.Favoured = dao.Favoured{}
		//ogss

		timesell := service.TimeSell.GetTimeSellByGoodsID(goods.ID)
		//计算价格以及优惠
		if timesell.IsEnable() {

			Price = uint64(util.Rounding45(float64(Price)-(float64(Price)*(float64(timesell.Discount)/float64(100))), 2))
			GoodsPrice = GoodsPrice + Price

			Favoured := uint64(util.Rounding45(float64(value.SellPrice)*(float64(timesell.Discount)/float64(100)), 2))
			FavouredPrice = FavouredPrice + (Favoured * uint64(value.Quantity))

			ogs.Favoured = dao.Favoured{Name: "限时抢购", Target: util.StructToJSON(timesell), TypeName: "TimeSell", Discount: uint64(timesell.Discount)}

			value.SellPrice = value.SellPrice - Favoured

		} else {

			collage := service.Collage.GetCollageByGoodsID(goods.ID)
			if collage.ID != 0 && collage.TotalNum > 0 {

				Price = uint64(util.Rounding45(float64(Price)-(float64(Price)*(float64(collage.Discount)/float64(100))), 2))
				GoodsPrice = GoodsPrice + Price

				Favoured := uint64(util.Rounding45(float64(value.SellPrice)*(float64(collage.Discount)/float64(100)), 2))
				FavouredPrice = FavouredPrice + (Favoured * uint64(value.Quantity))

				ogs.Favoured = dao.Favoured{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint64(collage.Discount)}

				value.SellPrice = value.SellPrice - Favoured

				//goodsInfo.Favoureds = append(goodsInfo.Favoureds, dao.Favoured{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint64(collage.Discount)})
			} else {
				GoodsPrice = GoodsPrice + Price
				FullCutPrice = FullCutPrice + Price
			}

		}
		ogs.OrdersGoods = *value
		oggs = append(oggs, ogs)
		//ogss=append(ogss,ogs)

		//计算快递费，重量要加上数量,先计算规格的重，再计算购买的重量
		weight := (specification.Num * specification.Weight) * uint64(value.Quantity)

		if goods.ExpressTemplateID == 0 {
			Error = errors.New("找不到快递模板")
			value.AddError(Error.Error())
			return
		} else {
			//为每个订单设置三种计价方式
			if _, o := expresstemplateMap[goods.ExpressTemplateID]; o == false {
				nmw := dao.ExpressTemplateNMW{}
				nmw.N = nmw.N + int(value.Quantity)
				nmw.M = nmw.M + int(Price) //市场价X数量
				nmw.W = nmw.W + int(weight)
				expresstemplateMap[goods.ExpressTemplateID] = nmw
			} else {
				nmw := expresstemplateMap[goods.ExpressTemplateID]
				nmw.N = nmw.N + int(value.Quantity)
				nmw.M = nmw.M + int(Price) //市场价X数量
				nmw.W = nmw.W + int(weight)
				expresstemplateMap[goods.ExpressTemplateID] = nmw
			}

		}

	}
	//计算快满减
	for index, value := range fullcuts {

		if FullCutPrice >= value.Amount {
			FullCutAll = value.CutAmount
			//返回满减的值
			fullcut = fullcuts[index]
			break
		}
	}

	//计算快递费
	if PostType == 1 && addressee.IsEmpty() == false {

		for ID, value := range expresstemplateMap {
			var expresstemplate dao.ExpressTemplate
			service.ExpressTemplate.Get(Orm, ID, &expresstemplate)

			etFree := make([]dao.ExpressTemplateFreeItem, 0)
			json.Unmarshal([]byte(expresstemplate.Free), &etFree)

			var expressTemplateFreeItem *dao.ExpressTemplateFreeItem

		al:
			//从包邮列表中，找出一个计费方式
			for _, exp_f_value := range etFree {

				for _, exp_f_a_value := range exp_f_value.Areas {
					if strings.EqualFold(addressee.ProvinceName, exp_f_a_value) {
						expressTemplateFreeItem = &exp_f_value
						break al
					}
				}

			}

			if expressTemplateFreeItem != nil && expressTemplateFreeItem.IsFree(expresstemplate, value) {
				//有包邮项目
				ExpressPrice = 0

			} else {
				//无包邮项目

				etTemplate := dao.ExpressTemplateTemplate{}
				json.Unmarshal([]byte(expresstemplate.Template), &etTemplate)

				var expressTemplateItem *dao.ExpressTemplateItem

			alt:
				for _, exp_f_value := range etTemplate.Items {

					for _, exp_f_a_value := range exp_f_value.Areas {
						if strings.EqualFold(addressee.ProvinceName, exp_f_a_value) {
							expressTemplateItem = &exp_f_value
							break alt
						}
					}

				}

				if expressTemplateItem != nil {
					ExpressPrice = ExpressPrice + expressTemplateItem.CalculateExpressPrice(expresstemplate, value)
				} else {
					ExpressPrice = ExpressPrice + etTemplate.Default.CalculateExpressPrice(expresstemplate, value)
				}

			}

		}

	} else {
		ExpressPrice = 0
	}

	return

}

//从购买车提交的订单，通过 ShoppingCart ID,生成  OrdersGoods 列表,添加到 play.SessionConfirmOrders
func (service OrdersService) AddCartOrdersByShoppingCartIDs(Session *gweb.Session, UserID uint64, IDs []uint64) error {
	//Orm := Orm()
	//var scs []dao.ShoppingCart
	scs := service.ShoppingCart.GetGSIDs(UserID, IDs)
	/*err := Orm.Where(IDs).Find(&scs).Error
	if err != nil {
		return err
	}*/
	ogs := make([]dao.OrdersGoods, 0)
	for _, value := range scs {

		ordersGoods := service.createOrdersGoods(value)

		ogs = append(ogs, ordersGoods)
	}

	Session.Attributes.Put(play.SessionConfirmOrders, &ogs)

	return nil

}
func (service OrdersService) AddCartOrders(UserID uint64, GoodsID, SpecificationID uint64, Quantity uint) error {
	//Orm := dao.Orm()
	shoppingCarts := service.ShoppingCart.FindShoppingCartByUserID(UserID)

	tx := dao.Orm().Begin()

	var goods dao.Goods
	err := service.Goods.Get(tx, GoodsID, &goods)
	if err != nil {
		tx.Rollback()
		return err
	}

	var specification dao.Specification
	err = service.Goods.Get(tx, SpecificationID, &specification)
	if err != nil {
		tx.Rollback()
		return err
	}

	if specification.GoodsID != GoodsID {
		tx.Rollback()
		return errors.New("产品规格不匹配")
	}
	if specification.ID == 0 {
		tx.Rollback()
		return errors.New("找不到规格")
	}
	if specification.Stock-Quantity < 0 {
		tx.Rollback()
		return errors.New(specification.Label + "库存不足")
	}

	have := false
	for _, value := range shoppingCarts {
		var mgoods dao.Goods
		var mspecification dao.Specification
		util.JSONToStruct(value.Goods, &mgoods)
		util.JSONToStruct(value.Specification, &mspecification)

		if mgoods.ID == goods.ID && mspecification.ID == specification.ID {

			//已经存在，添加数量
			value.Quantity = value.Quantity + Quantity
			if value.Quantity > specification.Stock {
				value.Quantity = specification.Stock
			}
			err := service.ChangeModel(tx, value.ID, value)
			if err != nil {
				tx.Rollback()
				return err
			}
			have = true

		}

	}

	if have == false {

		sc := dao.ShoppingCart{}
		sc.UserID = UserID
		sc.Quantity = Quantity
		sc.Specification = util.StructToJSON(specification)
		sc.Goods = util.StructToJSON(goods)
		sc.GSID = strconv.Itoa(int(goods.ID)) + strconv.Itoa(int(specification.ID))
		err := service.Add(tx, &sc)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil

}
