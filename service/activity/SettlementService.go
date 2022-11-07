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

//结算佣金，结算积分，结算成长值，是否送福利卷
func (service SettlementService) SettlementUser(Orm *gorm.DB, Brokerage uint, orders *model.Orders) error {
	var err error
	//用户自己。下单者

	//var orders model.Orders
	//service.Get(Orm, OrderID, &orders)

	//var user model.User
	//service.Get(Orm, orders.UserID, &user)
	user := dao.GetByPrimaryKey(Orm, &model.User{}, orders.UserID).(*model.User)
	//fmt.Println(user.Name)

	leve1 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageLeve1).V)
	leve2 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageLeve2).V)
	leve3 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageLeve3).V)
	leve4 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageLeve4).V)
	leve5 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageLeve5).V)
	leve6 := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyBrokerageLeve6).V)

	leves := []uint{leve1, leve2, leve3, leve4, leve5, leve6}

	GrowValue := object.ParseUint(service.Configuration.GetConfiguration(orders.OID, model.ConfigurationKeyScoreConvertGrowValue).V)

	user.Score = user.Score + orders.PayMoney
	user.Growth = user.Growth + (uint(math.Floor(float64(orders.PayMoney)/100+0.5)) * GrowValue)
	err = dao.UpdateByPrimaryKey(Orm, &model.User{}, user.ID, &model.User{Growth: user.Growth})
	if err != nil {
		return err
	}
	err = service.Journal.AddScoreJournal(Orm, user.ID, "积分", "购买商品", play.ScoreJournal_Type_GM, int64(user.Score), extends.KV{Key: "OrdersID", Value: orders.ID})
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
		_user := dao.GetByPrimaryKey(Orm, &model.User{}, user.SuperiorID).(*model.User)
		if _user.ID <= 0 {
			return nil
		}

		leveMenoy := int64(math.Floor(float64(value)/float64(100)*float64(Brokerage) + 0.5))
		err = service.Journal.AddUserJournal(Orm, _user.ID, "佣金", strconv.Itoa(index+1)+"级用户", play.UserJournal_Type_LEVE, leveMenoy, extends.KV{Key: "OrdersID", Value: orders.ID}, user.ID)
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

		user = _user
	}

	return nil
}
