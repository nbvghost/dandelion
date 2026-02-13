package leave_message

import (
	"fmt"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

//leave-message-list

type List struct {
	Organization *model.Organization `mapping:""`
	Put          struct {
		Query    ListQueryParam
		Order    dao.Sort
		PageNo   int
		PageSize int
	} `method:"put"`
}
type ListQueryParam struct {
	Name     string
	Email    string
	Content  string
	ClientIP string
}

func (m *List) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *List) HandlePut(ctx constrain.IContext) (r constrain.IResult, err error) {

	var pageSize = m.Put.PageSize
	var pageIndex = m.Put.PageNo - 1

	if pageSize <= 0 {
		pageSize = 10
	}

	if pageIndex < 0 {
		pageIndex = 0
	}

	Orm := db.GetDB(ctx)
	var contents []model.LeaveMessage

	db := Orm.Model(model.LeaveMessage{}).Where(`"OID"=?`, m.Organization.ID)

	if len(m.Put.Query.Name) > 0 {
		db = db.Where(`"Name" like ?`, fmt.Sprintf("%%%s%%", m.Put.Query.Name))
	}
	if len(m.Put.Query.Email) > 0 {
		db = db.Where(`"Email" like ?`, fmt.Sprintf("%%%s%%", m.Put.Query.Email))
	}
	if len(m.Put.Query.Content) > 0 {
		db = db.Where(`"Content" like ?`, fmt.Sprintf("%%%s%%", m.Put.Query.Content))
	}
	if len(m.Put.Query.ClientIP) > 0 {
		db = db.Where(`"ClientIP" like ?`, fmt.Sprintf("%%%s%%", m.Put.Query.ClientIP))
	}

	var recordsTotal int64

	db = db.Count(&recordsTotal).Limit(pageSize).Offset(pageSize * pageIndex).Order(m.Put.Order.OrderByColumn(`"CreatedAt"`, true)).Find(&contents)

	return result.NewData(map[string]any{
		"Pagination": result.NewPagination(m.Put.PageNo, pageSize, recordsTotal, contents),
	}), err
}
