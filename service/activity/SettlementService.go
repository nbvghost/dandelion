package activity

import (
	"log"
	"math"
	"strconv"
	"time"

	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/service/configuration"
	"github.com/nbvghost/dandelion/service/journal"
	"github.com/nbvghost/dandelion/service/user"
	"github.com/nbvghost/dandelion/service/wechat"

	"gorm.io/gorm"

	"github.com/nbvghost/tool/object"
)

type SettlementService struct {
	model.BaseDao
	Configuration configuration.ConfigurationService
	Journal       journal.JournalService
	GiveVoucher   GiveVoucherService
	CardItem      CardItemService
	Wx            wechat.WxService
	MessageNotify wechat.MessageNotify
	User          user.UserService
}

// 结算佣金，结算积分，结算成长值，是否送福利卷
func (service SettlementService) SettlementUser(Orm *gorm.DB, ordersGoodsList []*model.OrdersGoods, orders *model.Orders) error {
	var err error
	//用户自己。下单者

	//var orders model.Orders
	//service.Get(Orm, OrderID, &orders)

	//var user model.User
	//service.Get(Orm, orders.UserID, &user)
	u := dao.GetByPrimaryKey(Orm, &model.User{}, orders.UserID).(*model.User)
	//fmt.Println(user.Name)

	brokerage := service.Configuration.GetBrokerageConfiguration(orders.OID) //service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageType)

	var Brokerage uint
	for i := range ordersGoodsList {
		value := ordersGoodsList[i]
		//var specification model.Specification
		//util.JSONToStruct(value.Specification, &specification)

		if brokerage.Type == configuration.BrokeragePRODUCT {
			Brokerage = Brokerage + value.SellPrice
		}
		if brokerage.Type == configuration.BrokerageCUSTOM {
			Brokerage = Brokerage + value.TotalBrokerage
		}
	}

	leves := []uint{brokerage.Leve1, brokerage.Leve2, brokerage.Leve3, brokerage.Leve4, brokerage.Leve5, brokerage.Leve6}

	GrowValue := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyScoreConvertGrowValue).V)

	u.Score = u.Score + orders.PayMoney
	u.Growth = u.Growth + (uint(math.Floor(float64(orders.PayMoney)/100+0.5)) * GrowValue)
	err = dao.UpdateByPrimaryKey(Orm, &model.User{}, u.ID, &model.User{Growth: u.Growth})
	if err != nil {
		return err
	}
	err = service.Journal.AddScoreJournal(Orm, u.ID, "积分", "购买商品", play.ScoreJournal_Type_GM, int64(u.Score), extends.KV{Key: "OrdersID", Value: orders.ID})
	if err != nil {
		return err
	}

	gvs := service.GiveVoucher.FindASC()
	for _, value := range gvs {
		//主订单的金额来决定是否送卡卷
		if uint(orders.PayMoney) >= value.ScoreMaxValue {

			err := service.CardItem.AddVoucherCardItem(Orm, orders.OrderNo, orders.UserID, value.VoucherID)
			if err != nil {
				return err
			}
			break
		}
	}

	err = service.Journal.AddOrganizationJournal(Orm, orders.OID, "商品交易", "商品交易", play.OrganizationJournal_Goods, int64(orders.PayMoney), extends.KV{Key: "OrdersID", Value: orders.ID})

	if err != nil {
		return err
	}
	for index, value := range leves {
		if value <= 0 {
			break
		}
		//var _user model.User
		_user := dao.GetByPrimaryKey(Orm, &model.User{}, u.SuperiorID).(*model.User)
		if _user.ID <= 0 {
			return nil
		}

		leveMenoy := int64(math.Floor(float64(value)/float64(100)*float64(Brokerage) + 0.5))
		err = service.Journal.AddUserJournal(Orm, _user.ID, "佣金", strconv.Itoa(index+1)+"级用户", play.UserJournal_Type_LEVE, leveMenoy, extends.KV{Key: "OrdersID", Value: orders.ID}, u.ID)
		if err != nil {
			log.Println(err)
			continue
		}

		err = service.User.AddUserBlockAmount(Orm, _user.ID, -leveMenoy)
		if err != nil {
			log.Println(err)
			continue
		}

		err = service.Journal.AddOrganizationJournal(Orm, orders.OID, "商品交易", "推广佣金"+_user.Name, play.OrganizationJournal_Brokerage, -leveMenoy, extends.KV{Key: "OrdersID", Value: orders.ID})
		if err != nil {
			log.Println(err)
			continue
		}

		err = service.Journal.AddScoreJournal(Orm, _user.ID, "积分", "佣金积分", play.ScoreJournal_Type_LEVE, int64(leveMenoy), extends.KV{Key: "OrdersID", Value: orders.ID})
		if err != nil {
			log.Println(err)
			continue
		}

		workTime := time.Now().Unix() - orders.CreatedAt.Unix()

		service.MessageNotify.INComeNotify(_user, "来自"+strconv.Itoa(index+1)+"级用户，现金收入", strconv.Itoa(int(workTime/60/60))+"小时", "收入："+strconv.FormatFloat(float64(leveMenoy)/float64(100), 'f', 2, 64)+"元")

		u = _user
	}

	return nil
}
