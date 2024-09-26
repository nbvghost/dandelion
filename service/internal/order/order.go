package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/service/internal/payment"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service/internal/activity"
	"github.com/nbvghost/dandelion/service/internal/company"
	"github.com/nbvghost/dandelion/service/internal/configuration"
	"github.com/nbvghost/dandelion/service/internal/express"
	"github.com/nbvghost/dandelion/service/internal/goods"
	"github.com/nbvghost/dandelion/service/internal/journal"
	"github.com/nbvghost/dandelion/service/internal/user"
	"github.com/nbvghost/dandelion/service/internal/wechat"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"gorm.io/gorm/clause"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/server/redis"
	"gorm.io/gorm"

	"github.com/nbvghost/tool"
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
	MessageNotify   wechat.MessageNotify
	Journal         journal.JournalService
	CardItem        activity.CardItemService
	Organization    company.OrganizationService
	Configuration   configuration.ConfigurationService
	User            user.UserService
}

type ShoppingCartResult struct {
	ConfirmOrdersGoods *extends.ConfirmOrdersGoods
	ShoppingCartList   []*model.ShoppingCart
}

func (m OrdersService) FindShoppingCartListDetails(oid dao.PrimaryKey, userID dao.PrimaryKey, address *model.Address) (*ShoppingCartResult, error) {
	list := m.ShoppingCart.FindShoppingCartByUserID(userID)

	shoppingCartList := make([]*model.ShoppingCart, 0)
	orderGoodsList := make([]*model.OrdersGoods, 0)
	for i := range list {
		item := list[i].(*model.ShoppingCart)
		shoppingCartList = append(shoppingCartList, item)
		orderGoods, err := m.createOrdersGoods(item.GoodsID, item.SpecificationID, item.Quantity)
		if err != nil {
			orderGoodsList = append(orderGoodsList, &model.OrdersGoods{Error: err.Error()})
		} else {
			orderGoodsList = append(orderGoodsList, orderGoods)
		}

		//results[oredersGoods.OID]=append(results[oredersGoods.OID],oredersGoods)
	}

	confirmOrdersGoods, err := m.AnalyseOrdersGoodsList(oid, address, orderGoodsList)
	if err != nil {
		return nil, err
	}
	return &ShoppingCartResult{
		ConfirmOrdersGoods: confirmOrdersGoods,
		ShoppingCartList:   shoppingCartList,
	}, nil
}

