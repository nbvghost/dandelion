package order

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/entity/model"
)

type Query struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		PageSize  int
		Page      int
		Status    []model.OrdersStatus
		StartDate string
		EndDate   string
	} `method:"post"`
}

func (m *Query) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Query) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	/*startDate, err := time.ParseInLocation("2006-01-02 15:04:05", m.Post.StartDate, time.Local)
	if err != nil {
		return nil, err
	}
	endDate, err := time.ParseInLocation("2006-01-02 15:04:05", m.Post.EndDate, time.Local)
	if err != nil {
		return nil, err
	}*/

	list, err := service.Order.Orders.ListOrders(nil, m.Organization.ID, (&dao.Sort{}).OrderByColumn(`"CreatedAt"`, true), m.Post.Page, m.Post.PageSize)
	return result.NewData(map[string]any{
		"Pagination": list,
	}), err
}
