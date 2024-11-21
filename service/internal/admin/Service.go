package admin

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service/internal/company"
	"github.com/nbvghost/dandelion/service/internal/configuration"
	"github.com/nbvghost/dandelion/service/internal/content"
	"github.com/nbvghost/dandelion/service/internal/wechat"
	"github.com/nbvghost/tool/encryption"

	"log"
	"strings"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
)

type AdminService struct {
	model.BaseDao
	Organization  company.OrganizationService
	Configuration configuration.ConfigurationService
	Content       content.ContentService
	WxService     wechat.WxService
}

func (m AdminService) AddItem(OID dao.PrimaryKey, item *model.Admin) (err error) {
	item.OID = OID
	if strings.EqualFold(item.Account, "") {
		return errors.New("账号不允许为空")
	}

	item.Account = strings.TrimSpace(item.Account)
	item.PassWord = strings.TrimSpace(item.PassWord)
	//item.PassWord = encryption.Md5ByString(item.PassWord)

	if strings.EqualFold(item.Account, "admin") || strings.EqualFold(item.Account, "manager") || strings.EqualFold(item.Account, "administrator") {
		return errors.New("此账号不允许注册")
	}
	return dao.Create(db.Orm(), item)
}

func (m AdminService) GetItem(context constrain.IContext, ID dao.PrimaryKey) (r constrain.IResult, err error) {

	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])

	item := dao.GetByPrimaryKey(db.Orm(), &model.Admin{}, dao.PrimaryKey(ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", item)}, err
}
func (m AdminService) ListItem(context constrain.IContext, admin *model.Admin) (r constrain.IResult, err error) {
	//admin := context.Session.Attributes.Get(play.SessionAdmin).(*model.Admin)
	dts := &model.Datatables{}
	/*err = util.RequestBodyToJSON(context.Request.Body, dts)
	if err != nil {
		return nil, err
	}*/
	draw, recordsTotal, recordsFiltered, list := m.DatatablesListOrder(db.Orm(), dts, &[]model.Admin{}, admin.OID, "")
	return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, nil
}

func (m AdminService) DeleteItem(context constrain.IContext, ID dao.PrimaryKey) (r constrain.IResult, err error) {
	//ID := object.ParseUint(context.PathParams["ID"])

	Orm := db.Orm()
	item := dao.GetByPrimaryKey(Orm, &model.Admin{}, dao.PrimaryKey(ID)).(*model.Admin)
	if item.IsZero() {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, gorm.ErrRecordNotFound
	}
	if strings.EqualFold(item.Account, "admin") {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("admin不能删除"), "", nil)}, nil
	}

	err = dao.DeleteByPrimaryKey(Orm, item, dao.PrimaryKey(ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}, err
}
func (m AdminService) ChangeAuthority(context constrain.IContext, admin *model.Admin, ID dao.PrimaryKey) (r constrain.IResult, err error) {
	Orm := db.Orm()
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	//ID := object.ParseUint(context.PathParams["ID"])
	//item := &model.Admin{}
	/*err = util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}*/

	_admin := dao.GetByPrimaryKey(Orm, &model.Admin{}, dao.PrimaryKey(ID)).(*model.Admin)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(_admin.Account, "admin") {
		//admin := context.Session.Attributes.Get(play.SessionAdmin).(*model.Admin)
		if strings.EqualFold(admin.Account, _admin.Account) {

		} else {
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("无权修改admin账号权限"), "", nil)}, nil
		}
	}

	//err = dao.UpdateByPrimaryKey(Orm, &model.Admin{}, dao.PrimaryKey(ID), &model.Admin{Authority: item.Authority})
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}

/*func (service AdminService) ChangePassWork(context *gweb.Context) (r constrain.IResult, err error) {
	Orm := singleton.Orm()
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	ID := object.ParseUint(context.PathParams["ID"])
	item := &model.Admin{}
	err = util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}

	var _admin model.Admin
	err = service.Get(Orm, dao.PrimaryKey(ID), &_admin)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(_admin.Account, "admin") {
		admin := context.Session.Attributes.Get(play.SessionAdmin).(*model.Admin)
		if strings.EqualFold(admin.Account, _admin.Account) {

		} else {
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("无权修改admin账号密码"), "", nil)}, nil
		}
	}

	item.PassWord = encryption.Md5ByString(item.PassWord)

	err = service.ChangeModel(Orm, dao.PrimaryKey(ID), item)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}*/

func (m AdminService) DelAdmin(ID uint) error {
	Orm := db.Orm()
	err := Orm.Delete(model.Admin{}, "ID=?", ID).Error
	return err
}
func (m AdminService) FindAdmin() []model.Admin {
	Orm := db.Orm()
	var list []model.Admin

	Orm.Find(&list)

	return list
}

