package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/server/redis"
	"github.com/nbvghost/dandelion/service/activity"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/dandelion/service/configuration"
	"github.com/nbvghost/dandelion/service/express"
	"github.com/nbvghost/dandelion/service/goods"
	"github.com/nbvghost/dandelion/service/journal"
	"github.com/nbvghost/dandelion/service/user"
	"github.com/nbvghost/dandelion/service/wechat"

	"gorm.io/gorm"

	"github.com/nbvghost/glog"
	"github.com/nbvghost/gpa/types"
	"github.com/nbvghost/tool"
	"github.com/nbvghost/tool/object"
)

type OrdersService struct {
	model.BaseDao
	Goods           goods.GoodsService
	ShoppingCart    ShoppingCartService
	TimeSell        activity.TimeSellService
	Collage         activity.CollageService
	Settlement      activity.SettlementService
	ExpressTemplate express.ExpressTemplateService
	FullCut         activity.FullCutService
	Wx              wechat.WxService
	Journal         journal.JournalService
	CardItem        activity.CardItemService
	Organization    company.OrganizationService
	Configuration   configuration.ConfigurationService
	User            user.UserService
}

//如果订单未完成，或是退款，扣除相应的冻结金额，不用结算，佣金
func (service OrdersService) AfterSettlementUserBrokerage(tx *gorm.DB, orders model.Orders) error {
	var err error
	//用户自己。下单者
	//Orm:=singleton.Orm()

	//var orders model.Orders
	//service.Get(Orm, OrderID, &orders)

	ogs, err := service.FindOrdersGoodsByOrdersID(tx, orders.ID)
	var Brokerage uint
	for _, value := range ogs {
		//var specification model.Specification
		//util.JSONToStruct(value.Specification, &specification)
		Brokerage = Brokerage + value.TotalBrokerage
	}

	var user model.User
	service.Get(tx, orders.UserID, &user)

	leve1 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve1).V)
	leve2 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve2).V)
	leve3 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve3).V)
	leve4 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve4).V)
	leve5 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve5).V)
	leve6 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve6).V)

	leves := []uint{leve1, leve2, leve3, leve4, leve5, leve6}

	//var OutBrokerageMoney int64 = 0
	for _, value := range leves {
		if value <= 0 {
			break
		}
		var _user model.User
		service.Get(tx, user.SuperiorID, &_user)
		if _user.ID <= 0 {
			return nil
		}
		leveMenoy := int64(math.Floor(float64(value)/float64(100)*float64(Brokerage) + 0.5))
		err = service.User.AddUserBlockAmount(tx, _user.ID, -leveMenoy)
		if err != nil {
			glog.Error(err)
			continue
		}
		//OutBrokerageMoney = OutBrokerageMoney + leveMenoy
		//workTime := time.Now().Unix() - orders.CreatedAt.Unix()

		//service.Wx.INComeNotify(_user, "来自"+strconv.Itoa(index+1)+"级用户，预计现金收入", strconv.Itoa(int(workTime/60/60))+"小时", "预计收入："+strconv.FormatFloat(float64(leveMenoy)/float64(100), 'f', 2, 64)+"元")
		//fmt.Println("预计收入：" + strconv.FormatFloat(float64(leveMenoy)/float64(100), 'f', 2, 64) + "元")
		user = _user
	}

	return err
}
func (service OrdersService) FirstSettlementUserBrokerage(tx *gorm.DB, orders model.Orders) error {
	var err error
	//用户自己。下单者
	//Orm:=singleton.Orm()

	//var orders model.Orders
	//service.Get(Orm, OrderID, &orders)

	ogs, err := service.FindOrdersGoodsByOrdersID(tx, orders.ID)
	var Brokerage uint
	for _, value := range ogs {
		//var specification model.Specification
		//util.JSONToStruct(value.Specification, &specification)
		Brokerage = Brokerage + value.TotalBrokerage
	}

	var user model.User
	service.Get(tx, orders.UserID, &user)

	leve1 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve1).V)
	leve2 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve2).V)
	leve3 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve3).V)
	leve4 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve4).V)
	leve5 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve5).V)
	leve6 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve6).V)

	leves := []uint{leve1, leve2, leve3, leve4, leve5, leve6}

	//var OutBrokerageMoney int64 = 0
	for index, value := range leves {
		if value <= 0 {
			break
		}
		var _user model.User
		service.Get(tx, user.SuperiorID, &_user)
		if _user.ID <= 0 {
			return nil
		}
		leveMenoy := int64(math.Floor(float64(value)/float64(100)*float64(Brokerage) + 0.5))
		err = service.User.AddUserBlockAmount(tx, _user.ID, leveMenoy)
		if err != nil {
			glog.Error(err)
			continue
		}
		//OutBrokerageMoney = OutBrokerageMoney + leveMenoy
		workTime := time.Now().Unix() - orders.CreatedAt.Unix()

		service.Wx.INComeNotify(_user, "来自"+strconv.Itoa(index+1)+"级用户，预计现金收入", strconv.Itoa(int(workTime/60/60))+"小时", "预计收入："+strconv.FormatFloat(float64(leveMenoy)/float64(100), 'f', 2, 64)+"元")
		//fmt.Println("预计收入：" + strconv.FormatFloat(float64(leveMenoy)/float64(100), 'f', 2, 64) + "元")
		user = _user
	}

	return err
}
func (service OrdersService) MinusSettlementUserBrokerage(tx *gorm.DB, orders model.Orders) error {
	var err error
	//用户自己。下单者
	//Orm:=singleton.Orm()

	ogs, err := service.FindOrdersGoodsByOrdersID(tx, orders.ID)
	var Brokerage uint
	for _, value := range ogs {
		//var specification model.Specification
		//util.JSONToStruct(value.Specification, &specification)
		Brokerage = Brokerage + value.TotalBrokerage
	}

	//var orders model.Orders
	//service.Get(Orm, OrderID, &orders)

	var user model.User
	service.Get(tx, orders.UserID, &user)

	leve1 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve1).V)
	leve2 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve2).V)
	leve3 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve3).V)
	leve4 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve4).V)
	leve5 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve5).V)
	leve6 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, configuration.ConfigurationKeyBrokerageLeve6).V)

	leves := []uint{leve1, leve2, leve3, leve4, leve5, leve6}

	//var OutBrokerageMoney int64 = 0
	for _, value := range leves {
		if value <= 0 {
			break
		}
		var _user model.User
		service.Get(tx, user.SuperiorID, &_user)
		if _user.ID <= 0 {
			return nil
		}
		leveMenoy := int64(math.Floor(float64(value)/float64(100)*float64(Brokerage) + 0.5))
		err = service.User.AddUserBlockAmount(tx, _user.ID, -leveMenoy)
		if err != nil {
			return err
		}
		//OutBrokerageMoney = OutBrokerageMoney + leveMenoy
		//workTime := time.Now().Unix() - orders.CreatedAt.Unix()
		//service.Wx.INComeNotify(_user, "来自"+strconv.Itoa(index+1)+"级用户，预计现金收入", strconv.Itoa(int(workTime/60/60))+"小时", "预计收入："+strconv.FormatFloat(float64(leveMenoy)/float64(100), 'f', 2, 64)+"元")
		//fmt.Println("预计收入：" + strconv.FormatFloat(float64(-leveMenoy)/float64(100), 'f', 2, 64) + "元")
		user = _user
	}

	return err
}

