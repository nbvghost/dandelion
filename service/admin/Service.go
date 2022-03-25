package admin

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/dandelion/service/configuration"
	"github.com/nbvghost/dandelion/service/content"
	"github.com/nbvghost/gweb"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/library/util"

	"github.com/nbvghost/gpa/types"
	"github.com/nbvghost/tool/object"

	"github.com/nbvghost/tool/encryption"

	"github.com/nbvghost/glog"
)

type AdminService struct {
	model.BaseDao
	Organization  company.OrganizationService
	Configuration configuration.ConfigurationService
	Content       content.ContentService
}

func (service AdminService) AddItem(OID types.PrimaryKey, item *model.Admin) (err error) {
	item.OID = OID
	if strings.EqualFold(item.Account, "") {
		return errors.New("账号不允许为空")
	}

	item.Account = strings.ToLower(item.Account)
	item.PassWord = encryption.Md5ByString(item.PassWord)

	if strings.EqualFold(item.Account, "admin") || strings.EqualFold(item.Account, "manager") || strings.EqualFold(item.Account, "administrator") {
		return errors.New("此账号不允许注册")
	}
	return service.Add(singleton.Orm(), item)
}

func (service AdminService) GetItem(context *gweb.Context) (r constrain.IResult, err error) {

	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	ID := object.ParseUint(context.PathParams["ID"])
	item := &model.Admin{}
	err = service.Get(singleton.Orm(), types.PrimaryKey(ID), item)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", item)}, err
}
func (service AdminService) ListItem(context *gweb.Context) (r constrain.IResult, err error) {
	admin := context.Session.Attributes.Get(play.SessionAdmin).(*model.Admin)
	dts := &model.Datatables{}
	err = util.RequestBodyToJSON(context.Request.Body, dts)
	if err != nil {
		return nil, err
	}
	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(singleton.Orm(), dts, &[]model.Admin{}, admin.OID, "")
	return &result.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}, nil
}

func (service AdminService) DeleteItem(context *gweb.Context) (r constrain.IResult, err error) {
	ID := object.ParseUint(context.PathParams["ID"])
	item := &model.Admin{}
	Orm := singleton.Orm()

	err = service.Get(Orm, types.PrimaryKey(ID), item)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}
	if strings.EqualFold(item.Account, "admin") {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("admin不能删除"), "", nil)}, nil
	}

	err = service.Delete(Orm, item, types.PrimaryKey(ID))
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "删除成功", nil)}, err
}
func (service AdminService) ChangeAuthority(context *gweb.Context) (r constrain.IResult, err error) {
	Orm := singleton.Orm()
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	ID := object.ParseUint(context.PathParams["ID"])
	item := &model.Admin{}
	err = util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}

	var _admin model.Admin
	err = service.Get(Orm, types.PrimaryKey(ID), &_admin)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(_admin.Account, "admin") {
		admin := context.Session.Attributes.Get(play.SessionAdmin).(*model.Admin)
		if strings.EqualFold(admin.Account, _admin.Account) {

		} else {
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("无权修改admin账号权限"), "", nil)}, nil
		}
	}

	err = service.ChangeModel(Orm, types.PrimaryKey(ID), &model.Admin{Authority: item.Authority})
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}
func (service AdminService) ChangePassWork(context *gweb.Context) (r constrain.IResult, err error) {
	Orm := singleton.Orm()
	//ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	ID := object.ParseUint(context.PathParams["ID"])
	item := &model.Admin{}
	err = util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, err
	}

	var _admin model.Admin
	err = service.Get(Orm, types.PrimaryKey(ID), &_admin)
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

	err = service.ChangeModel(Orm, types.PrimaryKey(ID), item)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}

func (service AdminService) DelAdmin(ID uint) error {
	Orm := singleton.Orm()
	err := Orm.Delete(model.Admin{}, "ID=?", ID).Error
	return err
}
func (service AdminService) FindAdmin() []model.Admin {
	Orm := singleton.Orm()
	var list []model.Admin

	Orm.Find(&list)

	return list
}

