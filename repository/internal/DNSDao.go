package internal

import (
	"context"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type DNSDao struct{}

func (m DNSDao) ListDNS(ctx context.Context, OID dao.PrimaryKey) []model.DNS {
	var list []model.DNS
	db.GetDB(ctx).Model(model.DNS{}).Where(`"OID"=? and "Type"=?`, OID, model.DNSTypeA).Find(&list)
	return list
}
func (m DNSDao) GetOID(ctx context.Context, domainName string) dao.PrimaryKey {
	Orm := db.GetDB(ctx)
	var d model.DNS
	Orm.Model(model.DNS{}).Where(`"Type"=? and "Domain"=?`, model.DNSTypeA, domainName).First(&d)
	return d.OID
}
func (m DNSDao) GetDefaultDNS(ctx context.Context, OID dao.PrimaryKey) *model.DNS {
	Orm := db.GetDB(ctx)
	var list []model.DNS
	Orm.Model(model.DNS{}).Where(`"OID"=? and "Type"=?`, OID, model.DNSTypeA).Find(&list)
	if len(list) > 0 {
		return &list[0]
	}
	return &model.DNS{}
}