// AfterSettlementUserBrokerage 退款，扣除相应的冻结金额，不用结算，佣金
func (m OrdersService) AfterSettlementUserBrokerage(tx *gorm.DB, orders *model.Orders) error {
	var err error
	//用户自己。下单者
	//Orm:=singleton.Orm()

	//var orders model.Orders
	//service.Get(Orm, OrderID, &orders)

	brokerage := m.Configuration.GetBrokerageConfiguration(orders.OID)

	//var orderUser model.User
	orderUser := dao.GetByPrimaryKey(tx, &model.User{}, orders.UserID).(*model.User)
	if orderUser.IsZero() {
		return gorm.ErrRecordNotFound
	}

	//leve1 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageLeve1).V)
	//leve2 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageLeve2).V)
	//leve3 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageLeve3).V)
	//leve4 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageLeve4).V)
	//leve5 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageLeve5).V)
	//leve6 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageLeve6).V)

	//leves := []uint{leve1, leve2, leve3, leve4, leve5, leve6}
	leves := []float64{brokerage.Leve1, brokerage.Leve2, brokerage.Leve3, brokerage.Leve4, brokerage.Leve5, brokerage.Leve6}

	for range leves {
		_user := dao.GetByPrimaryKey(tx, &model.User{}, orderUser.SuperiorID).(*model.User)
		if _user.ID <= 0 {
			break
		}

		err = m.Journal.DisableFreezeUserAmount(tx, _user.ID, journal.NewDataTypeOrder(orders.ID), orders.UserID)
		if err != nil {
			log.Println(err)
			return err
		}

		orderUser = _user
	}

	return err
}
func (m OrdersService) FirstSettlementUserBrokerage(tx *gorm.DB, orders model.Orders) error {
	var err error
	//用户自己。下单者
	//Orm:=singleton.Orm()

	//var orders model.Orders
	//service.Get(Orm, OrderID, &orders)
	brokerage := m.Configuration.GetBrokerageConfiguration(orders.OID)

	ogs, err := m.FindOrdersGoodsByOrdersID(tx, orders.ID)
	var Brokerage uint
	for i := range ogs {
		value := ogs[i]
		//var specification model.Specification
		//util.JSONToStruct(value.Specification, &specification)
		//Brokerage = Brokerage + value.TotalBrokerage
		if brokerage.Type == configuration.BrokeragePRODUCT {
			Brokerage = Brokerage + value.SellPrice
		}
		if brokerage.Type == configuration.BrokerageCUSTOM {
			Brokerage = Brokerage + value.TotalBrokerage
		}
	}

	//var orderUser model.User
	orderUser := dao.GetByPrimaryKey(tx, &model.User{}, orders.UserID).(*model.User)
	if orderUser.IsZero() {
		return gorm.ErrRecordNotFound
	}

	leves := []float64{brokerage.Leve1, brokerage.Leve2, brokerage.Leve3, brokerage.Leve4, brokerage.Leve5, brokerage.Leve6}

	//var OutBrokerageMoney int64 = 0
	for index, value := range leves {
		if value <= 0 {
			break
		}
		/*var _user model.User
		err = service.Get(tx, orderUser.SuperiorID, &_user)
		if err != nil {
			return err
		}*/
		_user := dao.GetByPrimaryKey(tx, &model.User{}, orderUser.SuperiorID).(*model.User)
		if _user.ID <= 0 {
			return nil
		}
		leveAmount := int64(math.Floor(value/float64(100)*float64(Brokerage) + 0.5))
		/*err = service.User.AddUserBlockAmount(tx, _user.ID, leveMenoy)
		if err != nil {
			log.Println(err)
			continue
		}*/

		//AddUserJournal(Orm, _user.ID, "佣金", strconv.Itoa(index+1)+"级用户", play.UserJournal_Type_LEVE, leveMenoy, extends.KV{Key: "OrdersID", Value: orders.ID}, u.ID)
		err = m.Journal.FreezeUserAmount(tx, _user.ID, "佣金", strconv.Itoa(index+1)+"级用户", leveAmount, journal.NewDataTypeOrder(orders.ID), orders.UserID)
		if err != nil {
			log.Println(err)
			continue
		}

		//OutBrokerageMoney = OutBrokerageMoney + leveMenoy
		workTime := time.Now().Unix() - orders.CreatedAt.Unix()
		m.MessageNotify.INComeNotify(_user, "来自"+strconv.Itoa(index+1)+"级用户，预计现金收入", strconv.Itoa(int(workTime/60/60))+"小时", "预计收入："+strconv.FormatFloat(float64(leveAmount)/float64(100), 'f', 2, 64)+"元")
		//fmt.Println("预计收入：" + strconv.FormatFloat(float64(leveMenoy)/float64(100), 'f', 2, 64) + "元")
		orderUser = _user
	}

	return err
}
func (m OrdersService) OrdersStockManager(db *gorm.DB, orders *model.Orders, isMinus bool) error {

	if orders.PostType == 2 {
		//线下订单，不去维护在线商品库存
		log.Println("线下订单，不去维护在线商品库存")
		return nil
	}

	//管理商品库存
	//Orm := singleton.Orm()
	//list []model.OrdersGoods

	list, _ := m.FindOrdersGoodsByOrdersID(db, orders.ID)
	for i := range list {
		value := list[i]
		var specification model.Specification = value.Specification
		//service.Get(Orm, value.SpecificationID, &specification)
		//util.JSONToStruct(value.Specification, &specification)
		var g model.Goods = value.Goods
		//service.Get(Orm, value.GoodsID, &goods)
		//util.JSONToStruct(value.Goods, &g)

		if isMinus {
			//减
			//UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
			//db.Model(&product).Updates(map[string]interface{}{"price": gorm.Expr("price * ? + ?", 2, 100)})
			err := dao.UpdateBy(db, &model.Specification{}, map[string]interface{}{"Stock": gorm.Expr(`"Stock" - ?`, value.Quantity)}, `"ID"=? and "Stock">=?`, specification.ID, value.Quantity)
			if err != nil {
				return err
			}
			err = dao.UpdateBy(db, &model.Goods{}, map[string]interface{}{"Stock": gorm.Expr(`"Stock" - ?`, value.Quantity)}, `"ID"=? and "Stock">=?`, g.ID, value.Quantity)
			if err != nil {
				return err
			}
		} else {
			//添加
			//Stock := int64(specification.Stock + value.Quantity)
			err := dao.UpdateByPrimaryKey(db, &model.Specification{}, specification.ID, map[string]interface{}{"Stock": gorm.Expr(`"Stock" + ?`, value.Quantity)})
			if err != nil {
				return err
			}
			err = dao.UpdateByPrimaryKey(db, &model.Goods{}, g.ID, map[string]interface{}{"Stock": gorm.Expr(`"Stock" + ?`, value.Quantity)})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (m OrdersService) Situation(StartTime, EndTime int64) interface{} {

	st := time.Unix(StartTime/1000, 0)
	st = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
	et := time.Unix(EndTime/1000, 0).Add(24 * time.Hour)
	et = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())

	Orm := db.Orm()

	type Result struct {
		TotalMoney uint `gorm:"column:TotalMoney"`
		TotalCount uint `gorm:"column:TotalCount"`
	}

	var result Result

	Orm.Table("Orders").Select(`SUM("PayMoney") as "TotalMoney",COUNT("ID") as "TotalCount"`).Where(`"CreatedAt">=?`, st).Where(`"CreatedAt"<?`, et).Where(map[string]interface{}{"IsPay": 1}).Find(&result)
	//fmt.Println(result)
	return result
}

/*
	func (service OrdersService) AddOrderBrokerageTemp(UserID uint, OrderNo string, Amount int64) error {
		var orderbrokerage model.OrderBrokerageTemp
		orderbrokerage.OrderNo = OrderNo
		orderbrokerage.Brokerage = uint(Amount)
		orderbrokerage.UserID = UserID
		err := dao.Create((singleton.Orm(), &orderbrokerage)
		return err
	}
*/
func (m OrdersService) AddOrdersPackage(db *gorm.DB, TotalMoney uint, UserID dao.PrimaryKey) (model.OrdersPackage, error) {

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
	err := dao.Create(db, &orderbrokerage)
	return orderbrokerage, err
}

// 确认收货
func (m OrdersService) TakeDeliver(OrdersID dao.PrimaryKey) error {
	Orm := db.Orm()

	//var orders model.Orders
	orders := dao.GetByPrimaryKey(Orm, entity.Orders, OrdersID).(*model.Orders)
	if orders.ID == 0 {

		return errors.New("订单不存在")
	}
	//下单状态,只有邮寄才能确认收货
	if orders.Status == model.OrdersStatusDeliver {

		tx := Orm.Begin()

		err := dao.UpdateByPrimaryKey(tx, entity.Orders, orders.ID, map[string]interface{}{"Status": model.OrdersStatusOrderOk, "ReceiptTime": time.Now()})
		if err != nil {
			tx.Rollback()
			return err
		}

		ogs, err := m.FindOrdersGoodsByOrdersID(tx, orders.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		var ogsList []*model.OrdersGoods
		for i := range ogs {
			value := ogs[i]
			//var specification model.Specification
			//util.JSONToStruct(value.Specification, &specification)
			//Brokerage = Brokerage + value.TotalBrokerage
			ogsList = append(ogsList, value)
		}

		/*err = service.CardItem.AddOrdersGoodsCardItem(tx, orders, ogs)
		if err != nil {
			tx.Rollback()
			return err, ""
		}*/

		//err = service.CardItem.AddOrdersGoodsCardItem()

		//Orm *gorm.DB, UserID uint, Brokerage uint, TargetID uint, PayMenoy uint
		err = m.Settlement.SettlementUser(tx, ogsList, orders)
		if err != nil {
			tx.Rollback()
			return err
		} else {
			tx.Commit()
			go func(ogs []*model.OrdersGoods) {
				for i := range ogs {
					value := ogs[i]
					_goods := value.Goods
					if _goods.ID != 0 {
						err = dao.UpdateByPrimaryKey(db.Orm(), entity.Goods, _goods.ID, &model.Goods{CountSale: _goods.CountSale + uint(value.Quantity)})
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

// 检查订单状态
func (m OrdersService) AnalysisOrdersStatus(context constrain.IServiceContext, ordersID dao.PrimaryKey, oid dao.PrimaryKey) error {

	Orm := db.Orm()

	//var orders model.Orders
	orders := dao.GetByPrimaryKey(Orm, entity.Orders, ordersID).(*model.Orders)
	if orders.ID == 0 {
		//return errors.New("订单不存在"), ""
		return nil
	}
	if orders.Status == model.OrdersStatusOrder {

		if time.Now().Unix() >= orders.CreatedAt.Add(3*time.Hour*24).Unix() {
			//一直处于下单状态超过3天，没有付款，自动关闭订单，并加回库存
			err := dao.UpdateByPrimaryKey(Orm, entity.Orders, orders.ID, map[string]interface{}{"Status": model.OrdersStatusClosed})
			if err != nil {
				return err
			}
			//管理商品库存
			err = m.OrdersStockManager(Orm, orders, false)
			if err != nil {
				return err
			}
		}

	} else if orders.Status == model.OrdersStatusDeliver {
		if time.Now().Unix() >= orders.DeliverTime.Add(15*time.Hour*24).Unix() {
			//等待收货时间超过15天，自动订单完成
			//dao.UpdateByPrimaryKey(Orm, orders.ID, &model.Orders{}, map[string]interface{}{"Status": model.OrdersStatusOrderOk, "ReceiptTime": time.Now()})
			//管理商品库存
			//service.Goods.OrdersStockManager(orders, false)
			err := m.TakeDeliver(ordersID)
			if err != nil {
				return err
			}
		}

	} else if orders.Status == model.OrdersStatusCancel {
		if time.Now().Unix() >= orders.UpdatedAt.Add(5*time.Hour*24).Unix() {
			//订单已经支付，用户申请了取消订单，超过5天，自动取消
			_, err := m.CancelOk(context, ordersID)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// 发货
func (m OrdersService) Deliver(ShipTitle, ShipKey, ShipNo string, OrdersID dao.PrimaryKey) error {
	Orm := db.Orm().Begin()

	//var orders model.Orders
	orders := dao.GetByPrimaryKey(Orm, &model.Orders{}, OrdersID).(*model.Orders)
	if orders.ID == 0 {
		Orm.Rollback()
		return errors.New("订单不存在")
	}
	if orders.IsPay != 1 {
		Orm.Rollback()
		return errors.New("订单没有支付")
	}

	expressCompany := dao.GetBy(Orm, &model.ExpressCompany{}, map[string]any{"Key": ShipKey}).(*model.ExpressCompany)
	if expressCompany.IsZero() {
		Orm.Rollback()
		return errors.New("找不到快递信息")
	}

	shipping := &model.OrdersShipping{
		OID:     orders.OID,
		OrderNo: orders.OrderNo,
		Title:   ShipTitle,
		Image:   "",
		No:      ShipNo,
		Name:    expressCompany.Name,
		Key:     expressCompany.Key,
	}
	err := dao.Create(Orm, shipping)
	if err != nil {
		Orm.Rollback()
		return err
	}

	orders.DeliverTime = time.Now()
	orders.Status = model.OrdersStatusDeliver
	err = dao.UpdateByPrimaryKey(Orm, &model.Orders{}, OrdersID, &model.Orders{DeliverTime: orders.DeliverTime, Status: orders.Status})
	if err != nil {
		Orm.Rollback()
		return err
	}

	/*ogs, err := service.FindOrdersGoodsByOrdersID(db.Orm(), orders.ID)
	if err != nil {
		Orm.Rollback()
		return err
	}*/

	/*as := service.MessageNotify.OrderDeliveryNotify(orders, ogs, wxConfig)
	if as.Code != result.Success {

		err = errors.New(as.Message)
	}*/
	Orm.Commit()
	return err
}
func (m OrdersService) GetOrdersPackageByOrderNo(OrderNo string) model.OrdersPackage {
	Orm := db.Orm()
	var orders model.OrdersPackage
	Orm.Where(&model.OrdersPackage{OrderNo: OrderNo}).First(&orders)
	return orders
}

func (m OrdersService) GetSupplyOrdersByOrderNo(OrderNo string) model.SupplyOrders {
	Orm := db.Orm()
	var orders model.SupplyOrders
	Orm.Where(&model.SupplyOrders{OrderNo: OrderNo}).First(&orders)
	return orders
}

func (m OrdersService) ListOrdersStatusCount(UserID dao.PrimaryKey, Status []string) (TotalRecords int64) {
	Orm := db.Orm()
	var orders []model.Orders
	db := Orm.Model(model.Orders{})

	now := time.Now()
	ts := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	te := ts.Add(24 * time.Hour)

	db = db.Where(`"UpdatedAt">=? and "UpdatedAt"<?`, ts, te)
	db = db.Where(`"UserID"=?`, UserID)

	if len(Status) > 0 {
		db = db.Where(`"Status" = ?`, Status[0])
		for index, value := range Status {
			if index != 0 {
				db = db.Or(`"Status" = ?`, value)
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

func (m OrdersService) ListCollageRecord(UserID dao.PrimaryKey, Index int) []CollageRecord {
	Orm := db.Orm()

	db := Orm.Raw(`
SELECT
o."ID" AS "OrdersID",cr."No" as "No",cr."UserID" as "UserID",cr."Collager" as "Collager",og."Favoured" as "Favoured",og."Goods" as "Goods",og."Specification" as "Specification",o."Status" AS "Status",o."IsPay" AS "IsPay",og."Quantity" as "Quantity",
(SELECT mcr."CreatedAt" FROM "CollageRecord" mcr WHERE mcr."No"=cr."No" AND mcr."Collager">0) AS "CreatedAt",
(SELECT COUNT(mo."IsPay") FROM "CollageRecord" mcr,"Orders" mo WHERE mcr."No"=cr."No" AND mo."OrderNo"=mcr."OrderNo") AS COUNT,
(SELECT SUM(mo."IsPay") FROM "CollageRecord" mcr,"Orders" mo WHERE mcr."No"=cr."No" AND mo."OrderNo"=mcr."OrderNo") AS "IsPaySUM"
FROM
"User" u,"Orders" o,"CollageRecord" cr,"OrdersGoods" og
WHERE
cr."UserID"=? AND u."ID"=cr."UserID" AND o."OrderNo"=cr."OrderNo" AND og."OrdersGoodsNo"=cr."OrdersGoodsNo"
GROUP BY cr."No"
`, UserID)
	//db := Orm.Raw("SELECT o.ID AS OrdersID,cr.No,cr.UserID,cr.Collager,cr.IsPay,sdf.*,og.Favoured,og.Goods,cr.CreatedAt as CreatedAt from User u,Orders o,CollageRecord cr,OrdersGoods og,(SELECT COUNT(cr.NO) AS COUNT,SUM(cr.IsPay) AS SUM FROM CollageRecord cr GROUP BY cr.NO) AS sdf WHERE cr.UserID=? AND u.ID=cr.UserID AND o.OrderNo=cr.OrderNo AND og.OrdersGoodsNo=cr.OrdersGoodsNo GROUP BY cr.No", UserID)

	packs := make([]CollageRecord, 0)
	db = db.Limit(play.Paging).Offset(play.Paging * Index).Order(`"CreatedAt" desc`)
	db.Scan(&packs)

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

func (m OrdersService) ListOrders(queryParam *serviceargument.ListOrdersQueryParam, oid dao.PrimaryKey, fieldOrder clause.OrderByColumn, pageNo int, pageSize int) (*result.Pagination, error) {
	if pageSize <= 0 {
		pageSize = 10
	}

	pageIndex := pageNo - 1
	if pageIndex < 0 {
		pageIndex = 0
	}

	Orm := db.Orm()
	var orders []model.Orders

	db := Orm.Model(model.Orders{})
	if oid > 0 {
		db = db.Where(`"OID"=?`, oid)
	}

	if queryParam != nil {
		if queryParam.UserID != 0 {
			db = db.Where(`"UserID"=?`, queryParam.UserID)
		}

		if queryParam.StartDate.IsZero() == false && queryParam.EndDate.IsZero() == false {
			db = db.Where(`"CreatedAt">=? and "CreatedAt"<=?`, queryParam.StartDate, queryParam.EndDate)
		}

		if len(queryParam.Status) > 0 {
			db = db.Where(`"Status" in ?`, queryParam.Status)
		}
	}

	var recordsTotal int64

	db = db.Count(&recordsTotal).Limit(pageSize).Offset(pageSize * pageIndex).Order(fieldOrder).Find(&orders)

	results := make([]extends.OrdersDetail, 0)
	for _, value := range orders {
		/*	pack := struct {
			Orders          model.Orders
			User            model.User
			OrdersGoodsList []*model.OrdersGoods //[]model.OrdersGoods
		}{}*/
		pack := extends.OrdersDetail{}

		pack.Orders = value

		ordersShippingList := make([]*model.OrdersShipping, 0)
		dao.Find(Orm, &model.OrdersShipping{}).Where(`"OrderNo"=?`, value.OrderNo).Result(&ordersShippingList)
		pack.OrdersShippingList = ordersShippingList

		pack.User = *(dao.GetByPrimaryKey(Orm, &model.User{}, value.UserID).(*model.User))

		ogs, err := m.FindOrdersGoodsByOrdersID(Orm, value.ID)
		if err != nil {
			return nil, err
		}
		pack.OrdersGoodsList = ogs
		//:todo 拼单
		//og := ogs[0]
		//pack.CollageUsers = service.FindOrdersGoodsByCollageUser(og.CollageNo)
		results = append(results, pack)
	}
	return result.NewPagination(pageNo, pageSize, recordsTotal, results), nil
}

func (m OrdersService) OrderPaySuccess(totalFee uint, outTradeNo string, transactionId string, payTime time.Time, ordersType model.OrdersType) (string, error) {

	//Orm := singleton.Orm()

	//TotalFee, _ := strconv.ParseUint(result["total_fee"], 10, 64)
	//OrderNo := result["out_trade_no"]
	//TimeEnd := result["time_end"]
	//attach := result["attach"]

	if ordersType == model.OrdersTypeSupply {
		//充值的，目前只涉及到门店自主核销的时候，才需要用到充值
		orders := m.GetSupplyOrdersByOrderNo(outTradeNo)
		if orders.IsPay == 0 {
			tx := db.Orm().Begin()
			err := dao.UpdateByPrimaryKey(tx, entity.SupplyOrders, orders.ID, &model.SupplyOrders{PayTime: payTime, IsPay: 1, PayMoney: totalFee})
			if err != nil {
				tx.Rollback()
				return "", err
			} else {
				if strings.EqualFold(orders.Type, play.SupplyType_Store) {
					_, err := m.Journal.AddStoreJournal(tx, orders.StoreID, "门店", "充值", play.StoreJournal_Type_CZ, int64(totalFee), orders.ID)
					if err != nil {
						tx.Rollback()
						return "", err
					} else {
						tx.Commit()
						return "已经支付成功", nil
					}
				} else {
					tx.Commit()
					strings.EqualFold(orders.Type, play.SupplyType_User)
					return "", fmt.Errorf("未实现的数据类型:%s", orders.Type)
				}

			}
		} else {
			return "", errors.New("订单已经处理或过期")
		}

	} else if ordersType == model.OrdersTypeGoodsPackage { //合并商品订单
		tx := db.Orm().Begin()
		ordersPackage := m.GetOrdersPackageByOrderNo(outTradeNo)
		if ordersPackage.TotalPayMoney == totalFee {
			//var OrderNoList []string
			//util.JSONToStruct(ordersPackage.OrderList, &OrderNoList)

			err := dao.UpdateByPrimaryKey(tx, entity.OrdersPackage, ordersPackage.ID, &model.OrdersPackage{IsPay: 1})
			if err != nil {
				tx.Rollback()
				return "", err
			}

			OrderList := repository.OrdersDao.GetOrdersByOrdersPackageNo(ordersPackage.OrderNo)

			for index := range OrderList {
				//orders := service.GetOrdersByOrderNo(value)
				df, msg := m.ProcessingOrders(tx, OrderList[index], transactionId, payTime)
				if df == false {
					tx.Rollback()
					return "", errors.New(msg)
				}
			}
			tx.Commit()
			return "已经支付成功", nil
		} else {
			tx.Commit()
			return "", errors.New("金额不正确或订单不允许")
		}

	} else if ordersType == model.OrdersTypeGoods { //商品订单
		//orders.PayMoney == total_fee.
		tx := db.Orm().Begin()
		orders := repository.OrdersDao.GetOrdersByOrderNo(outTradeNo)
		if orders.PayMoney == totalFee {
			su, msg := m.ProcessingOrders(tx, orders, transactionId, payTime)
			if su == false {
				tx.Rollback()
				return "", errors.New(msg)
			}
			tx.Commit()
			return msg, nil
		} else {
			tx.Commit()
			return "", errors.New("金额不正确或订单不允许")
		}

	} else {
		return "", fmt.Errorf("未实现的订单类型:%s", ordersType)
	}

}

func (m OrdersService) ProcessingOrders(tx *gorm.DB, orders model.Orders, transactionId string, payTime time.Time) (Success bool, Message string) {
	if orders.IsPay == 0 {
		if orders.Status == model.OrdersStatusOrder {
			var err error

			err = dao.UpdateByPrimaryKey(tx, entity.Orders, orders.ID, &model.Orders{PayTime: payTime, TransactionID: transactionId, IsPay: 1, Status: model.OrdersStatusPay})
			if err != nil {
				return false, err.Error()
			}
			err = m.FirstSettlementUserBrokerage(tx, orders)
			if err != nil {

				return false, err.Error()
			}

			return true, "已经支付成功"

		} else {

			return false, "金额不正确或订单不允许"
		}
	} else {

		return false, "订单已经处理或过期"
	}
}

// BuyCollageOrders 拼单购买
func (m OrdersService) BuyCollageOrders(ctx constrain.IContext, UserID, GoodsID, SpecificationID dao.PrimaryKey, Quantity uint) error {
	Orm := db.Orm()
	//var goods model.Goods
	//var specification model.Specification
	//var expresstemplate model.ExpressTemplate

	goods := dao.GetByPrimaryKey(Orm, entity.Goods, GoodsID).(*model.Goods)
	if goods.IsZero() {
		return gorm.ErrRecordNotFound
	}
	specification := dao.GetByPrimaryKey(Orm, entity.Specification, SpecificationID).(*model.Specification)
	if specification.IsZero() {
		return gorm.ErrRecordNotFound
	}
	if specification.GoodsID != goods.ID {
		return errors.New("产品与规格不匹配")
	}

	//shoppingCart := model.ShoppingCart{}
	//shoppingCart.Quantity = Quantity
	//shoppingCart.Specification = util.StructToJSON(specification)
	//shoppingCart.Goods = util.StructToJSON(goods)
	//shoppingCart.UserID = UserID

	ordersGoods, err := m.createOrdersGoods(goods.ID, specification.ID, Quantity)
	if err != nil {
		return err
	}

	//ordersGoods.CollageNo = tool.UUID()
	collage := m.Collage.GetCollageByGoodsID(goods.ID, goods.OID)
	if collage.ID != 0 && collage.TotalNum > 0 {

		favoured := sqltype.Discount{Name: strconv.Itoa(collage.Num) + "人拼团", Target: util.StructToJSON(collage), TypeName: "Collage", Discount: uint(collage.Discount)}
		ordersGoods.Discounts = []sqltype.Discount{favoured}
	}

	ogs := make([]*model.OrdersGoods, 0)
	ogs = append(ogs, ordersGoods)
	//Session.Attributes.Put(gweb.AttributesKey(string(play.SessionConfirmOrders)), &ogs)
	ctx.Redis().Set(ctx, redis.NewConfirmOrders(ctx.UID()), &ogs, 24*time.Hour)
	return nil

}

// CreateOrdersGoods 从商品外直接购买，生成OrdersGoods，添加到 play.SessionConfirmOrders
func (m OrdersService) CreateOrdersGoods(ctx constrain.IContext, UserID, GoodsID, SpecificationID dao.PrimaryKey, Quantity uint) ([]*model.OrdersGoods, error) {
	Orm := db.Orm()
	//var goods model.Goods
	//var specification model.Specification
	//var expresstemplate model.ExpressTemplate

	g := dao.GetByPrimaryKey(Orm, &model.Goods{}, GoodsID).(*model.Goods)
	if g.IsZero() {
		return nil, gorm.ErrRecordNotFound
	}
	specification := dao.GetByPrimaryKey(Orm, &model.Specification{}, SpecificationID).(*model.Specification)
	if specification.IsZero() {
		return nil, gorm.ErrRecordNotFound
	}
	if specification.GoodsID != g.ID {
		return nil, errors.New("产品与规格不匹配")
	}

	//shoppingCart := model.ShoppingCart{}
	//shoppingCart.Quantity = Quantity
	//shoppingCart.Specification = util.StructToJSON(specification)
	//shoppingCart.Goods = util.StructToJSON(goods)
	//shoppingCart.UserID = UserID

	ordersGoods, err := m.createOrdersGoods(g.ID, specification.ID, Quantity)
	if err != nil {
		return nil, err
	}

	ogs := make([]*model.OrdersGoods, 0)
	ogs = append(ogs, ordersGoods)
	//Session.Attributes.Put(gweb.AttributesKey(string(play.SessionConfirmOrders)), &ogs)
	//ctx.Redis().Set(ctx, redis.NewConfirmOrders(ctx.UID()), &ogs, 24*time.Hour)
	return ogs, nil

}

// 从购买车提交的订单，通过 ShoppingCart ID,生成  OrdersGoods 列表,添加到 play.SessionConfirmOrders
/*func (service OrdersService) AddCartOrdersByShoppingCartIDs(ctx constrain.IContext, UserID dao.PrimaryKey, IDs []string) error {
	scs := service.ShoppingCart.GetGSIDs(UserID, IDs)
	ogs := make([]*extends.OrdersGoods, 0)
	for _, value := range scs {
		ordersGoods, err := service.createOrdersGoods(&value)
		if err != nil {
			return err
		}
		ogs = append(ogs, ordersGoods)
	}
	ctx.Redis().Set(ctx, redis.NewConfirmOrders(ctx.UID()), &ogs, 24*time.Hour)
	return nil
}*/
func (m OrdersService) ConvertOrdersGoods(data *model.OrdersGoods) (*model.OrdersGoods, error) {
	return data, nil
	/*ordersGoods := &extends.OrdersGoodsMix{}
	//var g model.Goods=data.Goods
	//var specification model.Specification=data.Specification
	//var discountList []sqltype.Discount = data.Discounts

	if specification.GoodsID != g.ID {
		return nil, errors.New("产品规格不匹配")
	}
	if specification.ID == 0 {
		return nil, errors.New("找不到规格")
	}
	if specification.Stock-data.Quantity < 0 {
		return nil, errors.New("库存不足")
	}

	ordersGoods.Specification = &specification
	ordersGoods.Goods = &g
	ordersGoods.OID = g.OID
	ordersGoods.Quantity = data.Quantity

	ordersGoods.Image = data.Image

	ordersGoods.OrdersGoodsNo = data.OrdersGoodsNo
	ordersGoods.Discounts = data.Discounts
	ordersGoods.Status = data.Status
	ordersGoods.RefundInfo = data.RefundInfo
	ordersGoods.OrdersID = data.OrdersID
	ordersGoods.CostPrice = data.CostPrice
	ordersGoods.SellPrice = data.SellPrice
	ordersGoods.TotalBrokerage = data.TotalBrokerage

	goodsSkuData := m.goodsSkuData(db.Orm(), g.ID, specification.LabelIndex)
	ordersGoods.SkuImages = goodsSkuData.SkuImages
	ordersGoods.SkuLabelMap = goodsSkuData.SkuLabelMap
	ordersGoods.SkuLabelDataMap = goodsSkuData.SkuLabelDataMap
	return ordersGoods, nil*/
}
func (m OrdersService) createOrdersGoods(goodsID dao.PrimaryKey, specificationID dao.PrimaryKey, quantity uint) (*model.OrdersGoods, error) {
	ordersGoods := &model.OrdersGoods{}
	//var goods model.Goods
	//var specification model.Specification
	//var timesell model.TimeSell

	g := dao.GetByPrimaryKey(db.Orm(), &model.Goods{}, goodsID).(*model.Goods)
	if g.IsZero() {
		return nil, errors.New("无效的商品或商品已经不存在")
	}

	specification := dao.GetByPrimaryKey(db.Orm(), &model.Specification{}, specificationID).(*model.Specification)
	if specification.ID == 0 {
		//return errors.New("找不到规格")
		return nil, errors.New("无效的商品规格或商品规格已经不存在")
	}
	/*err := util.JSONToStruct(shoppingCart.Goods, &goods)
	if err != nil {
		return nil, err
	}*/
	/*err = util.JSONToStruct(shoppingCart.Specification, &specification)
	if err != nil {
		return nil, err
	}*/

	/*err := service.Goods.Get(Orm, shoppingCart.GoodsID, &goods)
	if err != nil {
		ordersGoods.AddError(err.Error())
	}

	err = service.Goods.Get(Orm, shoppingCart.SpecificationID, &specification)
	if err != nil {
		ordersGoods.AddError(err.Error())
	}*/
	if specification.GoodsID != g.ID {
		//return errors.New("产品规格不匹配")
		return nil, errors.New("商品规格无效")
	}

	if specification.Stock-quantity < 0 {
		//return errors.New(specification.Label + "库存不足")
		//ordersGoods.AddError("库存不足")
		//shoppingCart.Quantity = specification.Stock
		return nil, errors.New("库存不足")
	}

	ordersGoods.Specification = *specification
	ordersGoods.Goods = *g
	ordersGoods.OID = g.OID
	//ordersGoods.GoodsID = goods.ID
	//ordersGoods.SpecificationID = specification.ID
	ordersGoods.Quantity = quantity
	ordersGoods.CostPrice = specification.MarketPrice
	ordersGoods.SellPrice = specification.MarketPrice
	ordersGoods.OrdersGoodsNo = tool.UUID()
	ordersGoods.Discounts = m.Goods.GetDiscounts(g.ID, g.OID)

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

	return ordersGoods, nil
}

func (m OrdersService) AddOrders(db *gorm.DB, orders *model.Orders, list []*model.OrdersGoods) error {
	err := dao.Create(db, orders)
	if err != nil {
		return err
	}
	for i := range list {
		value := list[i]
		//value.OrdersGoods.OrdersID = orders.ID
		//value.OrdersGoods.Discounts = value.Discounts //util.StructToJSON((&value).Discounts)
		//err = dao.Create(db, &((&value).OrdersGoods))
		/*var goodsSkuList []model.GoodsSku
		for ii := 0; ii < len(value.Specification.LabelIndex); ii++ {
			labelID := value.Specification.LabelIndex[ii]
			skuLabelData := value.SkuLabelDataMap[labelID]
			skuLabel := value.SkuLabelMap[skuLabelData.GoodsSkuLabelID]
			goodsSkuList = append(goodsSkuList, model.GoodsSku{
				GoodsSkuLabel:     skuLabel,
				GoodsSkuLabelData: skuLabelData,
			})
		}*/
		err = dao.Create(db, &model.OrdersGoods{
			OID:            value.OID,
			OrdersGoodsNo:  value.OrdersGoodsNo,
			Status:         value.Status,
			RefundID:       0,
			OrdersID:       orders.ID,
			Image:          value.Image,
			Goods:          value.Goods,
			Specification:  value.Specification,
			GoodsSkus:      value.GoodsSkus,
			Discounts:      value.Discounts,
			Quantity:       value.Quantity,
			CostPrice:      value.CostPrice,
			SellPrice:      value.SellPrice,
			TotalBrokerage: value.TotalBrokerage,
			//Error:          value.OrdersGoods,
		})
		if err != nil {
			return err
		}

		//var goods model.Goods
		//var specification model.Specification
		//err = util.JSONToStruct(value.OrdersGoods.Goods, &goods)
		/*if err != nil {
			return err
		}*/
		/*err = util.JSONToStruct(value.OrdersGoods.Specification, &specification)
		if err != nil {
			return err
		}*/
		err = m.ShoppingCart.DeleteByUserIDAndGoodsIDAndSpecificationID(db, orders.UserID, value.Goods.ID, value.Specification.ID)
		if err != nil {
			return err
		}
	}

	err = m.OrdersStockManager(db, orders, true)
	if err != nil {
		return err
	}
	return nil

}
func (m OrdersService) ChangeOrdersPayMoney(context constrain.IContext, PayMoney float64, OrdersID dao.PrimaryKey, oid dao.PrimaryKey) (Success result.ActionResultCode, Message string) {
	tx := db.Orm().Begin()

	orders := repository.OrdersDao.GetOrdersByID(OrdersID)
	pm := payment.NewPayment(context, oid, orders.PayMethod)

	if strings.EqualFold(orders.PrepayID, "") == false {
		err := pm.CloseOrder(orders.OrderNo)
		if err != nil {
			tx.Rollback()
			return result.Fail, err.Error()
		}
		/*success, message := pm.CloseOrder(orders.OrderNo)
		if success == false {
			tx.Rollback()
			return result.Fail, message
		}*/
	}

	err := dao.UpdateByPrimaryKey(tx, &model.Orders{}, OrdersID, map[string]interface{}{"PayMoney": uint(PayMoney * 100), "PrepayID": "", "OrderNo": tool.UUID()})
	if err != nil {
		tx.Rollback()
		return result.Fail, err.Error()
	}

	tx.Commit()

	return result.Success, "订单金额修改成功"

}

func (m OrdersService) AnalyseOrdersGoodsListByOrders(orders *model.Orders, address *model.Address) (*extends.ConfirmOrdersGoods, error) {
	orderGoodsList := make([]*model.OrdersGoods, 0)
	ordersGoodsList, _ := m.FindOrdersGoodsByOrdersID(db.Orm(), orders.ID)
	for i := 0; i < len(ordersGoodsList); i++ {
		ordersGoods, err := m.ConvertOrdersGoods(ordersGoodsList[i])
		if err != nil {
			return nil, err
		}
		orderGoodsList = append(orderGoodsList, ordersGoods)
	}
	confirmOrdersGoods, err := m.AnalyseOrdersGoodsList(orders.OID, address, orderGoodsList)
	if err != nil {
		return nil, err
	}
	return confirmOrdersGoods, err
}

// AnalyseOrdersGoodsList 订单分析，
func (m OrdersService) AnalyseOrdersGoodsList(oid dao.PrimaryKey, addressee *model.Address, orderGoods []*model.OrdersGoods) (*extends.ConfirmOrdersGoods, error) {

	/*oslist := make(map[dao.PrimaryKey][]*extends.OrdersGoods)
	for index, v := range orderGoods {
		items := oslist[v.OID]
		if items == nil {
			oslist[v.OID] = make([]*extends.OrdersGoods, 0)
		}
		oslist[v.OID] = append(oslist[v.OID], orderGoods[index])
	}*/

	outResult := &extends.ConfirmOrdersGoods{}

	//var org model.Organization
	//org := dao.GetByPrimaryKey(singleton.Orm(), &model.Organization{}, key).(*model.Organization)

	analyseResult, err := m.analyseOne(oid, addressee, orderGoods)
	if err != nil {
		return nil, err
	}

	//result.Error = Error
	outResult.OrdersGoodsInfos = analyseResult.OrdersGoodsInfo
	outResult.FavouredPrice = analyseResult.FavouredPrice
	outResult.FullCutAll = analyseResult.FullCutAll
	outResult.GoodsPrice = analyseResult.GoodsPrice
	outResult.ExpressPrice = analyseResult.ExpressPrice
	outResult.FullCut = analyseResult.FullCut
	outResult.Address = addressee

	outResult.TotalAmount = analyseResult.GoodsPrice - analyseResult.FullCutAll + analyseResult.ExpressPrice

	return outResult, nil
}

type GoodsSkuData struct {
	SkuImages       []string
	SkuLabelMap     map[dao.PrimaryKey]*model.GoodsSkuLabel
	SkuLabelDataMap map[dao.PrimaryKey]*model.GoodsSkuLabelData
}

func (m OrdersService) goodsSkuData(tx *gorm.DB, goodsID dao.PrimaryKey, specificationLabelIndex sqltype.Array[dao.PrimaryKey]) *GoodsSkuData {

	goodsSkuLabelMap := make(map[dao.PrimaryKey]*model.GoodsSkuLabel)
	{
		goodsSkuLabelList := dao.Find(tx, &model.GoodsSkuLabel{}).Where(`"GoodsID"=?`, goodsID).List()
		for i := range goodsSkuLabelList {
			goodsSkuLabelMap[goodsSkuLabelList[i].Primary()] = goodsSkuLabelList[i].(*model.GoodsSkuLabel)
		}
	}
	skuImages := make([]string, 0)
	goodsSkuLabelDataMap := make(map[dao.PrimaryKey]*model.GoodsSkuLabelData)
	{
		goodsSkuLabelDataList := dao.Find(tx, &model.GoodsSkuLabelData{}).Where(`"GoodsID"=?`, goodsID).List()
		for i := range goodsSkuLabelDataList {
			item := goodsSkuLabelDataList[i].(*model.GoodsSkuLabelData)

			for _, labelIndex := range specificationLabelIndex {
				if item.ID == labelIndex {
					goodsSkuLabelDataMap[item.ID] = item

					if v, ok := goodsSkuLabelMap[item.GoodsSkuLabelID]; ok {
						if len(item.Image) > 0 && v.Image {
							skuImages = append(skuImages, item.Image)
						}
					}
					break
				}
			}
		}
	}

	return &GoodsSkuData{
		SkuImages:       skuImages,
		SkuLabelMap:     goodsSkuLabelMap,
		SkuLabelDataMap: goodsSkuLabelDataMap,
	}
}

// 订单分析，
func (m OrdersService) analyseOne(OID dao.PrimaryKey, address *model.Address, list []*model.OrdersGoods) (*extends.AnalyseResult, error) {

	analyseResult := &extends.AnalyseResult{}

	Orm := db.Orm()

	fullCuts := m.FullCut.FindOrderByAmountDesc(Orm, OID)

	//可以使用满减的金额
	FullCutPrice := uint(0)
	//FavouredPrice := uint(0)

	oggs := make([]*model.OrdersGoods, 0)

	expresstemplateMap := make(map[dao.PrimaryKey]model.ExpressTemplateNMW)

	for index := range list {
		item := list[index]
		//value.ID = 5445
		//var goods model.Goods
		//var specification model.Specification

		//util.JSONToStruct(value.Goods, &goods)
		//util.JSONToStruct(value.Specification, &specification)

		/*if value.ElementStatus.IsError {
			oggs = append(oggs, value)
			continue
		}*/

		if len(item.Error) > 0 {
			oggs = append(oggs, item)
			continue
		}

		var value = item

		goodsSkuData := m.goodsSkuData(Orm, value.Goods.ID, value.Specification.LabelIndex)
		//value.SkuImages = goodsSkuData.SkuImages
		//value.SkuLabelMap = goodsSkuData.SkuLabelMap
		//value.SkuLabelDataMap = goodsSkuData.SkuLabelDataMap

		var goodsSkuList []model.GoodsSku
		for ii := 0; ii < len(value.Specification.LabelIndex); ii++ {
			labelID := value.Specification.LabelIndex[ii]
			skuLabelData := goodsSkuData.SkuLabelDataMap[labelID]
			skuLabel := goodsSkuData.SkuLabelMap[skuLabelData.GoodsSkuLabelID]
			goodsSkuList = append(goodsSkuList, model.GoodsSku{
				GoodsSkuLabel:     skuLabel,
				GoodsSkuLabelData: skuLabelData,
			})
		}

		value.GoodsSkus = goodsSkuList

		if len(goodsSkuData.SkuImages) > 0 {
			value.Image = goodsSkuData.SkuImages[len(goodsSkuData.SkuImages)-1]
		} else {
			if len(value.Goods.Images) > 0 {
				value.Image = value.Goods.Images[0]
			}
		}

		/*if PostType == 1 {
			//邮寄时，才判断库存
			if int64(value.Specification.Stock-value.Quantity) < 0 {
				//Error = errors.New(value.Specification.Label + "库存不足")
				//value.AddError(Error.Error())
				return nil, errors.New(value.Specification.Label + "库存不足")
			}
		}*/

		//value.Goods = util.StructToJSON(goods)
		//value.Specification = util.StructToJSON(specification)

		Price := value.Specification.MarketPrice * uint(value.Quantity)

		value.CostPrice = value.Specification.MarketPrice
		value.SellPrice = value.Specification.MarketPrice
		//value.TotalBrokerage =

		//ogss

		/*var discounts []extends.Discount
		if strings.EqualFold(value.Discounts, "") == false {
			util.JSONToStruct(value.Discounts, &discounts)
		}*/

		//计算价格以及优惠
		if len(value.Discounts) > 0 {
			for i := range value.Discounts {
				discount := value.Discounts[i]
				Price = uint(util.Rounding45(float64(Price)-(float64(Price)*(float64(discount.Discount)/float64(100))), 2))
				analyseResult.GoodsPrice = analyseResult.GoodsPrice + Price
				Favoured := uint(util.Rounding45(float64(value.SellPrice)*(float64(discount.Discount)/float64(100)), 2))
				analyseResult.FavouredPrice = analyseResult.FavouredPrice + (Favoured * uint(value.Quantity))

				value.SellPrice = value.SellPrice - Favoured
			}
			//value.Discounts = make([]sqltype.Discount, 0)
			//value.Discounts = value.Discounts

		} else {
			analyseResult.GoodsPrice = analyseResult.GoodsPrice + Price
			FullCutPrice = FullCutPrice + Price
		}
		oggs = append(oggs, value)
		//ogss=append(ogss,ogs)

		//计算快递费，重量要加上数量,先计算规格的重，再计算购买的重量
		weight := (value.Specification.Num * value.Specification.Weight) * uint(value.Quantity)

		if value.Goods.ExpressTemplateID == 0 {
			//Error = errors.New("找不到快递模板")
			//value.AddError(Error.Error())
			return nil, errors.New("找不到快递模板")
		} else {
			//为每个订单设置三种计价方式
			if _, o := expresstemplateMap[value.Goods.ExpressTemplateID]; o == false {
				nmw := model.ExpressTemplateNMW{}
				nmw.N = nmw.N + int(value.Quantity)
				nmw.M = nmw.M + int(Price) //市场价X数量
				nmw.W = nmw.W + int(weight)
				expresstemplateMap[value.Goods.ExpressTemplateID] = nmw
			} else {
				nmw := expresstemplateMap[value.Goods.ExpressTemplateID]
				nmw.N = nmw.N + int(value.Quantity)
				nmw.M = nmw.M + int(Price) //市场价X数量
				nmw.W = nmw.W + int(weight)
				expresstemplateMap[value.Goods.ExpressTemplateID] = nmw
			}
		}
	}

	//计算快满减
	for index, value := range fullCuts {
		if FullCutPrice >= value.Amount {
			analyseResult.FullCutAll = value.CutAmount
			//返回满减的值
			analyseResult.FullCut = fullCuts[index]
			break
		}
	}

	if address.IsEmpty() {
		//return nil, errors.New("地址不能为空")
	}
	//计算快递费

	for ID, value := range expresstemplateMap {
		//var expresstemplate model.ExpressTemplate
		expressTemplate := dao.GetByPrimaryKey(Orm, &model.ExpressTemplate{}, ID).(*model.ExpressTemplate)

		etFree := make([]model.ExpressTemplateFreeItem, 0)
		err := json.Unmarshal([]byte(expressTemplate.Free), &etFree)
		if err != nil {
			return nil, err
		}

		var expressTemplateFreeItem *model.ExpressTemplateFreeItem

	al:
		//从包邮列表中，找出一个计费方式
		for _, expFValue := range etFree {

			for _, expFAValue := range expFValue.Areas {
				if strings.EqualFold(address.ProvinceName, expFAValue) {
					expressTemplateFreeItem = &expFValue
					break al
				}
			}

		}

		if expressTemplateFreeItem != nil && expressTemplateFreeItem.IsFree(expressTemplate, value) {
			//有包邮项目
			analyseResult.ExpressPrice = 0

		} else {
			//无包邮项目

			etTemplate := model.ExpressTemplateTemplate{}
			err = json.Unmarshal([]byte(expressTemplate.Template), &etTemplate)
			if err != nil {
				return nil, err
			}

			var expressTemplateItem *model.ExpressTemplateItem

		alt:
			for _, expFValue := range etTemplate.Items {

				for _, expFAValue := range expFValue.Areas {
					if strings.EqualFold(address.ProvinceName, expFAValue) {
						expressTemplateItem = &expFValue
						break alt
					}
				}

			}

			if expressTemplateItem != nil {
				analyseResult.ExpressPrice = analyseResult.ExpressPrice + expressTemplateItem.CalculateExpressPrice(expressTemplate, value)
			} else {
				analyseResult.ExpressPrice = analyseResult.ExpressPrice + etTemplate.Default.CalculateExpressPrice(expressTemplate, value)
			}
		}
	}
	analyseResult.OrdersGoodsInfo = oggs
	return analyseResult, nil

}

func (m OrdersService) AddCartOrders(UserID dao.PrimaryKey, GoodsID, SpecificationID dao.PrimaryKey, Quantity uint) error {
	//Orm := singleton.Orm()
	shoppingCarts := m.ShoppingCart.FindShoppingCartByUserID(UserID)

	tx := db.Orm().Begin()

	//var goods model.Goods
	g := dao.GetByPrimaryKey(tx, entity.Goods, GoodsID).(*model.Goods)
	if g.IsZero() {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	//var specification model.Specification
	specification := dao.GetByPrimaryKey(tx, entity.Specification, SpecificationID).(*model.Specification)
	if specification.IsZero() {
		tx.Rollback()
		return gorm.ErrRecordNotFound
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
		shoppingCart := value.(*model.ShoppingCart)
		//var mgoods model.Goods
		//var mspecification model.Specification
		//util.JSONToStruct(shoppingCart.Goods, &mgoods)
		//util.JSONToStruct(shoppingCart.Specification, &mspecification)

		if shoppingCart.GoodsID == g.ID && shoppingCart.SpecificationID == specification.ID {
			//已经存在，添加数量
			shoppingCart.Quantity = shoppingCart.Quantity + Quantity
			if shoppingCart.Quantity > specification.Stock {
				shoppingCart.Quantity = specification.Stock
			}
			err := dao.UpdateByPrimaryKey(tx, entity.ShoppingCart, shoppingCart.ID, shoppingCart)
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
		//sc.Specification = util.StructToJSON(specification)
		//sc.Goods = util.StructToJSON(goods)
		//sc.GSID = strconv.Itoa(int(goods.ID)) + strconv.Itoa(int(specification.ID))
		sc.GoodsID = g.ID
		sc.SpecificationID = specification.ID
		err := dao.Create(tx, &sc)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil

}
func (m OrdersService) GetOrdersGoodsByOrdersGoodsNo(tx *gorm.DB, ordersGoodsNo string) *model.OrdersGoods {
	//var ogs []model.OrdersGoods
	//err := service.FindWhere(DB, &ogs, &model.OrdersGoods{OrdersID: OrdersID})
	ogs := dao.GetBy(tx, &model.OrdersGoods{}, map[string]any{"OrdersGoodsNo": ordersGoodsNo}).(*model.OrdersGoods)
	return ogs
}
func (m OrdersService) FindOrdersGoodsByOrdersID(DB *gorm.DB, OrdersID dao.PrimaryKey) ([]*model.OrdersGoods, error) {
	//var ogs []model.OrdersGoods
	//err := service.FindWhere(DB, &ogs, &model.OrdersGoods{OrdersID: OrdersID})
	//ogs := dao.Find(DB, &model.OrdersGoods{}).Where(`"OrdersID"=?`, OrdersID).List()
	var ordersGoodsList []*model.OrdersGoods
	DB.Model(&model.OrdersGoods{}).Where(`"OrdersID"=?`, OrdersID).Find(&ordersGoodsList)
	return ordersGoodsList, nil
}
func (m OrdersService) FindOrdersGoodsByCollageUser(CollageNo string) []model.User {
	orm := db.Orm()
	var user []model.User

	var sql = `SELECT u.* FROM "Orders" o,"OrdersGoods" og,"User" u WHERE og."CollageNo"=1 AND o."IsPay"=1 and o."ID"=og."OrdersID" AND u."ID"=o."UserID"`
	orm.Raw(sql, CollageNo).Scan(&user)
	//orm.Exec("SELECT u.* FROM Orders o,OrdersGoods og,USER u WHERE og.CollageNo=? AND o.ID=og.OrdersID AND u.ID=o.UserID", CollageNo).Find(&user)
	return user
}

func (m OrdersService) QueryOrdersTask(context constrain.IServiceContext, orders *model.Orders) error {
	pm := payment.NewPayment(context, orders.OID, orders.PayMethod)
	//if orders.IsPay == 0 {
	//当前状态为没有支付，去检测一下，订单状态。
	transaction, err := pm.OrderQuery(orders)
	if err != nil {
		return err
	}

	/*
		【交易状态】 交易状态，枚举值：
		* SUCCESS：支付成功
		* REFUND：转入退款
		* NOTPAY：未支付
		* CLOSED：已关闭
		* REVOKED：已撤销（仅付款码支付会返回）
		* USERPAYING：用户支付中（仅付款码支付会返回）
		* PAYERROR：支付失败（仅付款码支付会返回）
	*/
	switch transaction.State {
	case serviceargument.OrderQueryState_SUCCESS:
		_, err = m.OrderPaySuccess(uint(transaction.PayerTotalAmount), transaction.OutTradeNo, transaction.TransactionID, transaction.PayTime, orders.Type)
		if err != nil {
			return err
		}
	case serviceargument.OrderQueryState_REFUND:
		err = m.OrdersRefundSuccess(orders)
		if err != nil {
			return err
		}
	case serviceargument.OrderQueryState_NOTPAY:
	case serviceargument.OrderQueryState_CLOSED:
		err = dao.UpdateByPrimaryKey(db.Orm(), entity.Orders, orders.ID, map[string]interface{}{"Status": model.OrdersStatusClosed})
		if err != nil {
			return err
		}
	case serviceargument.OrderQueryState_REVOKED:
		err = dao.UpdateByPrimaryKey(db.Orm(), entity.Orders, orders.ID, map[string]interface{}{"Status": model.OrdersStatusClosed})
		if err != nil {
			return err
		}
	case serviceargument.OrderQueryState_USERPAYING:
	case serviceargument.OrderQueryState_PAYERROR:
	}
	//}
	err = m.AnalysisOrdersStatus(context, orders.ID, orders.OID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Cancel 申请取消
func (m OrdersService) Cancel(ctx constrain.IContext, OrdersID dao.PrimaryKey) (string, error) {
	Orm := db.Orm()

	//var orders model.Orders
	orders := dao.GetByPrimaryKey(Orm, entity.Orders, OrdersID).(*model.Orders)
	if orders.ID == 0 {

		return "", errors.New("订单不存在")
	}

	//下单状态
	if orders.Status == model.OrdersStatusOrder {
		if orders.IsPay == model.OrdersIsPayPayed {
			err := dao.UpdateByPrimaryKey(Orm, entity.Orders, OrdersID, map[string]interface{}{"Status": model.OrdersStatusCancel})
			return "申请取消，等待客服确认", err
		} else {
			//没支付的订单
			//管理商品库存
			err := m.OrdersStockManager(Orm, orders, false)
			if err != nil {
				return "", err
			}
			err = dao.UpdateByPrimaryKey(Orm, entity.Orders, OrdersID, map[string]interface{}{"Status": model.OrdersStatusCancelOk})
			return "取消成功", err

		}
	} else if orders.Status == model.OrdersStatusPay {
		if orders.IsPay == model.OrdersIsPayPayed {
			err := dao.UpdateByPrimaryKey(Orm, entity.Orders, OrdersID, map[string]interface{}{"Status": model.OrdersStatusCancel})
			if err != nil {
				return "", err
			}
			//已经支付的订单，发起退款
			_, err = m.CancelOk(ctx, orders.ID)
			if err != nil {
				return "", err
			}
			return "订单退款申请成功，退款资金已经按原路退回，请注意查收信息", nil
		} else {
			return "", errors.New("不允许取消订单,订单没有支付或已经过期")
		}
	} else {
		return "", errors.New("申请取消订单失败，请联系客服")
	}
}
