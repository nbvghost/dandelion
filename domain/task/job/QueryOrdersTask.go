package job

import (
	"context"
	"log"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/dandelion/service/wechat"
)

type QueryOrdersTask struct {
	OrdersService order.OrdersService
	WxService     wechat.WxService
}

func (m *QueryOrdersTask) Run() error {
	wxConfigList := m.WxService.MiniProgram(singleton.Orm())
	for _, config := range wxConfigList {
		Orm := singleton.Orm()
		//var ordersList []model.Orders
		ordersList := dao.Find(Orm, entity.Orders).Where(`"Status"<>? and "Status"<>? and "Status"<>? and "Status"<>?`, model.OrdersStatusOrderOk, model.OrdersStatusCancelOk, model.OrdersStatusDelete, model.OrdersStatusClosed).List()
		//service.FindWhere(Orm, &ordersList, `"Status"<>? and "Status"<>? and "Status"<>? and "Status"<>?`, model.OrdersStatusOrderOk, model.OrdersStatusCancelOk, model.OrdersStatusDelete, model.OrdersStatusClosed)

		for _, orders := range ordersList {
			err := m.OrdersService.QueryOrdersTask(config.(*model.WechatConfig), orders.(*model.Orders))
			if err != nil {
				log.Println(err)
			}
		}

	}
	return nil
}

func NewQueryOrdersTask(context context.Context) Job {
	return &QueryOrdersTask{}
}
