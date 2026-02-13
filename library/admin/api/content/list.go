package content

import (
	"fmt"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

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
	Title            string
	Content          string
	ContentItemID    dao.PrimaryKey
	ContentSubTypeID dao.PrimaryKey
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
	var contents []model.Content

	db := Orm.Model(model.Content{})

	db = db.Where(`"ContentItemID"=?`, m.Put.Query.ContentItemID)

	if len(m.Put.Query.Title) > 0 {
		db = db.Where(fmt.Sprintf(`"Title" like '%%%s%%'`, m.Put.Query.Title))
	}
	if len(m.Put.Query.Content) > 0 {
		db = db.Where(fmt.Sprintf(`"Content" like '%%%s%%'`, m.Put.Query.Content))
	}

	if m.Put.Query.ContentSubTypeID > 0 {
		db = db.Where(`"ContentSubTypeID"=?`, m.Put.Query.ContentSubTypeID)
	}

	var recordsTotal int64

	db = db.Count(&recordsTotal).Limit(pageSize).Offset(pageSize * pageIndex).Order(`"IsStickyTop"`).Order(m.Put.Order.OrderByColumn(`"CreatedAt"`, true)).Find(&contents)

	return result.NewData(map[string]any{
		"Pagination": result.NewPagination(m.Put.PageNo, pageSize, recordsTotal, contents),
	}), err
}