/*
Account
PassWord
Domain
*/
func (m AdminService) InitOrganizationInfo(account string,password string) (admin *model.Admin, err error) {
	//Orm := singleton.Orm()

	/*mDomain := util.ParseDomain(domain)
	if len(mDomain) == 0 {

		return nil, errors.Errorf("不是一个有效的域名:%s", domain)
	}*/

	tx := db.Orm().Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	/*_org := service.Organization.FindByDomain(tx, mDomain)
	if _org != nil && _org.ID > 0 {

		return nil, errors.Errorf("域名：" + mDomain + "已经被占用。")
	}*/

	admin = m.FindAdminByAccount(tx, account)

	shop := m.Organization.GetOrganization(admin.ID).(*model.Organization)
	if shop.IsZero() {
		shop.Name = ""
		shop.Expire = time.Now().Add((365 * 1) * 24 * time.Hour)
		if err = dao.Create(tx, shop); err != nil {
			return nil, err
		}
	}

	if admin.IsZero() {
		admin.Account = strings.TrimSpace(account)
		admin.PassWord = encryption.Md5ByString(strings.TrimSpace(password))
		//admin.PassWord = encryption.Md5ByString(PassWord)
		//admin.OID = shop.ID
		admin.Initiator = true
		admin.LastLoginAt = time.Now()
		admin.OID = shop.ID
		if err = dao.Create(tx, admin); err != nil {
			return nil, err
		}
	}

	domain := fmt.Sprintf("default")

	var dns model.DNS
	tx.Model(&model.DNS{}).Where(`"Type"=? and "OID"=?`, model.DNSTypeA, shop.ID).First(&dns)
	if dns.IsZero() {
		dns.Type = model.DNSTypeA
		dns.Domain = domain
		dns.OID = shop.ID
		if err = tx.Model(&model.DNS{}).Create(&dns).Error; err != nil {
			return nil, err
		}
		//return nil, fmt.Errorf("存在相同的DNS信息,Domain=%s,Type=%s", domain, model.DNSTypeA)
	}

	err = dao.Create(tx, &model.ExpressTemplate{
		Entity:   dao.Entity{},
		OID:      shop.ID,
		Name:     "默认",
		Drawee:   "BUSINESS",
		Type:     "ITEM",
		Template: `{"Default":{"Areas":[],"N":0,"M":0,"AN":0,"ANM":0},"Items":[]}`,
		Free:     "[]",
	})
	if err != nil {
		return nil, err
	}

	var config *model.Configuration

	config = dao.GetBy(tx, &model.Configuration{}, map[string]any{"K": model.ConfigurationKeyBrokerageLeve1, "OID": shop.ID}).(*model.Configuration)
	if config.IsZero() {
		a := model.Configuration{K: model.ConfigurationKeyBrokerageLeve1, V: "0"}
		a.OID = shop.ID
		err = dao.Create(tx, &a)
		if err != nil {
			return nil, err
		}
	}

	config = dao.GetBy(tx, &model.Configuration{}, map[string]any{"K": model.ConfigurationKeyBrokerageLeve2, "OID": shop.ID}).(*model.Configuration)
	if config.IsZero() {
		a := model.Configuration{K: model.ConfigurationKeyBrokerageLeve2, V: "0"}
		a.OID = shop.ID
		err = dao.Create(tx, &a)
		if err != nil {

			return nil, err
		}
	}

	config = dao.GetBy(tx, &model.Configuration{}, map[string]any{"K": model.ConfigurationKeyBrokerageLeve3, "OID": shop.ID}).(*model.Configuration)
	if config.IsZero() {
		a := model.Configuration{K: model.ConfigurationKeyBrokerageLeve3, V: "0"}
		a.OID = shop.ID
		err = dao.Create(tx, &a)
		if err != nil {

			return nil, err
		}
	}

	config = dao.GetBy(tx, &model.Configuration{}, map[string]any{"K": model.ConfigurationKeyBrokerageLeve4, "OID": shop.ID}).(*model.Configuration)
	if config.IsZero() {
		a := model.Configuration{K: model.ConfigurationKeyBrokerageLeve4, V: "0"}
		a.OID = shop.ID
		err = dao.Create(tx, &a)
		if err != nil {

			return nil, err
		}
	}

	config = dao.GetBy(tx, &model.Configuration{}, map[string]any{"K": model.ConfigurationKeyBrokerageLeve5, "OID": shop.ID}).(*model.Configuration)
	if config.IsZero() {
		a := model.Configuration{K: model.ConfigurationKeyBrokerageLeve5, V: "0"}
		a.OID = shop.ID
		err = dao.Create(tx, &a)
		if err != nil {

			return nil, err
		}
	}

	config = dao.GetBy(tx, &model.Configuration{}, map[string]any{"K": model.ConfigurationKeyBrokerageLeve6, "OID": shop.ID}).(*model.Configuration)
	if config.IsZero() {
		a := model.Configuration{K: model.ConfigurationKeyBrokerageLeve6, V: "0"}
		a.OID = shop.ID
		err = dao.Create(tx, &a)
		if err != nil {
			return nil, err
		}
	}
	err = repository.ContentConfigDao.AddContentConfig(tx, shop)
	if err != nil {
		return nil, err
	}

	err = m.WxService.InitWechatConfig(tx, shop.ID)
	if err != nil {
		return nil, err
	}
	return admin, err
}
func (m AdminService) GetAdmin(ID dao.PrimaryKey) *model.Admin {
	Orm := db.Orm()
	admin := &model.Admin{}
	err := Orm.Where(`"ID"=?`, ID).First(admin).Error //SelectOne(user, "select * from User where Email=?", Email)
	if err != nil {
		log.Println(err)
	}
	return admin
}
func (m AdminService) FindAdminByID(Orm *gorm.DB, ID dao.PrimaryKey) model.Admin {
	manager := model.Admin{}
	Orm.Where(map[string]interface{}{"ID": ID}).First(&manager) //SelectOne(user, "select * from User where Email=?", Email)
	return manager
}
func (m AdminService) FindAdminByAccount(Orm *gorm.DB, Account string) *model.Admin {
	manager := &model.Admin{}
	Orm.Where(map[string]interface{}{"Account": Account}).First(manager) //SelectOne(user, "select * from User where Email=?", Email)
	return manager
}
func (m AdminService) FindAdminByAccountAndPassWord(Orm *gorm.DB, Account string, PassWord string) *model.Admin {
	manager := &model.Admin{}
	Orm.Where(map[string]interface{}{"Account": Account, "PassWord": PassWord}).First(manager) //SelectOne(user, "select * from User where Email=?", Email)
	return manager
}

