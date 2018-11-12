package service

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"dandelion/app/service/dao"
	"dandelion/app/util"

	"dandelion/app/play"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type UserService struct {
	dao.BaseDao
	Configuration ConfigurationService
	GiveVoucher   GiveVoucherService
	CardItem      CardItemService
	Organization  OrganizationService
	Wx            WxService
	Goods         GoodsService
}

func (service UserService) Situation(StartTime, EndTime int64) interface{} {

	st := time.Unix(StartTime/1000, 0)
	st = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
	et := time.Unix(EndTime/1000, 0).Add(24 * time.Hour)
	et = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())

	Orm := dao.Orm()

	type Result struct {
		TotalCount  uint64 `gorm:"column:TotalCount"`
		OnlineCount int
	}

	var result Result

	Orm.Table("User").Select("COUNT(ID) as TotalCount").Where("CreatedAt>=?", st).Where("CreatedAt<?", et).Find(&result)
	//fmt.Println(result)
	result.OnlineCount = len(gweb.Sessions.Data)
	return result
}
func (service UserService) GetFromIDs(UserID uint64) dao.UserFormIds {
	var result dao.UserFormIds
	for {
		dao.Orm().Table("UserFormIds").Where("UserID=?", UserID).Order("CreatedAt asc").First(&result)
		if result.ID == 0 {
			break
		}
		if result.CreatedAt.Add(7*24*time.Hour).Unix() < time.Now().Unix() {
			service.Delete(dao.Orm(), &dao.UserFormIds{}, result.ID)
		} else {
			break
		}
	}

	return result
}
func (service UserService) AddUserBlockAmount(Orm *gorm.DB, UserID uint64, Menoy int64) error {

	var user dao.User
	err := service.Get(Orm, UserID, &user)
	if err != nil {
		return err
	}

	tm := int64(user.BlockAmount) + Menoy
	if tm < 0 {
		return errors.New("冻结金额不足，无法扣款")
	}

	err = service.ChangeMap(Orm, UserID, &dao.User{}, map[string]interface{}{"BlockAmount": tm})
	return err
}
func (service UserService) MinusSettlementUserBrokerage(tx *gorm.DB, orders dao.Orders) error {
	var err error
	//用户自己。下单者
	//Orm:=dao.Orm()

	ogs, err := GlobalService.Orders.FindOrdersGoodsByOrdersID(tx, orders.ID)
	var Brokerage uint64
	for _, value := range ogs {
		//var specification dao.Specification
		//util.JSONToStruct(value.Specification, &specification)
		Brokerage = Brokerage + value.TotalBrokerage
	}

	//var orders dao.Orders
	//service.Get(Orm, OrderID, &orders)

	var user dao.User
	service.Get(tx, orders.UserID, &user)

	leve1, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve1).V, 10, 64)
	leve2, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve2).V, 10, 64)
	leve3, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve3).V, 10, 64)
	leve4, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve4).V, 10, 64)
	leve5, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve5).V, 10, 64)
	leve6, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve6).V, 10, 64)

	leves := []uint64{leve1, leve2, leve3, leve4, leve5, leve6}

	//var OutBrokerageMoney int64 = 0
	for _, value := range leves {
		if value <= 0 {
			break
		}
		var _user dao.User
		service.Get(tx, user.SuperiorID, &_user)
		if _user.ID <= 0 {
			return nil
		}
		leveMenoy := int64(math.Floor(float64(value)/float64(100)*float64(Brokerage) + 0.5))
		err = service.AddUserBlockAmount(tx, _user.ID, -leveMenoy)
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

//如果订单未完成，或是退款，扣除相应的冻结金额，不用结算，佣金
func (service UserService) AfterSettlementUserBrokerage(tx *gorm.DB, orders dao.Orders) error {
	var err error
	//用户自己。下单者
	//Orm:=dao.Orm()

	//var orders dao.Orders
	//service.Get(Orm, OrderID, &orders)

	ogs, err := GlobalService.Orders.FindOrdersGoodsByOrdersID(tx, orders.ID)
	var Brokerage uint64
	for _, value := range ogs {
		//var specification dao.Specification
		//util.JSONToStruct(value.Specification, &specification)
		Brokerage = Brokerage + value.TotalBrokerage
	}

	var user dao.User
	service.Get(tx, orders.UserID, &user)

	leve1, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve1).V, 10, 64)
	leve2, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve2).V, 10, 64)
	leve3, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve3).V, 10, 64)
	leve4, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve4).V, 10, 64)
	leve5, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve5).V, 10, 64)
	leve6, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve6).V, 10, 64)

	leves := []uint64{leve1, leve2, leve3, leve4, leve5, leve6}

	//var OutBrokerageMoney int64 = 0
	for _, value := range leves {
		if value <= 0 {
			break
		}
		var _user dao.User
		service.Get(tx, user.SuperiorID, &_user)
		if _user.ID <= 0 {
			return nil
		}
		leveMenoy := int64(math.Floor(float64(value)/float64(100)*float64(Brokerage) + 0.5))
		err = service.AddUserBlockAmount(tx, _user.ID, -leveMenoy)
		if err != nil {
			tool.CheckError(err)
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
func (service UserService) FirstSettlementUserBrokerage(tx *gorm.DB, orders dao.Orders) error {
	var err error
	//用户自己。下单者
	//Orm:=dao.Orm()

	//var orders dao.Orders
	//service.Get(Orm, OrderID, &orders)

	ogs, err := GlobalService.Orders.FindOrdersGoodsByOrdersID(tx, orders.ID)
	var Brokerage uint64
	for _, value := range ogs {
		//var specification dao.Specification
		//util.JSONToStruct(value.Specification, &specification)
		Brokerage = Brokerage + value.TotalBrokerage
	}

	var user dao.User
	service.Get(tx, orders.UserID, &user)

	leve1, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve1).V, 10, 64)
	leve2, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve2).V, 10, 64)
	leve3, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve3).V, 10, 64)
	leve4, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve4).V, 10, 64)
	leve5, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve5).V, 10, 64)
	leve6, _ := strconv.ParseUint(service.Configuration.GetConfiguration(orders.OID, play.ConfigurationKey_BrokerageLeve6).V, 10, 64)

	leves := []uint64{leve1, leve2, leve3, leve4, leve5, leve6}

	//var OutBrokerageMoney int64 = 0
	for index, value := range leves {
		if value <= 0 {
			break
		}
		var _user dao.User
		service.Get(tx, user.SuperiorID, &_user)
		if _user.ID <= 0 {
			return nil
		}
		leveMenoy := int64(math.Floor(float64(value)/float64(100)*float64(Brokerage) + 0.5))
		err = service.AddUserBlockAmount(tx, _user.ID, leveMenoy)
		if err != nil {
			tool.CheckError(err)
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

//结算佣金，结算积分，结算成长值，是否送福利卷
func (service UserService) SettlementUser(Orm *gorm.DB, Brokerage uint64, orders dao.Orders) error {
	var err error
	//用户自己。下单者

	//var orders dao.Orders
	//service.Get(Orm, OrderID, &orders)

	var user dao.User
	service.Get(Orm, orders.UserID, &user)

	//fmt.Println(user.Name)

	Journal := JournalService{}
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
	err = Journal.AddScoreJournal(Orm, user.ID, "积分", "购买商品", play.ScoreJournal_Type_GM, int64(user.Score), dao.KV{Key: "OrdersID", Value: orders.ID})
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

	err = Journal.AddOrganizationJournal(Orm, orders.OID, "商品交易", "商品交易", play.OrganizationJournal_Goods, int64(orders.PayMoney), dao.KV{Key: "OrdersID", Value: orders.ID})

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
		err = Journal.AddUserJournal(Orm, _user.ID, "佣金", strconv.Itoa(index+1)+"级用户", play.UserJournal_Type_LEVE, leveMenoy, dao.KV{Key: "OrdersID", Value: orders.ID}, user.ID)
		if err != nil {
			tool.CheckError(err)
			continue
		}

		err = service.AddUserBlockAmount(Orm, _user.ID, -leveMenoy)
		if err != nil {
			tool.CheckError(err)
			continue
		}

		err = Journal.AddOrganizationJournal(Orm, orders.OID, "商品交易", "推广佣金"+_user.Name, play.OrganizationJournal_Brokerage, -leveMenoy, dao.KV{Key: "OrdersID", Value: orders.ID})
		if err != nil {
			tool.CheckError(err)
			continue
		}

		err = Journal.AddScoreJournal(Orm, _user.ID, "积分", "佣金积分", play.ScoreJournal_Type_LEVE, int64(leveMenoy), dao.KV{Key: "OrdersID", Value: orders.ID})
		if err != nil {
			tool.CheckError(err)
			continue
		}

		workTime := time.Now().Unix() - orders.CreatedAt.Unix()

		service.Wx.INComeNotify(_user, "来自"+strconv.Itoa(index+1)+"级用户，现金收入", strconv.Itoa(int(workTime/60/60))+"小时", "收入："+strconv.FormatFloat(float64(leveMenoy)/float64(100), 'f', 2, 64)+"元")

		user = _user
	}

	return nil
}
func (service UserService) FindUserByIDs(IDs []uint64) []dao.User {
	var users []dao.User
	err := dao.Orm().Where(IDs).Find(&users).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	tool.CheckError(err)
	return users
}
func (service UserService) LeveAll6(Orm *gorm.DB, OneSuperiorID uint64) string {
	//Orm := dao.Orm()
	var leveIDs []string

	var user1 dao.User
	Orm.Model(&dao.User{}).Where("ID=?", OneSuperiorID).First(&user1)
	leveIDs = append(leveIDs, strconv.Itoa(int(user1.ID)))

	var user2 dao.User
	Orm.Model(&dao.User{}).Where("ID=?", user1.SuperiorID).First(&user2)
	leveIDs = append(leveIDs, strconv.Itoa(int(user2.ID)))

	var user3 dao.User
	Orm.Model(&dao.User{}).Where("ID=?", user2.SuperiorID).First(&user3)
	leveIDs = append(leveIDs, strconv.Itoa(int(user3.ID)))

	var user4 dao.User
	Orm.Model(&dao.User{}).Where("ID=?", user3.SuperiorID).First(&user4)
	leveIDs = append(leveIDs, strconv.Itoa(int(user4.ID)))

	var user5 dao.User
	Orm.Model(&dao.User{}).Where("ID=?", user4.SuperiorID).First(&user5)
	leveIDs = append(leveIDs, strconv.Itoa(int(user5.ID)))

	var user6 dao.User
	Orm.Model(&dao.User{}).Where("ID=?", user5.SuperiorID).First(&user6)
	leveIDs = append(leveIDs, strconv.Itoa(int(user6.ID)))

	return strings.Join(leveIDs, ",")
}
func (service UserService) Leve1(UserID uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if UserID <= 0 {
		return levea
	}
	Orm.Model(&dao.User{}).Where("SuperiorID=?", UserID).Pluck("ID", &levea)
	return levea
}
func (service UserService) Leve2(Leve1IDs []uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if len(Leve1IDs) <= 0 {
		return levea
	}
	Orm.Model(&dao.User{}).Where("SuperiorID in (?)", Leve1IDs).Pluck("ID", &levea)
	return levea
}
func (service UserService) Leve3(Leve2IDs []uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if len(Leve2IDs) <= 0 {
		return levea
	}
	Orm.Model(&dao.User{}).Where("SuperiorID in (?)", Leve2IDs).Pluck("ID", &levea)
	return levea
}
func (service UserService) Leve4(Leve3IDs []uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if len(Leve3IDs) <= 0 {
		return levea
	}
	Orm.Model(&dao.User{}).Where("SuperiorID in (?)", Leve3IDs).Pluck("ID", &levea)
	return levea
}
func (service UserService) Leve5(Leve4IDs []uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if len(Leve4IDs) <= 0 {
		return levea
	}
	Orm.Model(&dao.User{}).Where("SuperiorID in (?)", Leve4IDs).Pluck("ID", &levea)
	return levea
}
func (service UserService) Leve6(Leve5IDs []uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if len(Leve5IDs) <= 0 {
		return levea
	}
	Orm.Model(&dao.User{}).Where("SuperiorID in (?)", Leve5IDs).Pluck("ID", &levea)
	return levea
}
func (service UserService) GetUserInfo(UserID uint64) dao.UserInfo {
	Orm := dao.Orm()
	//.First(&user, 10)
	var userInfo dao.UserInfo
	Orm.Where(&dao.UserInfo{UserID: UserID}).First(&userInfo)
	if userInfo.ID == 0 && UserID != 0 {
		userInfo.UserID = UserID
		service.Add(Orm, &userInfo)
	}
	return userInfo
}

func (service UserService) UserAction(context *gweb.Context) gweb.Result {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	action := context.Request.URL.Query().Get("action")
	switch action {
	case "list":
		dts := &dao.Datatables{}
		util.RequestBodyToJSON(context.Request.Body, dts)
		draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.User{}, company.ID)
		return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
	}

	return &gweb.JsonResult{Data: dao.ActionStatus{Success: false, Message: "", Data: nil}}
}
func (service UserService) FindUserByTel(Orm *gorm.DB, Tel string) *dao.User {
	user := &dao.User{}
	err := Orm.Where("Tel=?", Tel).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	tool.CheckError(err)
	return user
}

func (service UserService) FindUserByOpenID(Orm *gorm.DB, OpenID string) *dao.User {

	user := &dao.User{}
	//CompanyOpenID := user.GetCompanyOpenID(CompanyID, OpenID)
	err := Orm.Where("OpenID=?", OpenID).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	tool.CheckError(err)
	return user
}
func (service UserService) AddUserByOpenID(Orm *gorm.DB, OpenID string) *dao.User {
	//Orm := dao.Orm()
	user := &dao.User{}
	user = service.FindUserByOpenID(Orm, OpenID)
	if user.ID == 0 {
		user.OpenID = OpenID
		service.Add(Orm, user)
	} else {

	}
	//CompanyOpenID := user.GetCompanyOpenID(CompanyID, OpenID)
	//err := Orm.Where("OpenID=?", OpenID).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	//tool.CheckError(err)
	return user
}
