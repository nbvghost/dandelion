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
		Status string `form:"Status"`
		Index  int    `form:"Index"`
	} `method:"get"`
}

func (m *List) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*entity.Organization)
	//Status := context.Request.URL.Query().Get("Status")
	//Index, _ := strconv.Atoi(context.Request.URL.Query().Get("Index"))

	var StatusList []model.OrdersStatus
	if !strings.EqualFold(m.Get.Status, "") {
		list := strings.Split(m.Get.Status, ",")
		for i := 0; i < len(list); i++ {
			StatusList = append(StatusList, model.OrdersStatus(list[i]))
		}
	}

	list, _ := m.OrdersService.ListOrders(m.User.ID, m.User.OID, 0, StatusList, 10, 10*m.Get.Index)
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: list}}, nil
	//fullcuts := controller.FullCut.FindOrderByAmountASC(service.Orm)
	//return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "", fullcuts)}

}