/*func (service AdminService) ManagerAction(context *gweb.Context) (r constrain.IResult, err error) {
	Orm := singleton.Orm()
	admin := context.Session.Attributes.Get(play.SessionAdmin).(*model.Admin)

	action := context.Request.URL.Query().Get("action")
	switch action {
	case play.ActionKey_del:
		dts := &model.Admin{}
		err = util.RequestBodyToJSON(context.Request.Body, dts)
		if err != nil {
			return nil, err
		}
		//manager
		sd := &model.Admin{}
		err = service.Get(Orm, dts.ID, sd)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(sd.Account, "manager") {
			//self.ChangeModel(Orm, dts.ID, &model.Manager{Account: dts.Account, PassWord: tool.Md5(dts.PassWord)})
			return &result.JsonResult{Data: result.ActionResult{Code: result.Fail, Message: "这个用户不能删除", Data: nil}}, nil
		} else {
			err = service.Delete(Orm, &model.Admin{}, dts.ID)
			if err != nil {
				return nil, err
			}
			return &result.JsonResult{Data: result.ActionResult{Code: result.Success, Message: "删除成功", Data: nil}}, nil
		}

	case play.ActionKey_change:

		dts := &model.Admin{}
		err = util.RequestBodyToJSON(context.Request.Body, dts)
		if err != nil {
			return nil, err
		}

		err = service.ChangeModel(Orm, dts.ID, &model.Admin{PassWord: encryption.Md5ByString(dts.PassWord)})
		if err != nil {
			return nil, err
		}
		return &result.JsonResult{Data: result.ActionResult{Code: result.Success, Message: "修改成功", Data: nil}}, nil

	case play.ActionKey_add:

	case "list":
		dts := &model.Datatables{}
		err = util.RequestBodyToJSON(context.Request.Body, dts)
		if err != nil {
			return nil, err
		}
		draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]model.Admin{}, admin.OID, "")
		return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, nil
	}

	return &result.JsonResult{Data: result.ActionResult{Code: result.Fail, Message: "", Data: nil}}, nil
}*/
/*func (service AdminService) ChangeAdmin(Account, Password string, ID dao.PrimaryKey) error {
	Orm := singleton.Orm()
	return service.ChangeModel(Orm, ID, model.Admin{Account: Account, PassWord: encryption.Md5ByString(Password)})
}
*/
