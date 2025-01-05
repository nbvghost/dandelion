package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type List struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		Query    *serviceargument.ListOrdersQueryParam
		Order    dao.Sort
		PageNo   int
		PageSize int
	} `method:"post"`
}

func (m *List) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *List) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)

	//Orm := db.Orm()
	//dts := &model.Datatables{}
	//util.RequestBodyToJSON(context.Request.Body, dts)

	//UserID, _ := strconv.ParseUint(dts.Columns[0].Search.Value, 10, 64)
	//UserID := object.ParseUint(m.Post.Datatables.Columns[0].Search.Value)
	//PostType, _ := strconv.ParseInt(m.Post.Datatables.Columns[1].Search.Value, 10, 64)
	//Status := m.Post.Datatables.Columns[2].Search.Value

	/*var StatusList []model.OrdersStatus
	if !strings.EqualFold(Status, "") {
		list := strings.Split(Status, ",")
		for i := 0; i < len(list); i++ {
			StatusList = append(StatusList, model.OrdersStatus(list[i]))
		}
	}*/
	//fmt.Println(dts)
	d, err := service.Order.Orders.ListOrders(m.Post.Query, m.Organization.ID, m.Post.Order.OrderByColumn(`"CreatedAt"`, true), m.Post.PageNo, m.Post.PageSize)

	return result.NewData(map[string]any{
		"Pagination": d,
	}), err
}
