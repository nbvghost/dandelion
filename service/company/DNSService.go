package company

import (
	"strings"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/gpa/types"
)

type DNSService struct {
}

func (service DNSService) ListDNS(OID types.PrimaryKey) []model.DNS {
	var list []model.DNS
	singleton.Orm().Model(model.DNS{}).Where(`"OID"=? and "Type"=?`, OID, model.DNSTypeA).Find(&list)
	return list
}
func (service DNSService) GetOID(domainName string) types.PrimaryKey {
	Orm := singleton.Orm()
	var d model.DNS
	Orm.Model(model.DNS{}).Where(`"Type"=? and "Domain"=?`, model.DNSTypeA, domainName).First(&d)
	return d.OID
}
func (service DNSService) GetDefaultDNS(OID types.PrimaryKey) *model.DNS {
	Orm := singleton.Orm()
	var list []model.DNS
	Orm.Model(model.DNS{}).Where(`"OID"=? and "Type"=?`, OID, model.DNSTypeA).Find(&list)
	for i := range list {
		if strings.EqualFold(list[i].Domain, "default") {
			return &list[i]
		}
	}
	return &model.DNS{}
}
