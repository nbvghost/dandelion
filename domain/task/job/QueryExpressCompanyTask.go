package job

import (
	"context"
	"log"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/wechat"
)

type QueryExpressCompanyTask struct {
	WxService wechat.WxService
}

func (m *QueryExpressCompanyTask) Run() error {
	wxConfigList := m.WxService.MiniProgram(singleton.Orm())
	for _, config := range wxConfigList {
		Orm := singleton.Orm()
		deliverys, err := m.WxService.GetDeliveryList(m.WxService.GetAccessToken(config.(*model.WechatConfig)))
		if err != nil {
			log.Println(err)
			continue
		}
		for _, delivery := range deliverys {
			has := dao.GetBy(Orm, entity.ExpressCompany, map[string]any{"Key": delivery.DeliveryId})
			if has.IsZero() {
				err := dao.Create(Orm, &model.ExpressCompany{
					Key:  delivery.DeliveryId,
					Name: delivery.DeliveryName,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func NewQueryExpressCompanyTask(context context.Context) Job {
	return &QueryExpressCompanyTask{}
}
