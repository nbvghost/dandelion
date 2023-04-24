package order

import (
	"strings"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/order"
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
	if m.Get.PageSize == 0 {
		m.Get.PageSize = 10
	}
	list, totalRecords := m.OrdersService.ListOrders(m.User.ID, m.User.OID, 0, StatusList, m.Get.PageSize, m.Get.PageSize*m.Get.Index)
	return result.NewData(map[string]any{"List": list, "Total": totalRecords, "Index": m.Get.Index}), nil
}
