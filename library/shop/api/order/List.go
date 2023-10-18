package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/order"
	"strings"
)

type List struct {
	OrdersService order.OrdersService
	User          *model.User `mapping:""`
	Get           struct {
		Status   string `form:"status"`
		Index    int    `form:"index"`
		PageSize int    `form:"page-size"`
	} `method:"get"`
}

func (m *List) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	var StatusList []model.OrdersStatus
	if !strings.EqualFold(m.Get.Status, "") {
		list := strings.Split(m.Get.Status, ",")
		for i := 0; i < len(list); i++ {
			StatusList = append(StatusList, model.OrdersStatus(list[i]))
		}
	}

	params := &order.ListOrdersQueryParam{
		UserID: m.User.ID,
		Status: StatusList,
	}

	list, err := m.OrdersService.ListOrders(params, m.User.OID, (&extends.Order{}).OrderByColumn(`"CreatedAt"`, true), m.Get.Index+1, m.Get.PageSize)
	if err != nil {
		return nil, err
	}
	return result.NewData(map[string]any{"List": list.List, "Total": list.Total, "Index": m.Get.Index}), nil
}
