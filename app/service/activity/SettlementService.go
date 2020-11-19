package activity

import (
	"github.com/jinzhu/gorm"
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service/configuration"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/service/journal"
	"github.com/nbvghost/dandelion/app/service/user"
	"github.com/nbvghost/dandelion/app/service/wechat"
	"github.com/nbvghost/glog"
	"math"
	"strconv"
	"time"
)

type SettlementService struct {
	dao.BaseDao
	Configuration configuration.ConfigurationService
	Journal       journal.JournalService
	GiveVoucher   GiveVoucherService
	CardItem      CardItemService
	Wx            wechat.WxService
	User          user.UserService
}

//结算佣金，结算积分，结算成长值，是否送福利卷
func (service SettlementService) SettlementUser(Orm *gorm.DB, Brokerage uint64, orders dao.Orders) error {
	var err error
	//用户自己。下单者

	//var orders dao.Orders
	//service.Get(Orm, OrderID, &orders)

	var user dao.User
	service.Get(Orm, orders.UserID, &user)

	//fmt.Println(user.Name)

	leve1, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve1).V, 10, 64)
	leve2, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve2).V, 10, 64)
	leve3, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve3).V, 10, 64)
	leve4, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve4).V, 10, 64)
	leve5, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve5).V, 10, 64)
	leve6, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve6).V, 10, 64)

	leves := []uint64{leve1, leve2, leve3, leve4, leve5, leve6}

	GrowValue, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_ScoreConvertGrowValue).V, 10, 64)

	user.Score = user.Score + orders.PayMoney
	user.Growth = user.Growth + uint64(uint64(math.Floor(float64(orders.PayMoney)/100+0.5))*GrowValue)
	err = service.ChangeModel(Orm, user.ID, &dao.User{Growth: user.Growth})
	if err != nil {
		return err
	}
	err = service.Journal.AddScoreJournal(Orm, user.ID, "积分", "购买商品", play.ScoreJournal_Type_GM, int64(user.Score), dao.KV{Key: "OrdersID", Value: orders.ID})
	if err != nil {
		return err
	}

	gvs := service.GiveVoucher.FindASC()
	for _, value := range gvs {
		//主订单的金额来决定是否送卡卷
		if uint64(orders.PayMoney) >= value.ScoreMaxValue {

			err := service.CardItem.AddVoucherCardItem(Orm, orders.OrderNo, orders.UserID, value.VoucherID)
			if err != nil {
				return err
			}
			break
		}
	}

	err = service.Journal.AddOrganizationJournal(Orm, orders.OID, "商品交易", "商品交易", play.OrganizationJournal_Goods, int64(orders.PayMoney), dao.KV{Key: "OrdersID", Value: orders.ID})

	if err != nil {
		return err
	}
	for index, value := range leves {
		if value <= 0 {
			break
		}
		var _user dao.User
		service.Get(Orm, user.SuperiorID, &_user)
		if _user.ID <= 0 {
			return nil
		}

		leveMenoy := int64(math.Floor(float64(value)/float64(100)*float64(Brokerage) + 0.5))
		err = service.Journal.AddUserJournal(Orm, _user.ID, "佣金", strconv.Itoa(index+1)+"级用户", play.UserJournal_Type_LEVE, leveMenoy, dao.KV{Key: "OrdersID", Value: orders.ID}, user.ID)
		if err != nil {
			glog.Error(err)
			continue
		}

		err = service.User.AddUserBlockAmount(Orm, _user.ID, -leveMenoy)
		if err != nil {
			glog.Error(err)
			continue
		}

		err = service.Journal.AddOrganizationJournal(Orm, orders.OID, "商品交易", "推广佣金"+_user.Name, play.OrganizationJournal_Brokerage, -leveMenoy, dao.KV{Key: "OrdersID", Value: orders.ID})
		if err != nil {
			glog.Error(err)
			continue
		}

		err = service.Journal.AddScoreJournal(Orm, _user.ID, "积分", "佣金积分", play.ScoreJournal_Type_LEVE, int64(leveMenoy), dao.KV{Key: "OrdersID", Value: orders.ID})
		if err != nil {
			glog.Error(err)
			continue
		}

		workTime := time.Now().Unix() - orders.CreatedAt.Unix()

		service.Wx.INComeNotify(_user, "来自"+strconv.Itoa(index+1)+"级用户，现金收入", strconv.Itoa(int(workTime/60/60))+"小时", "收入："+strconv.FormatFloat(float64(leveMenoy)/float64(100), 'f', 2, 64)+"元")

		user = _user
	}

	return nil
}
