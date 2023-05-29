package user

import (
	"errors"
	"github.com/nbvghost/dandelion/library/db"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/internal/repository"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/dandelion/service/configuration"
	"github.com/nbvghost/dandelion/service/journal"

	"gorm.io/gorm"

	"github.com/nbvghost/gpa/types"
	"github.com/nbvghost/tool/encryption"

	"github.com/nbvghost/gweb"
)

type UserService struct {
	model.BaseDao
	Configuration configuration.ConfigurationService
	//GiveVoucher   activity.GiveVoucherService
	//CardItem      activity.CardItemService
	Organization company.OrganizationService
	//Wx           wechat.WxService
	//Goods        goods.GoodsService
	//Orders       order.OrdersService
	Journal journal.JournalService
}

func (service UserService) Login(account string) (user *model.User) {
	var err error
	if user, err = repository.User.GetByPhone(account); user.IsZero() {

		user, err = repository.User.GetByEmail(account)
		if err != nil {
			log.Println(err)
		}
	}
	return user

}
func (service UserService) UpdateLoginStatus(userID types.PrimaryKey) error {

	return repository.User.UpdateByID(userID, map[string]interface{}{"LastLoginAt": time.Now()}).Err

}
func (service UserService) AddUser(name, email, password string) error {

	hasUser, err := repository.User.GetByEmail(email)
	if hasUser.IsZero() == false {
		return errors.New("record is exist")
	}
	user := &model.User{Name: name, Email: email, Password: encryption.Md5ByString(password)}

	err = repository.User.Create(user)
	return err

}
func (service UserService) Situation(StartTime, EndTime int64) interface{} {

	st := time.Unix(StartTime/1000, 0)
	st = time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
	et := time.Unix(EndTime/1000, 0).Add(24 * time.Hour)
	et = time.Date(et.Year(), et.Month(), et.Day(), 0, 0, 0, 0, et.Location())

	Orm := db.Orm()

	type Result struct {
		TotalCount  uint `gorm:"column:TotalCount"`
		OnlineCount int
	}

	var result Result

	Orm.Table("User").Select(`COUNT(ID) as "TotalCount"`).Where(`"CreatedAt">=?`, st).Where(`"CreatedAt"<?`, et).Find(&result)
	//fmt.Println(result)
	//todo://
	//result.OnlineCount = len(gweb.Sessions.Data)
	return result
}
func (service UserService) AddUserBlockAmount(Orm *gorm.DB, UserID types.PrimaryKey, Menoy int64) error {

	user := dao.GetByPrimaryKey(Orm, &model.User{}, UserID).(*model.User)
	if user.IsZero() {
		return gorm.ErrRecordNotFound
	}

	tm := int64(user.BlockAmount) + Menoy
	if tm < 0 {
		return errors.New("冻结金额不足，无法扣款")
	}

	err := dao.UpdateByPrimaryKey(Orm, &model.User{}, UserID, map[string]interface{}{"BlockAmount": tm})
	return err
}

