package service

import (
	"dandelion/app/play"
	"dandelion/app/service/dao"
	"dandelion/app/util"
	"errors"
	"strconv"
	"strings"

	"time"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/tool"
)

type AdminService struct {
	dao.BaseDao
	Organization  OrganizationService
	Configuration ConfigurationService
}

func (service AdminService) AddItem(context *gweb.Context) gweb.Result {
	item := &dao.Admin{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	if strings.EqualFold(item.Account, "") {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("账号不允许为空"), "", nil)}
	}

	item.Account = strings.ToLower(item.Account)
	item.PassWord = tool.Md5ByString(item.PassWord)

	if strings.EqualFold(item.Account, "admin") || strings.EqualFold(item.Account, "manager") || strings.EqualFold(item.Account, "administrator") {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("此账号不允许注册"), "", nil)}
	}

	err = service.Add(dao.Orm(), item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "添加成功", nil)}
}

func (service AdminService) GetItem(context *gweb.Context) gweb.Result {

	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Admin{}
	err := service.Get(dao.Orm(), ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", item)}
}
func (service AdminService) ListItem(context *gweb.Context) gweb.Result {
	admin := context.Session.Attributes.Get(play.SessionAdmin).(*dao.Admin)
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(dao.Orm(), dts, &[]dao.Admin{}, admin.OID)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}

func (service AdminService) DeleteItem(context *gweb.Context) gweb.Result {
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Admin{}
	Orm := dao.Orm()

	err := service.Get(Orm, ID, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}
	if strings.EqualFold(item.Account, "admin") {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("admin不能删除"), "", nil)}
	}

	err = service.Delete(Orm, item, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (service AdminService) ChangeAuthority(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Admin{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	var _admin dao.Admin
	service.Get(Orm, ID, &_admin)
	if strings.EqualFold(_admin.Account, "admin") {
		admin := context.Session.Attributes.Get(play.SessionAdmin).(*dao.Admin)
		if strings.EqualFold(admin.Account, _admin.Account) {

		} else {
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("无权修改admin账号权限"), "", nil)}
		}
	}

	err = service.ChangeModel(Orm, ID, &dao.Admin{Authority: item.Authority})
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}
func (service AdminService) ChangePassWork(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	item := &dao.Admin{}
	err := util.RequestBodyToJSON(context.Request.Body, item)
	if err != nil {
		return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "", nil)}
	}

	var _admin dao.Admin
	service.Get(Orm, ID, &_admin)
	if strings.EqualFold(_admin.Account, "admin") {
		admin := context.Session.Attributes.Get(play.SessionAdmin).(*dao.Admin)
		if strings.EqualFold(admin.Account, _admin.Account) {

		} else {
			return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(errors.New("无权修改admin账号密码"), "", nil)}
		}
	}

	item.PassWord = tool.Md5ByString(item.PassWord)

	err = service.ChangeModel(Orm, ID, item)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "修改成功", nil)}
}

func (service AdminService) DelAdmin(ID uint64) error {
	Orm := dao.Orm()
	err := Orm.Delete(dao.Admin{}, "ID=?", ID).Error
	return err
}
func (service AdminService) FindAdmin() []dao.Admin {
	Orm := dao.Orm()
	var list []dao.Admin

	Orm.Find(&list)

	return list
}

