package job

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service/internal/order"
	"github.com/nbvghost/dandelion/service/internal/wechat"
	"github.com/pkg/errors"
	"log"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type QueryOrdersTask struct {
	OrdersService order.OrdersService
	WxService     wechat.WxService
	context       constrain.IServiceContext
}

func (m *QueryOrdersTask) Run() error {
	wxConfigList := m.WxService.MiniProgram(db.Orm())
	for i := range wxConfigList {
		config := wxConfigList[i].(*model.WechatConfig)
		Orm := db.Orm()
		//var ordersList []model.Orders
		ordersList := dao.Find(Orm, entity.Orders).
			Where(`"OID"=?`, config.OID).
			Where(`"Status"=? or "Status"=? or "Status"=?`, model.OrdersStatusOrder, model.OrdersStatusRefund, model.OrdersStatusCancel).List()
		//Where(`"Status"<>? and "Status"<>? and "Status"<>? and "Status"<>?`, model.OrdersStatusOrderOk, model.OrdersStatusCancelOk, model.OrdersStatusDelete, model.OrdersStatusClosed).List()
		//service.FindWhere(Orm, &ordersList, `"Status"<>? and "Status"<>? and "Status"<>? and "Status"<>?`, model.OrdersStatusOrderOk, model.OrdersStatusCancelOk, model.OrdersStatusDelete, model.OrdersStatusClosed)

		for ii := range ordersList {
			orders := ordersList[ii].(*model.Orders)
			err := m.OrdersService.QueryOrdersTask(m.context, orders)
			if err != nil {
				log.Println(errors.WithMessage(err, fmt.Sprintf("订单ID:%s", orders.ID)))
			}
		}
	}
	return nil
}

func NewQueryOrdersTask(context constrain.IServiceContext) Job {
	return &QueryOrdersTask{context: context}
}