func (service OrdersService) OrdersStockManager(db *gorm.DB, orders model.Orders, isMinus bool) error {

	if orders.PostType == 2 {
		//线下订单，不去维护在线商品库存
		log.Println("线下订单，不去维护在线商品库存")
		return nil
	}

	//管理商品库存
	//Orm := singleton.Orm()
	//list []model.OrdersGoods

	list, _ := service.FindOrdersGoodsByOrdersID(db, orders.ID)
	for _, value := range list {
		var specification model.Specification
		//service.Get(Orm, value.SpecificationID, &specification)
		util.JSONToStruct(value.Specification, &specification)
		var goods model.Goods
		//service.Get(Orm, value.GoodsID, &goods)
		util.JSONToStruct(value.Goods, &goods)

		if isMinus {
			//减
			Stock := int64(specification.Stock - value.Quantity)
			if Stock < 0 {
				Stock = 0
			}
			err := service.ChangeMap(db, specification.ID, &model.Specification{}, map[string]interface{}{"Stock": uint(Stock)})
			if err != nil {
				return err
			}
			Stock = int64(goods.Stock - value.Quantity)
			if Stock < 0 {
				Stock = 0
			}
			err = service.ChangeMap(db, goods.ID, &model.Goods{}, map[string]interface{}{"Stock": uint(Stock)})
			if err != nil {
				return err
			}
		} else {
			//添加
			Stock := int64(specification.Stock + value.Quantity)
			if Stock < 0 {
				Stock = 0
			}
			err := service.ChangeMap(db, specification.ID, &model.Specification{}, map[string]interface{}{"Stock": uint(Stock)})
			if err != nil {
				return err
			}
			Stock = int64(goods.Stock + value.Quantity)
			if Stock < 0 {
				Stock = 0
			}
			err = service.ChangeMap(db, goods.ID, &model.Goods{}, map[string]interface{}{"Stock": uint(Stock)})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (service OrdersService) Situation(StartTime, EndTime int64) interface{} {

	st := time.Unix(StartTime/1000, 0)
	st = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
	et := time.Unix(EndTime/1000, 0).Add(24 * time.Hour)
	et = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())

	Orm := singleton.Orm()

	type Result struct {
		TotalMoney uint `gorm:"column:TotalMoney"`
		TotalCount uint `gorm:"column:TotalCount"`
	}

	var result Result

	Orm.Table("Orders").Select(`SUM("PayMoney") as TotalMoney,COUNT("ID") as TotalCount`).Where(`"CreatedAt">=?`, st).Where(`"CreatedAt"<?`, et).Where(map[string]interface{}{"IsPay": 1}).Find(&result)
	//fmt.Println(result)
	return result
}
func (service OrdersService) RefundInfo(OrdersGoodsID types.PrimaryKey, ShipName, ShipNo string) (error, string) {
	Orm := singleton.Orm()

	var ordersGoods model.OrdersGoods
	service.Get(Orm, OrdersGoodsID, &ordersGoods)

	var RefundInfo model.RefundInfo
	util.JSONToStruct(ordersGoods.RefundInfo, &RefundInfo)
	RefundInfo.ShipName = ShipName
	RefundInfo.ShipNo = ShipNo

	err := service.ChangeMap(Orm, OrdersGoodsID, &model.OrdersGoods{}, map[string]interface{}{"RefundInfo": util.StructToJSON(&RefundInfo), "Status": model.OrdersGoodsStatusOGRefundInfo})
	if err != nil {

		return err, ""
	}
	return nil, "快递信息填写成功"
}
func (service OrdersService) RefundComplete(OrdersGoodsID types.PrimaryKey, RefundType uint, wxConfig *model.WechatConfig) (error, string) {
	tx := singleton.Orm().Begin()

	var ordersGoods model.OrdersGoods
	service.Get(tx, OrdersGoodsID, &ordersGoods)

	var orders model.Orders
	service.Get(tx, ordersGoods.OrdersID, &orders)

	ordersPackage := service.GetOrdersPackageByOrderNo(orders.OrdersPackageNo)

	//RefundPrice := int64(ordersGoods.SellPrice) - int64(math.Floor(((float64(ordersGoods.SellPrice)*float64(ordersGoods.Quantity))/float64(orders.GoodsMoney)*float64(orders.DiscountMoney))+0.5))
	RefundPrice := ordersGoods.SellPrice * uint(ordersGoods.Quantity)
	if RefundPrice < 0 {
		RefundPrice = 0
	}
	var RefundInfo model.RefundInfo
	util.JSONToStruct(ordersGoods.RefundInfo, &RefundInfo)
	RefundInfo.RefundPrice = RefundPrice

	err := service.ChangeMap(tx, OrdersGoodsID, &model.OrdersGoods{}, map[string]interface{}{"RefundInfo": util.StructToJSON(&RefundInfo), "Status": model.OrdersGoodsStatusOGRefundComplete})
	if err != nil {
		tx.Rollback()
		return err, ""
	}

	Success, Message := service.Wx.Refund(orders, ordersPackage, orders.PayMoney, RefundInfo.RefundPrice, "用户申请退款", RefundType, wxConfig)
	if !Success {
		tx.Rollback()
		return errors.New(Message), ""
	}

	ogs, err := service.FindOrdersGoodsByOrdersID(tx, ordersGoods.OrdersID)
	if err != nil {
		tx.Rollback()
		return err, ""
	}
	haveRefunc := false
	//totalBrokerage := uint(0)
	for _, value := range ogs {
		//totalBrokerage = totalBrokerage + (value.TotalBrokerage * uint(value.Quantity))
		if !(value.Status == model.OrdersGoodsStatusOGRefundComplete) && !(value.Status == model.OrdersGoodsStatusOGNone) {
			haveRefunc = true
			break
		}
	}

	if haveRefunc == false {
		//orders 所有的子单品订单，已经全部退款成功。改orders为完成

		//err := service.ChangeMap(tx, orders.ID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusOrderOk})
		err := service.ChangeMap(tx, orders.ID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusRefundOk})
		if err != nil {
			tx.Rollback()
			return err, ""
		}
		//扣除佣金
		err = service.AfterSettlementUserBrokerage(tx, orders)
		if err != nil {
			tx.Rollback()
			return err, ""
		}
	}

	tx.Commit()

	//err := service.ChangeMap(Orm, OrdersGoodsID, &model.OrdersGoods{}, map[string]interface{}{"Status": model.OrdersStatusOGRefundOk})
	return nil, "已经同意,并已退款"
}
func (service OrdersService) RefundOk(OrdersGoodsID types.PrimaryKey) (error, string) {
	Orm := singleton.Orm()
	err := service.ChangeMap(Orm, OrdersGoodsID, &model.OrdersGoods{}, map[string]interface{}{"Status": model.OrdersGoodsStatusOGRefundOk})
	return err, "已经同意"
}
func (service OrdersService) RefundNo(OrdersGoodsID types.PrimaryKey) (error, string) {
	Orm := singleton.Orm()
	err := service.ChangeMap(Orm, OrdersGoodsID, &model.OrdersGoods{}, map[string]interface{}{"Status": model.OrdersGoodsStatusOGRefundNo})
	return err, "已经拒绝"
}
func (service OrdersService) AskRefund(OrdersGoodsID types.PrimaryKey, RefundInfo model.RefundInfo) (error, string) {
	tx := singleton.Orm().Begin()

	var ordersGoods model.OrdersGoods
	service.Get(tx, OrdersGoodsID, &ordersGoods)

	var orders model.Orders
	service.Get(tx, ordersGoods.OrdersID, &orders)

	if ordersGoods.ID == 0 {

		return errors.New("订单不存在"), ""
	}
	if orders.ID == 0 {

		return errors.New("订单不存在"), ""
	}
	//下单状态,如果订单状态为，已经发货状态或正在退款中
	if (orders.Status == model.OrdersStatusDeliver) || (orders.Status == model.OrdersStatusRefund) {

		err := service.ChangeMap(tx, OrdersGoodsID, &model.OrdersGoods{}, map[string]interface{}{"RefundInfo": util.StructToJSON(&RefundInfo), "Status": model.OrdersGoodsStatusOGAskRefund})
		if err != nil {
			tx.Rollback()
			return err, ""
		} else {
			var err error
			if orders.Status == model.OrdersStatusDeliver {
				err = service.ChangeMap(tx, orders.ID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusRefund, "RefundTime": time.Now()})
			} else {
				err = service.ChangeMap(tx, orders.ID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusRefund})
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

/*func (service OrdersService) AddOrderBrokerageTemp(UserID uint, OrderNo string, Amount int64) error {
	var orderbrokerage model.OrderBrokerageTemp
	orderbrokerage.OrderNo = OrderNo
	orderbrokerage.Brokerage = uint(Amount)
	orderbrokerage.UserID = UserID
	err := service.Add(singleton.Orm(), &orderbrokerage)
	return err
}*/
func (service OrdersService) AddOrdersPackage(TotalMoney uint, UserID types.PrimaryKey) (error, model.OrdersPackage) {

	//OrderNo       string    `gorm:"column:OrderNo;unique"` //订单号
	//OrderList string `gorm:"column:OrderList;type:LONGTEXT"`//json []
	//PayMoney      uint    `gorm:"column:PayMoney"`      //支付价
	//IsPay         uint    `gorm:"column:IsPay"`          //是否支付成功,0=未支付，1，支付成功，2过期
	//PrepayID      string    `gorm:"column:PrepayID"`
	//UserID        uint    `gorm:"column:UserID"`         //用户ID

	var orderbrokerage model.OrdersPackage
	orderbrokerage.OrderNo = tool.UUID()

	//orderbrokerage.OrderList = util.StructToJSON(OrderList)

	/*var totalMoney uint = 0
	for _,v := range OrderList{
		totalMoney=totalMoney+v.PayMoney
	}*/
	orderbrokerage.TotalPayMoney = TotalMoney
	orderbrokerage.IsPay = 0
	orderbrokerage.UserID = UserID
	err := service.Add(singleton.Orm(), &orderbrokerage)
	return err, orderbrokerage
}

//确认收货
func (service OrdersService) TakeDeliver(OrdersID types.PrimaryKey) error {
	Orm := singleton.Orm()

	var orders model.Orders
	service.Get(Orm, OrdersID, &orders)
	if orders.ID == 0 {

		return errors.New("订单不存在")
	}
	//下单状态,只有邮寄才能确认收货
	if (orders.Status == model.OrdersStatusDeliver) && orders.PostType == 1 {

		tx := Orm.Begin()

		err := service.ChangeMap(tx, orders.ID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusOrderOk, "ReceiptTime": time.Now()})
		if err != nil {
			tx.Rollback()
			return err
		}

		ogs, err := service.FindOrdersGoodsByOrdersID(tx, orders.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		var Brokerage uint
		for _, value := range ogs {
			//var specification model.Specification
			//util.JSONToStruct(value.Specification, &specification)
			Brokerage = Brokerage + value.TotalBrokerage
		}

		/*err = service.CardItem.AddOrdersGoodsCardItem(tx, orders, ogs)
		if err != nil {
			tx.Rollback()
			return err, ""
		}*/

		//err = service.CardItem.AddOrdersGoodsCardItem()

		//Orm *gorm.DB, UserID uint, Brokerage uint, TargetID uint, PayMenoy uint
		err = service.Settlement.SettlementUser(tx, Brokerage, orders)
		if err != nil {
			tx.Rollback()
			return err
		} else {

			tx.Commit()
			go func(ogs []model.OrdersGoods) {
				for _, value := range ogs {
					var _goods model.Goods
					//service.Goods.Get(singleton.Orm(), value.GoodsID, &_goods)
					err = util.JSONToStruct(value.Goods, &_goods)
					if err != nil {
						return
					}
					if _goods.ID != 0 {
						err = service.Goods.ChangeModel(singleton.Orm(), _goods.ID, &model.Goods{CountSale: _goods.CountSale + uint(value.Quantity)})
						if err != nil {
							return
						}
					}
				}

			}(ogs)
			return nil
		}

	}
	return errors.New("不允许收货")
}

//检查订单状态
func (service OrdersService) AnalysisOrdersStatus(OrdersID types.PrimaryKey, wxConfig *model.WechatConfig) error {

	Orm := singleton.Orm()

	var orders model.Orders
	service.Get(Orm, OrdersID, &orders)
	if orders.ID == 0 {
		//return errors.New("订单不存在"), ""
		return nil
	}
	if orders.Status == model.OrdersStatusOrder {

		if time.Now().Unix() >= orders.CreatedAt.Add(3*time.Hour*24).Unix() {
			//一直处于下单状态超过3天，没有付款，自动关闭订单，并加回库存
			err := service.ChangeMap(Orm, orders.ID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusClosed})
			if err != nil {
				return err
			}
			//管理商品库存
			err = service.OrdersStockManager(Orm, orders, false)
			if err != nil {
				return err
			}
		}

	} else if orders.Status == model.OrdersStatusDeliver {
		if time.Now().Unix() >= orders.DeliverTime.Add(15*time.Hour*24).Unix() {
			//等待收货时间超过15天，自动订单完成
			//service.ChangeMap(Orm, orders.ID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusOrderOk, "ReceiptTime": time.Now()})
			//管理商品库存
			//service.Goods.OrdersStockManager(orders, false)
			err := service.TakeDeliver(OrdersID)
			if err != nil {
				return err
			}
		}

	} else if orders.Status == model.OrdersStatusCancel {
		if time.Now().Unix() >= orders.UpdatedAt.Add(5*time.Hour*24).Unix() {
			//订单已经支付，用户申请了取消订单，超过5天，自动取消
			_, err := service.CancelOk(OrdersID, 0, wxConfig)
			if err != nil {
				_, err = service.CancelOk(OrdersID, 1, wxConfig)
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}
func (service OrdersService) CancelOk(OrdersID types.PrimaryKey, Type uint, wxConfig *model.WechatConfig) (string, error) {
	Orm := singleton.Orm()

	var orders model.Orders
	service.Get(Orm, OrdersID, &orders)
	if orders.ID == 0 {

		return "", errors.New("订单不存在")
	}

	ordersPackage := service.GetOrdersPackageByOrderNo(orders.OrdersPackageNo)

	//下单状态
	if orders.Status == model.OrdersStatusCancel {
		if orders.IsPay == sqltype.OrdersIsPayPayed {

			//邮寄
			if orders.PostType == sqltype.OrdersPostTypePost {
				Success, Message := service.Wx.Refund(orders, ordersPackage, orders.PayMoney, orders.PayMoney, "用户取消", Type, wxConfig)
				if Success {
					err := service.ChangeMap(Orm, orders.ID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusCancelOk})
					if err != nil {
						return Message, err
					}
					//管理商品库存
					err = service.OrdersStockManager(Orm, orders, false)
					if err != nil {
						return Message, err
					}
					err = service.MinusSettlementUserBrokerage(Orm, orders)
					if err != nil {
						return Message, err
					}
				} else {
					return "", errors.New(Message)
				}
			}
			if orders.PostType == sqltype.OrdersPostTypeOffline {
				Success, Message := service.Wx.Refund(orders, ordersPackage, orders.PayMoney, orders.PayMoney, "用户取消", Type, wxConfig)
				if Success {
					tx := Orm.Begin()
					err := service.ChangeMap(tx, orders.ID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusCancelOk})
					if err != nil {
						tx.Rollback()
						return "", err
					}
					ogs, err := service.FindOrdersGoodsByOrdersID(tx, orders.ID)
					if err != nil {
						tx.Rollback()
						return "", err
					}
					err = service.CardItem.CancelOrdersGoodsCardItem(tx, orders.UserID, ogs)
					if err != nil {
						tx.Rollback()
						return "", err
					}
					tx.Commit()

					//管理商品库存
					err = service.OrdersStockManager(singleton.Orm(), orders, false)
					if err != nil {
						return "", err
					}

					return Message, err
				} else {
					return "", errors.New(Message)
				}
			}

		}

	}
	return "", errors.New("不允许取消订单")
}

//申请取消
func (service OrdersService) Cancel(OrdersID types.PrimaryKey, wxConfig *model.WechatConfig) (string, error) {
	Orm := singleton.Orm()

	var orders model.Orders
	service.Get(Orm, OrdersID, &orders)
	if orders.ID == 0 {

		return "", errors.New("订单不存在")
	}

	ordersPackage := service.GetOrdersPackageByOrderNo(orders.OrdersPackageNo)

	//下单状态
	if orders.Status == model.OrdersStatusOrder {
		if orders.IsPay == 1 {
			err := service.ChangeMap(Orm, OrdersID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusCancel})
			return "申请取消，等待客服确认", err
		} else {
			Success, _ := service.Wx.OrderQuery(orders.OrderNo, wxConfig)
			if Success {
				//如果查询订单已经支付，由客服确认
				err := service.ChangeMap(Orm, OrdersID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusCancel})
				return "申请取消，等待客服确认", err
			} else {
				//没支付的订单
				Success, Message1 := service.Wx.Refund(orders, ordersPackage, orders.PayMoney, orders.PayMoney, "用户取消", 0, wxConfig)
				glog.Trace("Orders", "Cancel", Message1)
				if Success == false {
					Success, Message1 = service.Wx.Refund(orders, ordersPackage, orders.PayMoney, orders.PayMoney, "用户取消", 1, wxConfig)
					glog.Trace("Orders", "Cancel", Message1)
				}

				if Success {
					//管理商品库存
					err := service.OrdersStockManager(Orm, orders, false)
					if err != nil {
						return "", err
					}
					err = service.ChangeMap(Orm, OrdersID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusCancelOk})
					return "取消成功", err
				} else {

					//管理商品库存
					err := service.OrdersStockManager(Orm, orders, false)
					if err != nil {
						return "", err
					}
					err = service.ChangeMap(Orm, OrdersID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusCancelOk})
					return "取消成功", err

					//return errors.New(Message1), ""
					/*Success, Message2 := service.Wx.Refund(orders, orders.PayMoney, orders.PayMoney, "用户取消", 1)
					if Success {

					} else {

					}*/
				}
			}

		}

	} else if orders.Status == model.OrdersStatusPay {
		err := service.ChangeMap(Orm, OrdersID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusCancel})
		return "申请取消，等待客服确认", err
	} else {
		return "", errors.New("不允许取消订单")
	}

}

