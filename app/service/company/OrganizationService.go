package company

import (
	"errors"
	"github.com/nbvghost/glog"

	"github.com/nbvghost/dandelion/app/service/dao"

	"github.com/jinzhu/gorm"
)

type OrganizationService struct {
	dao.BaseDao
}

func (service OrganizationService) AddOrganizationBlockAmount(Orm *gorm.DB, OID uint64, Menoy int64) error {

	var org dao.Organization
	err := service.Get(Orm, OID, &org)
	if err != nil {
		return err
	}

	tm := int64(org.BlockAmount) + Menoy
	if tm < 0 {
		return errors.New("冻结金额不足，无法扣款")
	}

	err = service.ChangeMap(dao.Orm(), OID, &dao.Organization{}, map[string]interface{}{"BlockAmount": tm})
	return err
}
func (service OrganizationService) FindByDomain(Orm *gorm.DB, Domain string) *dao.Organization {
	manager := &dao.Organization{}
	//err := Orm.Where("Domain=?", Domain).First(manager).Error //SelectOne(user, "select * from User where Email=?", Email)
	Orm.Where(map[string]interface{}{"Domain": Domain}).First(manager) //SelectOne(user, "select * from User where Email=?", Email)
	if manager.ID == 0 {
		return nil
	}
	return manager
}
func (service OrganizationService) GetOrganization(ID uint64) *dao.Organization {
	Orm := dao.Orm()
	target := &dao.Organization{}
	service.Get(Orm, ID, target)
	return target
}
func (service OrganizationService) AddOrganization(Orm *gorm.DB, shop *dao.Organization) error {
	org := service.FindByDomain(Orm, shop.Domain)
	if org != nil && org.ID > 0 {
		return errors.New("域名：" + shop.Domain + "已经被占用，请试试其它域名")
	}
	return service.Add(Orm, shop)
}
func (service OrganizationService) DelCompany(ID uint64) error {
	Orm := dao.Orm()
	return service.Delete(Orm, dao.Organization{}, ID)
}
func (service OrganizationService) ChangeOrganization(ID uint64, shop *dao.Organization) error {
	Orm := dao.Orm()
	//return Orm.Save(article).Error
	//err := db.Orm.Save(shop).Error
	org := service.FindByDomain(Orm, shop.Domain)
	if org.ID != 0 && org.ID != shop.ID {
		return errors.New("企业信息不存在")
	} else {
		shop.Amount = org.Amount
		shop.BlockAmount = org.BlockAmount
		shop.Vip = org.Vip
		shop.Expire = org.Expire
		err := service.ChangeModel(Orm, ID, shop)
		if glog.Error(err) {
			return err
		} else {
			return nil
		}
	}

}

//Execute(Session *Session,Request *http.Request)(bool,Result)
/*func (self OrganizationService) ReadOrganization(Context *gweb.Context) (bool, gweb.Result) {

	var domain string
	fmt.Println(Context.Request.Host)
	if strings.Contains(Context.Request.Host, ".d.") {
		domains := strings.Split(Context.Request.Host, ".d.")
		fmt.Println(domains)
		domain = domains[0]
	} else {
		domain = ""
	}

	organization := self.FindByDomain(dao.Orm(), domain)
	if organization.ID == 0 {
		//context.Response.Header().Add("Login-Status", "0")
		//context.Response.Write([]byte(util.StructToJSON(&result.ActionResult{Code: result.ActionFail, Message: "找不到组织信息", Data: nil})))
		return false, &gweb.JsonResult{Data: &result.ActionResult{Code: result.ActionFail, Message: "找不到组织信息", Data: nil}}
	}

	if Context.Session.Attributes.Get(play.SessionOrganization) == nil {
		Context.Session.Attributes.Put(play.SessionOrganization, organization)
	}
	//context.Response.Header().Add("Login-Status", "1")
	return true, nil

}
*/
