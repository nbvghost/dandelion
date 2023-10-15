package job

import (
	"context"
	"fmt"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/pkg/errors"
	"log"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service/order"
	"github.com/nbvghost/dandelion/service/wechat"
)

type QueryOrdersTask struct {
	OrdersService order.OrdersService
	WxService     wechat.WxService
}

func (m *QueryOrdersTask) Run() error {
	wxConfigList := m.WxService.MiniProgram(db.Orm())
	for i := range wxConfigList {
		config := wxConfigList[i].(*model.WechatConfig)
		Orm := db.Orm()
		//var ordersList []model.Orders
		ordersList := dao.Find(Orm, entity.Orders).Where(`"OID"=?`, config.OID).Where(`"Status"<>? and "Status"<>? and "Status"<>? and "Status"<>?`, model.OrdersStatusOrderOk, model.OrdersStatusCancelOk, model.OrdersStatusDelete, model.OrdersStatusClosed).List()
		//service.FindWhere(Orm, &ordersList, `"Status"<>? and "Status"<>? and "Status"<>? and "Status"<>?`, model.OrdersStatusOrderOk, model.OrdersStatusCancelOk, model.OrdersStatusDelete, model.OrdersStatusClosed)

		for ii := range ordersList {
			orders := ordersList[ii].(*model.Orders)
			err := m.OrdersService.QueryOrdersTask(config, orders)
			if err != nil {
				log.Println(errors.WithMessage(err, fmt.Sprintf("订单ID:%s", orders.ID)))
			}
		}
	}
	return nil
}

func NewQueryOrdersTask(context context.Context) Job {
	return &QueryOrdersTask{}
}