//发货
func (service OrdersService) Deliver(ShipName, ShipNo string, OrdersID types.PrimaryKey, wxConfig *model.WechatConfig) error {
	Orm := singleton.Orm().Begin()

	var orders model.Orders
	service.Get(Orm, OrdersID, &orders)
	if orders.ID == 0 {
		Orm.Rollback()
		return errors.New("订单不存在")
	}
	if orders.IsPay != 1 {
		Orm.Rollback()
		return errors.New("订单没有支付")
	}

	err := service.ChangeModel(Orm, OrdersID, &model.Orders{ShipName: ShipName, ShipNo: ShipNo, DeliverTime: time.Now(), Status: model.OrdersStatusDeliver})
	if err != nil {
		Orm.Rollback()
		return err
	}
	orders.ShipName = ShipName
	orders.ShipNo = ShipNo
	orders.DeliverTime = time.Now()
	orders.Status = model.OrdersStatusDeliver

	ogs, err := service.FindOrdersGoodsByOrdersID(singleton.Orm(), orders.ID)
	if glog.Error(err) {
		Orm.Rollback()
		return err
	}

	as := service.Wx.OrderDeliveryNotify(orders, ogs, wxConfig)
	if as.Code != result.Success {

		err = errors.New(as.Message)
	}
	Orm.Commit()
	return err
}
func (service OrdersService) GetOrdersPackageByOrderNo(OrderNo string) model.OrdersPackage {
	Orm := singleton.Orm()
	var orders model.OrdersPackage
	Orm.Where(&model.OrdersPackage{OrderNo: OrderNo}).First(&orders)
	return orders
}
func (service OrdersService) GetOrdersByOrderNo(OrderNo string) model.Orders {
	Orm := singleton.Orm()
	var orders model.Orders
	Orm.Where(&model.Orders{OrderNo: OrderNo}).First(&orders)
	return orders
}
func (service OrdersService) GetOrdersByOrdersPackageNo(OrdersPackageNo string) []model.Orders {
	Orm := singleton.Orm()
	var orders []model.Orders
	Orm.Where(&model.Orders{OrdersPackageNo: OrdersPackageNo}).Find(&orders)
	return orders
}
func (service OrdersService) GetSupplyOrdersByOrderNo(OrderNo string) model.SupplyOrders {
	Orm := singleton.Orm()
	var orders model.SupplyOrders
	Orm.Where(&model.SupplyOrders{OrderNo: OrderNo}).First(&orders)
	return orders
}
func (service OrdersService) GetOrdersByID(ID types.PrimaryKey) model.Orders {
	Orm := singleton.Orm()
	var orders model.Orders
	Orm.First(&orders, ID)
	return orders
}
func (service OrdersService) ListOrdersStatusCount(UserID types.PrimaryKey, Status []string) (TotalRecords int64) {
	Orm := singleton.Orm()
	var orders []model.Orders
	db := Orm.Model(model.Orders{})

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

type CollageRecord struct {
	OrdersID      uint      `gorm:"column:OrdersID"`
	No            string    `gorm:"column:No"`
	UserID        uint      `gorm:"column:UserID"`
	Collager      uint      `gorm:"column:Collager"`
	Favoured      string    `gorm:"column:Favoured"`
	Goods         string    `gorm:"column:Goods"`
	Specification string    `gorm:"column:Specification"`
	Status        string    `gorm:"column:Status"`
	IsPay         uint      `gorm:"column:IsPay"`
	Quantity      uint      `gorm:"column:Quantity"`
	CreatedAt     time.Time `gorm:"column:CreatedAt"`
	COUNT         uint      `gorm:"column:COUNT"`
	IsPaySUM      uint      `gorm:"column:SUM"`
	//OrdersGoods model.OrdersGoods
}

func (service OrdersService) ListCollageRecord(UserID types.PrimaryKey, Index int) []CollageRecord {
	Orm := singleton.Orm()

	db := Orm.Raw(`
SELECT
o.ID AS OrdersID,cr.No,cr.UserID,cr.Collager,og.Favoured,og.Goods,og.Specification,o.Status AS Status,o.IsPay AS IsPay,og.Quantity as Quantity,
(SELECT mcr.CreatedAt FROM CollageRecord mcr WHERE mcr.NO=cr.NO AND mcr.Collager>0) AS CreatedAt,
(SELECT COUNT(mo.IsPay) FROM CollageRecord mcr,Orders mo WHERE mcr.NO=cr.NO AND mo.OrderNo=mcr.OrderNo) AS COUNT,
(SELECT SUM(mo.IsPay) FROM CollageRecord mcr,Orders mo WHERE mcr.NO=cr.NO AND mo.OrderNo=mcr.OrderNo) AS IsPaySUM
FROM
User u,Orders o,CollageRecord cr,OrdersGoods og
WHERE
cr.UserID=? AND u.ID=cr.UserID AND o.OrderNo=cr.OrderNo AND og.OrdersGoodsNo=cr.OrdersGoodsNo
GROUP BY cr.No
`, UserID)
	//db := Orm.Raw("SELECT o.ID AS OrdersID,cr.No,cr.UserID,cr.Collager,cr.IsPay,sdf.*,og.Favoured,og.Goods,cr.CreatedAt as CreatedAt from User u,Orders o,CollageRecord cr,OrdersGoods og,(SELECT COUNT(cr.NO) AS COUNT,SUM(cr.IsPay) AS SUM FROM CollageRecord cr GROUP BY cr.NO) AS sdf WHERE cr.UserID=? AND u.ID=cr.UserID AND o.OrderNo=cr.OrderNo AND og.OrdersGoodsNo=cr.OrdersGoodsNo GROUP BY cr.No", UserID)

	packs := make([]CollageRecord, 0)
	db = db.Limit(play.Paging).Offset(play.Paging * Index).Order("CreatedAt desc")
	db.Scan(&packs)

	fmt.Println(packs)
	//var recordsTotal = 0
	if Index >= 0 {
		//db = db.Limit(play.Paging).Offset(play.Paging * Index).Order("CreatedAt desc").Offset(0).Count(&recordsTotal)
		//db = db.Limit(play.Paging).Offset(play.Paging * Index).Order("CreatedAt desc")
	} else {
		//db = db.Order("CreatedAt desc").Count(&recordsTotal)
		//db = db.Order("CreatedAt desc")
	}

	return packs
}
func (service OrdersService) ListOrdersDate(UserID, OID types.PrimaryKey, PostType int, Status []model.OrdersStatus, startDate, endDate time.Time, Limit int, Offset int) (List []interface{}, TotalRecords int64) {
	Orm := singleton.Orm()
	var orders []model.Orders

	db := Orm.Model(model.Orders{})

	if UserID != 0 {
		db = db.Where(`"UserID"=?`, UserID)
	}
	if OID > 0 {
		db = db.Where(`"OID"=?`, OID)
	}
	if PostType != 0 {
		db = db.Where(`"PostType"=?`, PostType)
	}

	if startDate.Unix() != 0 && endDate.Unix() != 0 {
		db = db.Where(`"UpdatedAt">=? and "UpdatedAt"<=?`, startDate, endDate)
	}

	if len(Status) > 0 {
		db = db.Where(`"Status" in ?`, Status)
		/*for index, value := range Status {
			if index != 0 {
				db = db.Or("Status = ?", value)
			}
		}*/
	}

	var recordsTotal int64 = 0
	if Limit > 0 {
		db = db.Count(&recordsTotal).Limit(Limit).Offset(Offset).Order(`"CreatedAt" desc`).Find(&orders)
	} else {
		db = db.Count(&recordsTotal).Order(`"CreatedAt" desc`).Find(&orders)
	}

	results := make([]interface{}, 0)
	for _, value := range orders {

		pack := struct {
			Orders          model.Orders
			User            model.User
			OrdersGoodsList []model.OrdersGoods
			CollageUsers    []model.User
		}{}

		pack.Orders = value

		service.User.Get(Orm, value.UserID, &pack.User)

		ogs, _ := service.FindOrdersGoodsByOrdersID(Orm, value.ID)
		pack.OrdersGoodsList = ogs
		//:todo 拼单
		//og := ogs[0]
		//pack.CollageUsers = service.FindOrdersGoodsByCollageUser(og.CollageNo)
		results = append(results, pack)
	}
	return results, recordsTotal
}
func (service OrdersService) ListOrders(UserID, OID types.PrimaryKey, PostType int, Status []model.OrdersStatus, Limit int, Offset int) (List []interface{}, TotalRecords int64) {

	return service.ListOrdersDate(UserID, OID, PostType, Status, time.Unix(0, 0), time.Unix(0, 0), Limit, Offset)
}

//func (service OrdersService) OrderNotify(result util.Map) (Success bool, Message string) {
func (service OrdersService) OrderNotify(total_fee uint, out_trade_no, pay_time, attach string) (Success bool, Message string) {

	//Orm := singleton.Orm()

	//TotalFee, _ := strconv.ParseUint(result["total_fee"], 10, 64)
	//OrderNo := result["out_trade_no"]
	//TimeEnd := result["time_end"]
	//attach := result["attach"]

	if strings.EqualFold(attach, play.OrdersType_Supply) {
		//充值的，目前只涉及到门店自主核销的时候，才需要用到充值
		orders := service.GetSupplyOrdersByOrderNo(out_trade_no)
		if orders.IsPay == 0 {
			tx := singleton.Orm().Begin()
			t, _ := time.ParseInLocation("20060102150405", pay_time, time.Local)
			err := service.ChangeModel(tx, orders.ID, &model.SupplyOrders{PayTime: t, IsPay: 1, PayMoney: total_fee})
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
		tx := singleton.Orm().Begin()
		ordersPackage := service.GetOrdersPackageByOrderNo(out_trade_no)
		if ordersPackage.TotalPayMoney == total_fee {
			//var OrderNoList []string
			//util.JSONToStruct(ordersPackage.OrderList, &OrderNoList)

			err := service.ChangeModel(tx, ordersPackage.ID, &model.OrdersPackage{IsPay: 1})
			if err != nil {
				tx.Rollback()
				return false, err.Error()
			}

			OrderList := service.GetOrdersByOrdersPackageNo(ordersPackage.OrderNo)

			for index := range OrderList {
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
		tx := singleton.Orm().Begin()
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

func (service OrdersService) ProcessingOrders(tx *gorm.DB, orders model.Orders, pay_time string) (Success bool, Message string) {

	//orders := service.GetOrdersByOrderNo(out_trade_no)
	if orders.IsPay == 0 {
		if orders.Status == model.OrdersStatusOrder {

			t, _ := time.ParseInLocation("20060102150405", pay_time, time.Local)
			//var TotalBrokerage uint
			var err error
			if orders.PostType == 1 {
				//邮寄
				err = service.ChangeModel(tx, orders.ID, &model.Orders{PayTime: t, IsPay: 1, Status: model.OrdersStatusPay})
				if err != nil {

					return false, err.Error()
				}
				/*ogs, err := service.OrdersGoods.FindByOrdersID(tx, orders.ID)
				if err != nil {

					return false, err.Error()
				}

				for _, value := range ogs {
					//var specification model.Specification
					//util.JSONToStruct(value.Specification, &specification)
					TotalBrokerage = TotalBrokerage + value.TotalBrokerage
				}*/

			} else {
				//线下使用
				err = service.ChangeModel(tx, orders.ID, &model.Orders{PayTime: t, IsPay: 1, Status: model.OrdersStatusPay})
				if err != nil {

					return false, err.Error()
				}

				/*ogs, err := service.OrdersGoods.FindByOrdersID(tx, orders.ID)
				if err != nil {

					return false, err.Error()
				}

				for _, value := range ogs {
					//var specification model.Specification
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

				err := service.FirstSettlementUserBrokerage(tx, orders)
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

// BuyCollageOrders 拼单购买
func (service OrdersService) BuyCollageOrders(ctx constrain.IContext, UserID, GoodsID, SpecificationID types.PrimaryKey, Quantity uint) error {
	Orm := singleton.Orm()
	var goods model.Goods
	var specification model.Specification
	//var expresstemplate model.ExpressTemplate

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

	shoppingCart := model.ShoppingCart{}
	shoppingCart.Quantity = Quantity
	shoppingCart.Specification = util.StructToJSON(specification)
	shoppingCart.Goods = util.StructToJSON(goods)
	shoppingCart.UserID = UserID

	ordersGoods := service.createOrdersGoods(shoppingCart)

	//ordersGoods.CollageNo = tool.UUID()
	collage := service.Collage.GetCollageByGoodsID(goods.ID, goods.OID)
	if collage.ID != 0 && collage.TotalNum > 0 {

		favoured := extends.Discount{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint(collage.Discount)}
		ordersGoods.Discounts = util.StructToJSON([]extends.Discount{favoured})
	}

	ogs := make([]model.OrdersGoods, 0)
	ogs = append(ogs, ordersGoods)
	//Session.Attributes.Put(gweb.AttributesKey(string(play.SessionConfirmOrders)), &ogs)
	ctx.Redis().Set(ctx, redis.NewConfirmOrders(ctx.UID()), &ogs, 24*time.Hour)
	return nil

}

//从商品外直接购买，生成OrdersGoods，添加到 play.SessionConfirmOrders
func (service OrdersService) BuyOrders(ctx constrain.IContext, UserID, GoodsID, SpecificationID types.PrimaryKey, Quantity uint) error {
	Orm := singleton.Orm()
	var goods model.Goods
	var specification model.Specification
	//var expresstemplate model.ExpressTemplate

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

	shoppingCart := model.ShoppingCart{}
	shoppingCart.Quantity = Quantity
	shoppingCart.Specification = util.StructToJSON(specification)
	shoppingCart.Goods = util.StructToJSON(goods)
	shoppingCart.UserID = UserID

	ordersGoods := service.createOrdersGoods(shoppingCart)

	ogs := make([]model.OrdersGoods, 0)
	ogs = append(ogs, ordersGoods)
	//Session.Attributes.Put(gweb.AttributesKey(string(play.SessionConfirmOrders)), &ogs)
	ctx.Redis().Set(ctx, redis.NewConfirmOrders(ctx.UID()), &ogs, 24*time.Hour)
	return nil

}

//从购买车提交的订单，通过 ShoppingCart ID,生成  OrdersGoods 列表,添加到 play.SessionConfirmOrders
func (service OrdersService) AddCartOrdersByShoppingCartIDs(ctx constrain.IContext, UserID types.PrimaryKey, IDs []string) error {
	//Orm := Orm()
	//var scs []model.ShoppingCart
	scs := service.ShoppingCart.GetGSIDs(UserID, IDs)
	/*err := Orm.Where(IDs).Find(&scs).Error
	if err != nil {
		return err
	}*/
	ogs := make([]model.OrdersGoods, 0)
	for _, value := range scs {

		ordersGoods := service.createOrdersGoods(value)

		ogs = append(ogs, ordersGoods)
	}

	ctx.Redis().Set(ctx, redis.NewConfirmOrders(ctx.UID()), &ogs, 24*time.Hour)
	//Session.Attributes.Put(gweb.AttributesKey(string(play.SessionConfirmOrders)), &ogs)

	return nil

}
func (service OrdersService) createOrdersGoods(shoppingCart model.ShoppingCart) model.OrdersGoods {
	//Orm := Orm()

	ordersGoods := model.OrdersGoods{}
	var goods model.Goods
	var specification model.Specification
	//var timesell model.TimeSell

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
	ordersGoods.Discounts = util.StructToJSON(service.Goods.GetDiscounts(goods.ID, goods.OID))

	/*//限时抢购
	timesell := service.TimeSell.GetTimeSellByGoodsID(goods.ID)
	if timesell.IsEnable() {
		favoured := model.Favoured{Name: "限时抢购", Target: util.StructToJSON(timesell), TypeName: "TimeSell", Discount: uint(timesell.Discount)}
		ordersGoods.Favoured = util.StructToJSON(favoured)
	}*/

	/*if strings.EqualFold("Collage", Type) {
		//拼单
		ordersGoods.CollageNo = tool.UUID()
		collage := service.Collage.GetCollageByGoodsID(goods.ID)
		if collage.ID != 0 && collage.TotalNum > 0 {

			favoured := model.Favoured{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint(collage.Discount)}
			ordersGoods.Favoured = util.StructToJSON(favoured)
		}
	} else {

	}*/
	ordersGoods.TotalBrokerage = uint(ordersGoods.Quantity) * specification.Brokerage

	return ordersGoods
}
func (service OrdersService) AddOrders(orders *model.Orders, list []extends.OrdersGoodsInfo) error {
	Orm := singleton.Orm()

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
		}
	}()
	for _, value := range list {
		(&value).OrdersGoods.OrdersID = orders.ID
		(&value).OrdersGoods.Discounts = util.StructToJSON((&value).Discounts)
		err = service.Add(tx, &((&value).OrdersGoods))
		if err != nil {
			return err
		}

		var goods model.Goods
		var specification model.Specification
		err = util.JSONToStruct(value.OrdersGoods.Goods, &goods)
		if err != nil {
			return err
		}
		err = util.JSONToStruct(value.OrdersGoods.Specification, &specification)
		if err != nil {
			return err
		}
		err = service.ShoppingCart.DeleteByUserIDAndGoodsIDAndSpecificationID(tx, orders.UserID, goods.ID, specification.ID)
		if err != nil {
			return err
		}
	}

	err = service.OrdersStockManager(tx, *orders, true)
	if err != nil {
		return err
	}

	return nil

}
func (service OrdersService) ChangeOrdersPayMoney(PayMoney float64, OrdersID types.PrimaryKey, wxConfig *model.WechatConfig) (Success result.ActionResultCode, Message string) {
	tx := singleton.Orm().Begin()

	orders := service.GetOrdersByID(OrdersID)

	if strings.EqualFold(orders.PrepayID, "") == false {

		success, message := service.Wx.CloseOrder(orders.OrderNo, orders.OID, wxConfig)
		if success == false {
			tx.Rollback()
			return result.Fail, message
		}
	}

	err := service.ChangeMap(tx, OrdersID, &model.Orders{}, map[string]interface{}{"PayMoney": uint(PayMoney * 100), "PrepayID": "", "OrderNo": tool.UUID()})
	if err != nil {
		tx.Rollback()
		return result.Fail, err.Error()
	}

	tx.Commit()

	return result.Success, "订单金额修改成功"

}

type AnalyseOrdersGoods struct {
	Organization     model.Organization
	Error            error
	OrdersGoodsInfos []extends.OrdersGoodsInfo
	FavouredPrice    uint
	FullCutAll       uint
	GoodsPrice       uint
	ExpressPrice     uint
	FullCut          model.FullCut
}

//订单分析，
func (service OrdersService) AnalyseOrdersGoodsList(UserID types.PrimaryKey, addressee extends.Address, PostType int, AllList []model.OrdersGoods) (error, []AnalyseOrdersGoods, uint) {

	oslist := make(map[types.PrimaryKey][]model.OrdersGoods)
	for index, v := range AllList {
		items := oslist[v.OID]
		if items == nil {
			oslist[v.OID] = make([]model.OrdersGoods, 0)
		}
		oslist[v.OID] = append(oslist[v.OID], AllList[index])
	}

	out_result := make([]AnalyseOrdersGoods, 0)

	var golErr error
	var TotalPrice uint = 0

	for key := range oslist {
		result := AnalyseOrdersGoods{}

		var org model.Organization
		service.Organization.Get(singleton.Orm(), key, &org)

		Error, fullcut, oggs, FavouredPrice, FullCutAll, GoodsPrice, ExpressPrice := service.analyseOne(UserID, org.ID, addressee, PostType, oslist[key])
		if Error != nil {
			golErr = Error
		}
		result.Error = Error
		result.Organization = org
		result.OrdersGoodsInfos = oggs
		result.FavouredPrice = FavouredPrice
		result.FullCutAll = FullCutAll
		result.GoodsPrice = GoodsPrice
		result.ExpressPrice = ExpressPrice
		result.FullCut = fullcut

		TotalPrice = TotalPrice + (GoodsPrice - FullCutAll + ExpressPrice)
		out_result = append(out_result, result)
	}

	return golErr, out_result, TotalPrice
}

//订单分析，
func (service OrdersService) analyseOne(UserID, OID types.PrimaryKey, addressee extends.Address, PostType int, list []model.OrdersGoods) (Error error, fullcut model.FullCut, oggs []extends.OrdersGoodsInfo, FavouredPrice, FullCutAll uint, GoodsPrice uint, ExpressPrice uint) {
	Orm := singleton.Orm()

	fullcuts := service.FullCut.FindOrderByAmountDesc(Orm, OID)

	//可以使用满减的金额
	FullCutPrice := uint(0)
	//FavouredPrice := uint(0)

	oggs = make([]extends.OrdersGoodsInfo, 0)

	expresstemplateMap := make(map[types.PrimaryKey]model.ExpressTemplateNMW)

	for index := range list {
		value := &list[index]
		//value.ID = 5445
		var goods model.Goods
		var specification model.Specification

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

		//value.Goods = util.StructToJSON(goods)
		//value.Specification = util.StructToJSON(specification)

		Price := specification.MarketPrice * uint(value.Quantity)

		value.CostPrice = specification.MarketPrice
		value.SellPrice = specification.MarketPrice
		//value.TotalBrokerage =

		ogs := extends.OrdersGoodsInfo{}
		ogs.Discounts = make([]extends.Discount, 0)
		//ogss

		var discounts []extends.Discount
		if strings.EqualFold(value.Discounts, "") == false {
			util.JSONToStruct(value.Discounts, &discounts)
		}
		//计算价格以及优惠
		if len(discounts) > 0 {
			for index := range discounts {
				favoured := discounts[index]
				Price = uint(util.Rounding45(float64(Price)-(float64(Price)*(float64(favoured.Discount)/float64(100))), 2))
				GoodsPrice = GoodsPrice + Price
				Favoured := uint(util.Rounding45(float64(value.SellPrice)*(float64(favoured.Discount)/float64(100)), 2))
				FavouredPrice = FavouredPrice + (Favoured * uint(value.Quantity))
				value.SellPrice = value.SellPrice - Favoured
			}
			ogs.Discounts = discounts

		} else {
			GoodsPrice = GoodsPrice + Price
			FullCutPrice = FullCutPrice + Price
		}

		/*timesell := service.TimeSell.GetTimeSellByGoodsID(goods.ID)
		if timesell.IsEnable() {

			Price = uint(util.Rounding45(float64(Price)-(float64(Price)*(float64(timesell.Discount)/float64(100))), 2))
			GoodsPrice = GoodsPrice + Price

			Favoured := uint(util.Rounding45(float64(value.SellPrice)*(float64(timesell.Discount)/float64(100)), 2))
			FavouredPrice = FavouredPrice + (Favoured * uint(value.Quantity))

			ogs.Favoured = model.Favoured{Name: "限时抢购", Target: util.StructToJSON(timesell), TypeName: "TimeSell", Discount: uint(timesell.Discount)}

			value.SellPrice = value.SellPrice - Favoured

		} else {

			collage := service.Collage.GetCollageByGoodsID(goods.ID)
			if collage.ID != 0 && collage.TotalNum > 0 {

				Price = uint(util.Rounding45(float64(Price)-(float64(Price)*(float64(collage.Discount)/float64(100))), 2))
				GoodsPrice = GoodsPrice + Price

				Favoured := uint(util.Rounding45(float64(value.SellPrice)*(float64(collage.Discount)/float64(100)), 2))
				FavouredPrice = FavouredPrice + (Favoured * uint(value.Quantity))

				ogs.Favoured = model.Favoured{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint(collage.Discount)}

				value.SellPrice = value.SellPrice - Favoured

				//goodsInfo.Favoureds = append(goodsInfo.Favoureds, model.Favoured{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint(collage.Discount)})
			} else {
				GoodsPrice = GoodsPrice + Price
				FullCutPrice = FullCutPrice + Price
			}

		}*/
		ogs.OrdersGoods = *value
		oggs = append(oggs, ogs)
		//ogss=append(ogss,ogs)

		//计算快递费，重量要加上数量,先计算规格的重，再计算购买的重量
		weight := (specification.Num * specification.Weight) * uint(value.Quantity)

		if goods.ExpressTemplateID == 0 {
			Error = errors.New("找不到快递模板")
			value.AddError(Error.Error())
			return
		} else {
			//为每个订单设置三种计价方式
			if _, o := expresstemplateMap[goods.ExpressTemplateID]; o == false {
				nmw := model.ExpressTemplateNMW{}
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
			var expresstemplate model.ExpressTemplate
			service.ExpressTemplate.Get(Orm, ID, &expresstemplate)

			etFree := make([]model.ExpressTemplateFreeItem, 0)
			json.Unmarshal([]byte(expresstemplate.Free), &etFree)

			var expressTemplateFreeItem *model.ExpressTemplateFreeItem

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

				etTemplate := model.ExpressTemplateTemplate{}
				json.Unmarshal([]byte(expresstemplate.Template), &etTemplate)

				var expressTemplateItem *model.ExpressTemplateItem

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

func (service OrdersService) AddCartOrders(UserID types.PrimaryKey, GoodsID, SpecificationID types.PrimaryKey, Quantity uint) error {
	//Orm := singleton.Orm()
	shoppingCarts := service.ShoppingCart.FindShoppingCartByUserID(UserID)

	tx := singleton.Orm().Begin()

	var goods model.Goods
	err := service.Goods.Get(tx, GoodsID, &goods)
	if err != nil {
		tx.Rollback()
		return err
	}

	var specification model.Specification
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
		var mgoods model.Goods
		var mspecification model.Specification
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

		sc := model.ShoppingCart{}
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
func (service OrdersService) FindOrdersGoodsByOrdersID(DB *gorm.DB, OrdersID types.PrimaryKey) ([]model.OrdersGoods, error) {
	var ogs []model.OrdersGoods
	err := service.FindWhere(DB, &ogs, &model.OrdersGoods{OrdersID: OrdersID})
	return ogs, err
}
func (service OrdersService) FindOrdersGoodsByCollageUser(CollageNo string) []model.User {
	orm := singleton.Orm()
	var user []model.User

	orm.Raw("SELECT u.* FROM Orders o,OrdersGoods og,USER u WHERE og.CollageNo=? AND o.IsPay=1 and o.ID=og.OrdersID AND u.ID=o.UserID", CollageNo).Scan(&user)
	//orm.Exec("SELECT u.* FROM Orders o,OrdersGoods og,USER u WHERE og.CollageNo=? AND o.ID=og.OrdersID AND u.ID=o.UserID", CollageNo).Find(&user)
	return user
}
