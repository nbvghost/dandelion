package view

import (
	"github.com/nbvghost/dandelion/config"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/repository"
)

type Index struct {
	Organization *model.Organization `mapping:""`
	Admin        *model.Admin        `mapping:""`
}
type IndexView struct {
	extends.ViewBase
	Organization *model.Organization
	Admin        *model.Admin
	OSSHost      string
	DNSList      []model.DNS
}

func (m *Index) Render(ctx constrain.IContext) (r constrain.IViewResult, err error) {
	ossHost, err := ctx.Etcd().SelectOutsideServer(config.MicroServerOSS)
	if err != nil {
		return nil, err
	}
	dnsList := repository.DNSDao.ListDNS(ctx, m.Organization.ID)
	return &IndexView{Organization: m.Organization, OSSHost: ossHost, Admin: m.Admin, DNSList: dnsList}, nil
}
