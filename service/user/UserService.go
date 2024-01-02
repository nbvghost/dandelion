package user

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/tool/object"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/dandelion/service/configuration"
	"github.com/nbvghost/dandelion/service/journal"

	"gorm.io/gorm"

	"github.com/nbvghost/tool/encryption"
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
	if user = service.GetByPhone(db.Orm(), account); user.IsZero() {
		user = service.GetByEmail(db.Orm(), account)
	}
	return user
}
func (service UserService) UpdateLoginStatus(userID dao.PrimaryKey) error {

	return dao.UpdateByPrimaryKey(db.Orm(), &model.User{}, userID, map[string]interface{}{"LastLoginAt": time.Now()}) //repository.User.UpdateByID(userID, map[string]interface{}{"LastLoginAt": time.Now()}).Err

}
func (service UserService) AddUser(name, email, password string) error {

	hasUser := service.GetByEmail(db.Orm(), email)
	if hasUser.IsZero() == false {
		return errors.New("record is exist")
	}
	user := &model.User{Name: name, Email: email, Password: encryption.Md5ByString(password)}

	err := dao.Create(db.Orm(), user) //repository.User.Create(user)
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

func (service UserService) FindUserByIDs(IDs []uint) []model.User {
	var users []model.User
	if len(IDs) == 0 {
		return users
	}
	err := db.Orm().Where(IDs).Find(&users).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	log.Println(err)
	return users
}
func (service UserService) LeveAll6(Orm *gorm.DB, OneSuperiorID dao.PrimaryKey) string {
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
func (service UserService) Leve1(UserID dao.PrimaryKey) []uint {
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
func (service UserService) GetUserInfo(UserID dao.PrimaryKey) *UserInfoValue {
	Orm := db.Orm()
	//.First(&user, 10)
	m := make(map[model.UserInfoKey]string)
	oldD := make(map[model.UserInfoKey]string)
	var userInfo []*model.UserInfo
	Orm.Where(`"UserID"=?`, UserID).Find(&userInfo)
	for _, v := range userInfo {
		m[v.Key] = v.Value
		oldD[v.Key] = v.Value
	}
	return &UserInfoValue{userID: UserID, data: m, oldData: oldD}
}

type UserInfoValue struct {
	userID  dao.PrimaryKey
	data    map[model.UserInfoKey]string
	oldData map[model.UserInfoKey]string
}

func (m *UserInfoValue) Data() map[model.UserInfoKey]string {
	return m.data
}
func (m *UserInfoValue) Update(db *gorm.DB) error {
	change := make(map[model.UserInfoKey]string)
	for key, s := range m.data {
		if v, ok := m.oldData[key]; ok {
			if strings.EqualFold(v, s) {
				continue
			}
		}
		change[key] = m.data[key]
	}
	if len(change) > 0 {
		for key, s := range change {
			has := dao.GetBy(db, &model.UserInfo{}, map[string]any{"UserID": m.userID, "Key": key})
			if has.IsZero() {
				err := dao.Create(db, &model.UserInfo{
					UserID: m.userID,
					Key:    key,
					Value:  s,
				})
				if err != nil {
					return err
				}
			} else {
				err := dao.UpdateBy(db, &model.UserInfo{}, map[string]any{"Value": s}, map[string]any{"UserID": m.userID, "Key": key})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func (m *UserInfoValue) GetDaySignTime() time.Time {
	if v, ok := m.data[model.UserInfoKeyDaySignTime]; ok {
		t, _ := time.Parse(time.RFC3339, v)
		return t
	}
	return time.Time{}
}
func (m *UserInfoValue) SetDaySignTime(v time.Time) {
	m.data[model.UserInfoKeyDaySignTime] = v.Format(time.RFC3339)
}

func (m *UserInfoValue) GetDaySignCount() int {
	return object.ParseInt(m.data[model.UserInfoKeyDaySignCount])
}

func (m *UserInfoValue) GetLastIP() string {
	return object.ParseString(m.data[model.UserInfoKeyLastIP])
}
func (m *UserInfoValue) GetAllowAssistance() bool {
	return object.ParseBool(m.data[model.UserInfoKeyAllowAssistance])
}
func (m *UserInfoValue) GetSubscribe() bool {
	return object.ParseBool(m.data[model.UserInfoKeySubscribe])
}
func (m *UserInfoValue) GetAgent() bool {
	return object.ParseBool(m.data[model.UserInfoKeyAgent])
}

func (m *UserInfoValue) SetDaySignCount(v int) {
	m.data[model.UserInfoKeyDaySignCount] = object.ParseString(v)
}
func (m *UserInfoValue) SetLastIP(v string) {
	m.data[model.UserInfoKeyLastIP] = v
}
func (m *UserInfoValue) SetAllowAssistance(v bool) {
	m.data[model.UserInfoKeyAllowAssistance] = object.ParseString(v)
}
func (m *UserInfoValue) SetSubscribe(v bool) {
	m.data[model.UserInfoKeySubscribe] = object.ParseString(v)
}

type UserInfoKeyStateType string

const (
	UserInfoKeyStateTypeNormal  UserInfoKeyStateType = ""        //
	UserInfoKeyStateTypeClosure UserInfoKeyStateType = "closure" //封闭
)

func (m *UserInfoValue) SetState(v UserInfoKeyStateType) {
	m.data[model.UserInfoKeyState] = string(v)
}
func (m *UserInfoValue) SetAgent(v bool) {
	m.data[model.UserInfoKeyAgent] = strconv.FormatBool(v)
}
func (m *UserInfoValue) GetState() UserInfoKeyStateType {
	return UserInfoKeyStateType(m.data[model.UserInfoKeyState])
}

func (service UserService) UserAction(context constrain.IContext) (r constrain.IResult, err error) {
	/*company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	Orm := db.Orm()
	action := context.Request.URL.Query().Get("action")
	switch action {
	case "list":
		dts := &model.Datatables{}
		util.RequestBodyToJSON(context.Request.Body, dts)
		draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]model.User{}, company.ID, "")
		return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, nil
	}*/

	//return &result.JsonResult{Data: result.ActionResult{Code: result.Fail, Message: "", Data: nil}}, nil
	return result.NewData(nil), nil
}
func (service UserService) GetByPhone(Orm *gorm.DB, Tel string) *model.User {
	user := &model.User{}
	Orm.Where(`"Phone"=?`, Tel).First(user) //SelectOne(user, "select * from User where Tel=?", Tel)
	return user
}

func (service UserService) GetByEmail(Orm *gorm.DB, email string) *model.User {
	user := &model.User{}
	Orm.Where(`"Email"=?`, email).First(user) //SelectOne(user, "select * from User where Tel=?", Tel)
	return user
}

func (service UserService) FindUserByOpenID(Orm *gorm.DB, OID dao.PrimaryKey, OpenID string) *model.User {
	user := &model.User{}
	//CompanyOpenID := user.GetCompanyOpenID(CompanyID, OpenID)
	Orm.Where(`"OpenID"=? and "OID"=?`, OpenID, OID).First(user) //SelectOne(user, "select * from User where Tel=?", Tel)
	return user
}
func (service UserService) AddUserByOpenID(Orm *gorm.DB, OID dao.PrimaryKey, OpenID string) (*model.User, error) {
	//Orm := singleton.Orm()
	user := service.FindUserByOpenID(Orm, OID, OpenID)
	if user.IsZero() {
		user.OID = OID
		user.OpenID = OpenID
		user.Name = fmt.Sprintf("用户%d", time.Now().Unix())
		user.Portrait = "https://oss.sites.ink/assets/default"
		err := dao.Create(Orm, user)
		if err != nil {
			return nil, err
		}
	}
	//CompanyOpenID := user.GetCompanyOpenID(CompanyID, OpenID)
	//err := Orm.Where("OpenID=?", OpenID).First(user).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	//log.Println(err)
	return user, nil
}