/*
Account
PassWord
Domain
*/
func (service AdminService) InitOrganizationInfo(Account, PassWord, _Domain string) *result.ActionResult {
	//Orm := singleton.Orm()
	as := &result.ActionResult{}

	Domain := util.ParseDomain(_Domain)
	if len(Domain) == 0 {
		as = result.NewError(fmt.Errorf("不是一个有效的域名:%s", _Domain))
		return as
	}

	tx := singleton.Orm().Begin()
	defer func() {
		if as.Code == 0 {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	admin := service.FindAdminByAccount(tx, Account)
	if !admin.IsZero() {
		return result.NewError(fmt.Errorf("账号:%s被占用", Account))
	}

	_org := service.Organization.FindByDomain(tx, Domain)
	if _org != nil && _org.ID > 0 {

		as.Code = result.Fail
		as.Message = "域名：" + Domain + "已经被占用。"
		return as
	}

	shop := &model.Organization{}
	shop.Name = ""
	shop.Domain = Domain
	shop.Expire = time.Now().Add((365 * 1) * 24 * time.Hour)
	//shop.AdminID = types.PrimaryKey(adminID)
	if err := service.Organization.AddOrganization(tx, shop); err != nil {

		as.Code = result.Fail
		as.Message = err.Error()
		return as
	}
	admin.Account = strings.ToLower(Account)
	admin.PassWord = encryption.Md5ByString(PassWord)
	admin.OID = shop.ID
	if err := service.Add(tx, admin); err != nil {

		return result.NewError(err)
	}
	//tx.Model(model.Organization{}).Where(map[string]interface{}{"ID": shop.ID}).Updates(map[string]interface{}{"AdminID": admin.ID})

	as.Code = result.Success
	as.Message = "添加成功"

	var _Configuration model.Configuration

	err := service.Organization.FindWhere(tx, &_Configuration, `"K"=? and "OID"=?`, configuration.ConfigurationKeyBrokerageLeve1, shop.ID)
	if err != nil {

		return result.NewError(err)
	}
	if _Configuration.ID == 0 {
		a := model.Configuration{K: configuration.ConfigurationKeyBrokerageLeve1, V: "0"}
		a.OID = shop.ID
		err = service.Organization.Add(tx, &a)
		if err != nil {

			return result.NewError(err)
		}
	}

	_Configuration = model.Configuration{}
	err = service.Organization.FindWhere(tx, &_Configuration, `"K"=? and "OID"=?`, configuration.ConfigurationKeyBrokerageLeve2, shop.ID)
	if err != nil {

		return result.NewError(err)
	}
	if _Configuration.ID == 0 {
		a := model.Configuration{K: configuration.ConfigurationKeyBrokerageLeve2, V: "0"}
		a.OID = shop.ID
		err = service.Organization.Add(tx, &a)
		if err != nil {

			return result.NewError(err)
		}
	}

	_Configuration = model.Configuration{}
	err = service.Organization.FindWhere(tx, &_Configuration, `"K"=? and "OID"=?`, configuration.ConfigurationKeyBrokerageLeve3, shop.ID)
	if err != nil {

		return result.NewError(err)
	}
	if _Configuration.ID == 0 {
		a := model.Configuration{K: configuration.ConfigurationKeyBrokerageLeve3, V: "0"}
		a.OID = shop.ID
		err = service.Organization.Add(tx, &a)
		if err != nil {

			return result.NewError(err)
		}
	}

	_Configuration = model.Configuration{}
	err = service.Organization.FindWhere(tx, &_Configuration, `"K"=? and "OID"=?`, configuration.ConfigurationKeyBrokerageLeve4, shop.ID)
	if err != nil {

		return result.NewError(err)
	}
	if _Configuration.ID == 0 {
		a := model.Configuration{K: configuration.ConfigurationKeyBrokerageLeve4, V: "0"}
		a.OID = shop.ID
		err = service.Organization.Add(tx, &a)
		if err != nil {

			return result.NewError(err)
		}
	}

	_Configuration = model.Configuration{}
	err = service.Organization.FindWhere(tx, &_Configuration, `"K"=? and "OID"=?`, configuration.ConfigurationKeyBrokerageLeve5, shop.ID)
	if err != nil {

		return result.NewError(err)
	}
	if _Configuration.ID == 0 {
		a := model.Configuration{K: configuration.ConfigurationKeyBrokerageLeve5, V: "0"}
		a.OID = shop.ID
		err = service.Organization.Add(tx, &a)
		if err != nil {

			return result.NewError(err)
		}
	}

	_Configuration = model.Configuration{}
	err = service.Organization.FindWhere(tx, &_Configuration, `"K"=? and "OID"=?`, configuration.ConfigurationKeyBrokerageLeve6, shop.ID)
	if err != nil {

		return result.NewError(err)
	}
	if _Configuration.ID == 0 {
		a := model.Configuration{K: configuration.ConfigurationKeyBrokerageLeve6, V: "0"}
		a.OID = shop.ID
		err = service.Organization.Add(tx, &a)
		if err != nil {

			return result.NewError(err)
		}
	}
	err = service.Organization.FindWhere(tx, &_Configuration, `"K"=? and "OID"=?`, configuration.ConfigurationKeyBrokerageLeve6, shop.ID)
	if err != nil {
		return result.NewError(err)
	}
	if _Configuration.ID == 0 {
		a := model.Configuration{K: configuration.ConfigurationKeyBrokerageLeve6, V: "0"}
		a.OID = shop.ID
		err = service.Organization.Add(tx, &a)
		if err != nil {
			return result.NewError(err)
		}
	}

	err = service.Content.AddContentConfig(tx, shop)
	if glog.Error(err) {
		as.Code = result.SQLError
		as.Message = err.Error()
		return as
	}
	as.Code = 0
	return as
}
func (service AdminService) GetAdmin(ID uint) *model.Admin {
	Orm := singleton.Orm()
	admin := &model.Admin{}
	err := Orm.Where("ID=?", ID).First(admin).Error //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)

	return admin
}
func (service AdminService) FindAdminByID(Orm *gorm.DB, ID types.PrimaryKey) model.Admin {
	manager := model.Admin{}
	Orm.Where(map[string]interface{}{"ID": ID}).First(&manager) //SelectOne(user, "select * from User where Email=?", Email)
	return manager
}
func (service AdminService) FindAdminByAccount(Orm *gorm.DB, Account string) *model.Admin {
	manager := &model.Admin{}
	Orm.Where(map[string]interface{}{"Account": Account}).First(manager) //SelectOne(user, "select * from User where Email=?", Email)
	return manager
}
func (service AdminService) ManagerAction(context *gweb.Context) (r constrain.IResult, err error) {
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
}
func (service AdminService) ChangeAdmin(Account, Password string, ID types.PrimaryKey) error {
	Orm := singleton.Orm()
	return service.ChangeModel(Orm, ID, model.Admin{Account: Account, PassWord: encryption.Md5ByString(Password)})
}
