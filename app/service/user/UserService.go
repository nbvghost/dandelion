package user

import (
	"errors"
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service/company"
	"github.com/nbvghost/dandelion/app/service/configuration"
	"github.com/nbvghost/dandelion/app/service/journal"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/util"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
)

type UserService struct {
	dao.BaseDao
	Configuration configuration.ConfigurationService
	//GiveVoucher   activity.GiveVoucherService
	//CardItem      activity.CardItemService
	Organization company.OrganizationService
	//Wx           wechat.WxService
	//Goods        goods.GoodsService
	//Orders       order.OrdersService
	Journal journal.JournalService
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
	//todo://
	//result.OnlineCount = len(gweb.Sessions.Data)
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

func (service UserService) FindUserByIDs(IDs []uint64) []dao.User {
	var users []dao.User
	err := dao.Orm().Where(IDs).Find(&users).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	glog.Error(err)
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
		draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.User{}, company.ID, "")
		return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
	}

	return &gweb.JsonResult{Data: result.ActionResult{Code: result.ActionFail, Message: "", Data: nil}}
}
func (service UserService) FindUserByTel(Orm *gorm.DB, Tel string) *dao.User {
	user := &dao.User{}
	err := Orm.Where("Tel=?", Tel).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	glog.Error(err)
	return user
}

func (service UserService) FindUserByOpenID(Orm *gorm.DB, OpenID string) *dao.User {

	user := &dao.User{}
	//CompanyOpenID := user.GetCompanyOpenID(CompanyID, OpenID)
	err := Orm.Where("OpenID=?", OpenID).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	glog.Error(err)
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
	//glog.Error(err)
	return user
}
