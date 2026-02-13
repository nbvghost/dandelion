package express

import (
	"context"
	"errors"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/wechat"
)

type ExpressTemplateService struct {
	WxService wechat.WxService
}

func (b ExpressTemplateService) GetExpressInfo(context constrain.IContext, OrdersID dao.PrimaryKey, OrdersShipping *model.OrdersShipping) (map[string]interface{}, error) {
	waybill, err := b.WxService.GetTraceWaybill(context, OrdersID, OrdersShipping)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"WaybillToken": waybill,
	}, nil
}
func (b ExpressTemplateService) GetExpressTemplateByName(ctx context.Context, Name string) model.ExpressTemplate {
	Orm := db.GetDB(ctx)
	var list model.ExpressTemplate
	Orm.Model(&model.ExpressTemplate{}).Where(`"Name"=?`, Name).Find(&list)
	return list
}
func (b ExpressTemplateService) GetExpressTemplateByOID(ctx context.Context, OID dao.PrimaryKey) model.ExpressTemplate {
	Orm := db.GetDB(ctx)
	var list model.ExpressTemplate
	Orm.Model(&model.ExpressTemplate{}).Where(`"OID"=?`, OID).First(&list)
	return list
}
func (b ExpressTemplateService) SaveExpressTemplate(ctx context.Context, target *model.ExpressTemplate) error {
	Orm := db.GetDB(ctx)
	have := b.GetExpressTemplateByName(ctx, target.Name)
	if have.ID == 0 {
		return dao.Save(Orm, target)
	} else {
		if have.ID == target.ID {
			return dao.Save(Orm, target)
		} else {
			return errors.New("名称已经存在")
		}

	}
}
