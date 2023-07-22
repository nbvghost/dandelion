package company

import (
	"errors"
	"github.com/nbvghost/dandelion/library/db"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type OrganizationService struct {
	model.BaseDao
}

func (service OrganizationService) AddOrganizationBlockAmount(Orm *gorm.DB, OID dao.PrimaryKey, Menoy int64) error {

	org := dao.GetByPrimaryKey(Orm, &model.Organization{}, OID).(*model.Organization)
	if org.IsZero() {
		return errors.New("没找到数据")
	}

	tm := int64(org.BlockAmount) + Menoy
	if tm < 0 {
		return errors.New("冻结金额不足，无法扣款")
	}

	err := dao.UpdateByPrimaryKey(db.Orm(), &model.Organization{}, OID, map[string]interface{}{"BlockAmount": tm})
	return err
}
func (service OrganizationService) FindByName(Orm *gorm.DB, Name string) *model.Organization {
	manager := &model.Organization{}
	//err := Orm.Where("Domain=?", Domain).First(manager).Error //SelectOne(user, "select * from User where Email=?", Email)
	Orm.Where(map[string]interface{}{"Name": Name}).First(manager) //SelectOne(user, "select * from User where Email=?", Email)
	return manager
}
func (service OrganizationService) FindByDomain(Orm *gorm.DB, Domain string) *model.Organization {
	manager := &model.Organization{}
	var dns model.DNS
	Orm.Model(model.DNS{}).Where(`"Type"=? and "Domain"=?`, model.DNSTypeA, Domain).First(&dns)
	if dns.IsZero() {
		return manager
	}
	//err := Orm.Where("Domain=?", Domain).First(manager).Error //SelectOne(user, "select * from User where Email=?", Email)
	Orm.Where(map[string]interface{}{"ID": dns.OID}).First(manager) //SelectOne(user, "select * from User where Email=?", Email)
	return manager
}
func (service OrganizationService) GetOrganization(ID dao.PrimaryKey) dao.IEntity {
	Orm := db.Orm()
	//target := model.Organization{}
	//service.Get(Orm, ID, &target)
	return dao.GetByPrimaryKey(Orm, &model.Organization{}, ID)
}

func (service OrganizationService) DelCompany(ID dao.PrimaryKey) error {
	Orm := db.Orm()
	return dao.DeleteByPrimaryKey(Orm, &model.Organization{}, ID)
}
func (service OrganizationService) ChangeOrganization(ID dao.PrimaryKey, shop *model.Organization) error {
	Orm := db.Orm()
	//return Orm.Save(article).Error
	//err := db.Orm.Save(shop).Error
	org := service.GetOrganization(ID).(*model.Organization)
	if org.IsZero() {
		return errors.New("企业信息不存在")
	} else {
		shop.Amount = org.Amount
		shop.BlockAmount = org.BlockAmount
		shop.Vip = org.Vip
		shop.Expire = org.Expire
		err := dao.UpdateByPrimaryKey(Orm, &model.Organization{}, ID, shop)
		if err != nil {
			return err
		} else {
			return nil
		}
	}

}

/*func (service OrganizationService) FindByAdminID(Orm *gorm.DB, adminID dao.PrimaryKey) *model.Organization {
	manager := &model.Organization{}
	//err := Orm.Where("Domain=?", Domain).First(manager).Error //SelectOne(user, "select * from User where Email=?", Email)
	Orm.Where(map[string]interface{}{"AdminID": adminID}).First(manager) //SelectOne(user, "select * from User where Email=?", Email)
	if manager.ID == 0 {
		return nil
	}
	return manager
}*/

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

	organization := self.FindByDomain(singleton.Orm(), domain)
	if organization.ID == 0 {
		//context.Response.Header().Add("Login-Status", "0")
		//context.Response.Write([]byte(util.StructToJSON(&result.ActionResult{Code: result.Fail, Message: "找不到组织信息", Data: nil})))
		return false, &gweb.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: "找不到组织信息", Data: nil}}
	}

	if Context.Session.Attributes.Get(play.SessionOrganization) == nil {
		Context.Session.Attributes.Put(play.SessionOrganization, organization)
	}
	//context.Response.Header().Add("Login-Status", "1")
	return true, nil

}
*/