func (self AdminService) AddAdmin(Name, Password, Domain string) *dao.ActionStatus {
	Orm := dao.Orm()
	as := &dao.ActionStatus{}

	tx := Orm.Begin()

	admin := &dao.Admin{}
	admin.Account = Name
	//admin.Email = Email
	//admin.Tel = Tel

	admin.LastLoginAt = time.Now()
	if haveAdmin := self.FindAdminByAccount(tx, admin.Account); haveAdmin.ID != 0 {
		tx.Rollback()
		as.Success = false
		as.Message = "这个账号已经存在"
		return as
	}

	admin.PassWord = tool.Md5ByString(Password)

	_org := self.Organization.FindByDomain(tx, Domain)
	if _org.ID != 0 {
		tx.Rollback()
		as.Success = false
		as.Message = "域名：" + Domain + "已经被占用。"
		return as
	}

	shop := &dao.Organization{}
	shop.Name = Name + "的店铺"
	shop.Domain = Domain
	shop.Expire = time.Now().Add((365 * 1) * 24 * time.Hour)

	if err := self.Organization.AddOrganization(tx, shop); err != nil {
		tx.Rollback()
		as.Success = false
		as.Message = err.Error()
		return as
	}

	admin.OID = shop.ID

	if err := self.Add(tx, admin); err != nil {
		tx.Rollback()
		as.Success = false
		as.Message = err.Error()
		return as
	}

	tx.Commit()
	as.Success = true
	as.Message = "添加成功"

	var _Configuration dao.Configuration

	self.Organization.FindWhere(tx, &_Configuration, "K=? and OID=?", play.ConfigurationKey_BrokerageLeve1, shop.ID)
	if _Configuration.ID == 0 {
		a := dao.Configuration{K: play.ConfigurationKey_BrokerageLeve1, V: "0"}
		a.OID = shop.ID
		self.Organization.Add(tx, &a)
	}

	_Configuration = dao.Configuration{}
	self.Organization.FindWhere(tx, &_Configuration, "K=? and OID=?", play.ConfigurationKey_BrokerageLeve2, shop.ID)
	if _Configuration.ID == 0 {
		a := dao.Configuration{K: play.ConfigurationKey_BrokerageLeve2, V: "0"}
		a.OID = shop.ID
		self.Organization.Add(tx, &a)
	}

	_Configuration = dao.Configuration{}
	self.Organization.FindWhere(tx, &_Configuration, "K=? and OID=?", play.ConfigurationKey_BrokerageLeve3, shop.ID)
	if _Configuration.ID == 0 {
		a := dao.Configuration{K: play.ConfigurationKey_BrokerageLeve3, V: "0"}
		a.OID = shop.ID
		self.Organization.Add(tx, &a)
	}

	_Configuration = dao.Configuration{}
	self.Organization.FindWhere(tx, &_Configuration, "K=? and OID=?", play.ConfigurationKey_BrokerageLeve4, shop.ID)
	if _Configuration.ID == 0 {
		a := dao.Configuration{K: play.ConfigurationKey_BrokerageLeve4, V: "0"}
		a.OID = shop.ID
		self.Organization.Add(tx, &a)
	}

	_Configuration = dao.Configuration{}
	self.Organization.FindWhere(tx, &_Configuration, "K=? and OID=?", play.ConfigurationKey_BrokerageLeve5, shop.ID)
	if _Configuration.ID == 0 {
		a := dao.Configuration{K: play.ConfigurationKey_BrokerageLeve5, V: "0"}
		a.OID = shop.ID
		self.Organization.Add(tx, &a)
	}

	_Configuration = dao.Configuration{}
	self.Organization.FindWhere(tx, &_Configuration, "K=? and OID=?", play.ConfigurationKey_BrokerageLeve6, shop.ID)
	if _Configuration.ID == 0 {
		a := dao.Configuration{K: play.ConfigurationKey_BrokerageLeve6, V: "0"}
		a.OID = shop.ID
		self.Organization.Add(tx, &a)
	}

	return as
}
func (service AdminService) GetAdmin(ID uint64) *dao.Admin {
	Orm := dao.Orm()
	admin := &dao.Admin{}
	err := Orm.Where("ID=?", ID).First(admin).Error //SelectOne(user, "select * from User where Email=?", Email)
	tool.CheckError(err)
	return admin
}

func (service AdminService) FindAdminByAccount(Orm *gorm.DB, Account string) *dao.Admin {
	manager := &dao.Admin{}
	err := Orm.Where(map[string]interface{}{"Account": Account}).First(manager).Error //SelectOne(user, "select * from User where Email=?", Email)
	tool.CheckError(err)
	return manager
}
func (service AdminService) ManagerAction(context *gweb.Context) gweb.Result {
	Orm := dao.Orm()
	admin := context.Session.Attributes.Get(play.SessionAdmin).(*dao.Admin)

	action := context.Request.URL.Query().Get("action")
	switch action {
	case play.ActionKey_del:
		dts := &dao.Admin{}
		util.RequestBodyToJSON(context.Request.Body, dts)
		//manager
		sd := &dao.Admin{}
		service.Get(Orm, dts.ID, sd)
		if strings.EqualFold(sd.Account, "manager") {
			//self.ChangeModel(Orm, dts.ID, &dao.Manager{Account: dts.Account, PassWord: tool.Md5(dts.PassWord)})
			return &gweb.JsonResult{Data: dao.ActionStatus{Success: false, Message: "这个用户不能删除", Data: nil}}
		} else {
			service.Delete(Orm, &dao.Admin{}, dts.ID)
			return &gweb.JsonResult{Data: dao.ActionStatus{Success: true, Message: "删除成功", Data: nil}}
		}

	case play.ActionKey_change:

		dts := &dao.Admin{}
		util.RequestBodyToJSON(context.Request.Body, dts)

		service.ChangeModel(Orm, dts.ID, &dao.Admin{PassWord: tool.Md5ByString(dts.PassWord)})
		return &gweb.JsonResult{Data: dao.ActionStatus{Success: true, Message: "修改成功", Data: nil}}

	case play.ActionKey_add:
		dts := &dao.Admin{}
		util.RequestBodyToJSON(context.Request.Body, dts)

		manager := context.Session.Attributes.Get(play.SessionAdmin).(*dao.Admin)

		if !strings.EqualFold(manager.Account, "manager") {
			return &gweb.JsonResult{Data: dao.ActionStatus{Success: false, Message: "您没有添加账号的权限", Data: nil}}
		}

		sk := service.FindAdminByAccount(Orm, dts.Account)
		if sk.ID == 0 {
			service.Add(Orm, dts)
			return &gweb.JsonResult{Data: dao.ActionStatus{Success: true, Message: "添加成功", Data: nil}}
		} else {
			return &gweb.JsonResult{Data: dao.ActionStatus{Success: false, Message: "账号重复", Data: nil}}
		}
	case "list":
		dts := &dao.Datatables{}
		util.RequestBodyToJSON(context.Request.Body, dts)
		draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.Admin{}, admin.OID)
		return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
	}

	return &gweb.JsonResult{Data: dao.ActionStatus{Success: false, Message: "", Data: nil}}
}
func (service AdminService) ChangeAdmin(Name, Password string, ID uint64) error {
	Orm := dao.Orm()
	return service.ChangeModel(Orm, ID, dao.Admin{Account: Name, PassWord: tool.Md5ByString(Password)})
}
