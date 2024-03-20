package express

import (
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

func (b ExpressTemplateService) GetExpressInfo(context constrain.IContext, OrdersID dao.PrimaryKey) (map[string]interface{}, error) {
	waybill, err := b.WxService.GetTraceWaybill(context, OrdersID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"WaybillToken": waybill,
	}, nil
}
func (b ExpressTemplateService) GetExpressTemplateByName(Name string) model.ExpressTemplate {
	Orm := db.Orm()
	var list model.ExpressTemplate
	Orm.Model(&model.ExpressTemplate{}).Where(`"Name"=?`, Name).Find(&list)
	return list
}
func (b ExpressTemplateService) GetExpressTemplateByOID(OID dao.PrimaryKey) model.ExpressTemplate {
	Orm := db.Orm()
	var list model.ExpressTemplate
	Orm.Model(&model.ExpressTemplate{}).Where(`"OID"=?`, OID).First(&list)
	return list
}
func (b ExpressTemplateService) SaveExpressTemplate(target *model.ExpressTemplate) error {
	Orm := db.Orm()
	have := b.GetExpressTemplateByName(target.Name)
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
