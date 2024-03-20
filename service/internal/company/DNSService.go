package company

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type DNSService struct {
}

func (service DNSService) ListDNS(OID dao.PrimaryKey) []model.DNS {
	var list []model.DNS
	db.Orm().Model(model.DNS{}).Where(`"OID"=? and "Type"=?`, OID, model.DNSTypeA).Find(&list)
	return list
}
func (service DNSService) GetOID(domainName string) dao.PrimaryKey {
	Orm := db.Orm()
	var d model.DNS
	Orm.Model(model.DNS{}).Where(`"Type"=? and "Domain"=?`, model.DNSTypeA, domainName).First(&d)
	return d.OID
}
func (service DNSService) GetDefaultDNS(OID dao.PrimaryKey) *model.DNS {
	Orm := db.Orm()
	var list []model.DNS
	Orm.Model(model.DNS{}).Where(`"OID"=? and "Type"=?`, OID, model.DNSTypeA).Find(&list)
	if len(list) > 0 {
		return &list[0]
	}
	return &model.DNS{}
}
