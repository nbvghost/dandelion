package internal

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type DNSDao struct{}

func (m DNSDao) ListDNS(OID dao.PrimaryKey) []model.DNS {
	var list []model.DNS
	db.Orm().Model(model.DNS{}).Where(`"OID"=? and "Type"=?`, OID, model.DNSTypeA).Find(&list)
	return list
}
func (m DNSDao) GetOID(domainName string) dao.PrimaryKey {
	Orm := db.Orm()
	var d model.DNS
	Orm.Model(model.DNS{}).Where(`"Type"=? and "Domain"=?`, model.DNSTypeA, domainName).First(&d)
	return d.OID
}
func (m DNSDao) GetDefaultDNS(OID dao.PrimaryKey) *model.DNS {
	Orm := db.Orm()
	var list []model.DNS
	Orm.Model(model.DNS{}).Where(`"OID"=? and "Type"=?`, OID, model.DNSTypeA).Find(&list)
	if len(list) > 0 {
		return &list[0]
	}
	return &model.DNS{}
}
