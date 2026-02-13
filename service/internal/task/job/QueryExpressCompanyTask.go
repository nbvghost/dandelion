package job

import (
	"context"
	"log"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/wechat"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type QueryExpressCompanyTask struct {
	WxService          wechat.WxService
	AccessTokenService wechat.AccessTokenService
	Ctx                context.Context
}

func (m *QueryExpressCompanyTask) Run() error {
	wxConfigList := m.WxService.MiniProgram(db.GetDB(m.Ctx))
	for _, config := range wxConfigList {
		Orm := db.GetDB(m.Ctx)
		deliverys, err := m.WxService.GetDeliveryList(m.AccessTokenService.GetAccessToken(config.(*model.WechatConfig)))
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
	return &QueryExpressCompanyTask{Ctx: context}
}