func (service UserService) FindUserByIDs(IDs []uint) []model.User {
	var users []model.User
	if len(IDs) == 0 {
		return users
	}
	err := db.Orm().Where(IDs).Find(&users).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	log.Println(err)
	return users
}
func (service UserService) LeveAll6(Orm *gorm.DB, OneSuperiorID types.PrimaryKey) string {
	//Orm := singleton.Orm()
	var leveIDs []string

	var user1 model.User
	Orm.Model(&model.User{}).Where(`"ID"=?`, OneSuperiorID).First(&user1)
	leveIDs = append(leveIDs, strconv.Itoa(int(user1.ID)))

	var user2 model.User
	Orm.Model(&model.User{}).Where(`"ID"=?`, user1.SuperiorID).First(&user2)
	leveIDs = append(leveIDs, strconv.Itoa(int(user2.ID)))

	var user3 model.User
	Orm.Model(&model.User{}).Where(`"ID"=?`, user2.SuperiorID).First(&user3)
	leveIDs = append(leveIDs, strconv.Itoa(int(user3.ID)))

	var user4 model.User
	Orm.Model(&model.User{}).Where(`"ID"=?`, user3.SuperiorID).First(&user4)
	leveIDs = append(leveIDs, strconv.Itoa(int(user4.ID)))

	var user5 model.User
	Orm.Model(&model.User{}).Where(`"ID"=?`, user4.SuperiorID).First(&user5)
	leveIDs = append(leveIDs, strconv.Itoa(int(user5.ID)))

	var user6 model.User
	Orm.Model(&model.User{}).Where(`"ID"=?`, user5.SuperiorID).First(&user6)
	leveIDs = append(leveIDs, strconv.Itoa(int(user6.ID)))

	return strings.Join(leveIDs, ",")
}
func (service UserService) Leve1(UserID types.PrimaryKey) []uint {
	Orm := db.Orm()
	var levea []uint
	if UserID <= 0 {
		return levea
	}
	Orm.Model(&model.User{}).Where(`"SuperiorID"=?`, UserID).Pluck(`"ID"`, &levea)
	return levea
}
func (service UserService) Leve2(Leve1IDs []uint) []uint {
	Orm := db.Orm()
	var levea []uint
	if len(Leve1IDs) <= 0 {
		return levea
	}
	Orm.Model(&model.User{}).Where(`"SuperiorID" in (?)`, Leve1IDs).Pluck(`"ID"`, &levea)
	return levea
}
func (service UserService) Leve3(Leve2IDs []uint) []uint {
	Orm := db.Orm()
	var levea []uint
	if len(Leve2IDs) <= 0 {
		return levea
	}
	Orm.Model(&model.User{}).Where(`"SuperiorID" in (?)`, Leve2IDs).Pluck(`"ID"`, &levea)
	return levea
}
func (service UserService) Leve4(Leve3IDs []uint) []uint {
	Orm := db.Orm()
	var levea []uint
	if len(Leve3IDs) <= 0 {
		return levea
	}
	Orm.Model(&model.User{}).Where(`"SuperiorID" in (?)`, Leve3IDs).Pluck(`"ID"`, &levea)
	return levea
}
func (service UserService) Leve5(Leve4IDs []uint) []uint {
	Orm := db.Orm()
	var levea []uint
	if len(Leve4IDs) <= 0 {
		return levea
	}
	Orm.Model(&model.User{}).Where(`"SuperiorID" in (?)`, Leve4IDs).Pluck(`"ID"`, &levea)
	return levea
}
func (service UserService) Leve6(Leve5IDs []uint) []uint {
	Orm := db.Orm()
	var levea []uint
	if len(Leve5IDs) <= 0 {
		return levea
	}
	Orm.Model(&model.User{}).Where(`"SuperiorID" in (?)`, Leve5IDs).Pluck(`"ID"`, &levea)
	return levea
}
func (service UserService) GetUserInfo(UserID types.PrimaryKey) model.UserInfo {
	Orm := db.Orm()
	//.First(&user, 10)
	var userInfo model.UserInfo
	Orm.Where(&model.UserInfo{UserID: UserID}).First(&userInfo)
	if userInfo.ID == 0 && UserID != 0 {
		userInfo.UserID = UserID
		dao.Create(Orm, &userInfo)
	}
	return userInfo
}

func (service UserService) UserAction(context *gweb.Context) (r gweb.Result, err error) {
	company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	Orm := db.Orm()
	action := context.Request.URL.Query().Get("action")
	switch action {
	case "list":
		dts := &model.Datatables{}
		util.RequestBodyToJSON(context.Request.Body, dts)
		draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]model.User{}, company.ID, "")
		return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, nil
	}

	return &gweb.JsonResult{Data: result.ActionResult{Code: result.Fail, Message: "", Data: nil}}, nil
}
func (service UserService) FindUserByTel(Orm *gorm.DB, Tel string) *model.User {
	user := &model.User{}
	err := Orm.Where(`"Phone"=?`, Tel).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	log.Println(err)
	return user
}

func (service UserService) FindUserByOpenID(Orm *gorm.DB, OID types.PrimaryKey, OpenID string) *model.User {
	user := &model.User{}
	//CompanyOpenID := user.GetCompanyOpenID(CompanyID, OpenID)
	err := Orm.Where(`"OpenID"=? and "OID"=?`, OpenID, OID).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	log.Println(err)
	return user
}
func (service UserService) AddUserByOpenID(Orm *gorm.DB, OID types.PrimaryKey, OpenID string) *model.User {
	//Orm := singleton.Orm()

	user := service.FindUserByOpenID(Orm, OID, OpenID)
	if user.ID == 0 {
		user.OID = OID
		user.OpenID = OpenID
		dao.Create(Orm, user)
	} else {

	}
	//CompanyOpenID := user.GetCompanyOpenID(CompanyID, OpenID)
	//err := Orm.Where("OpenID=?", OpenID).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	//log.Println(err)
	return user
}
