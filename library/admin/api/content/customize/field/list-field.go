package field

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type ListField struct {
	Organization *model.Organization `mapping:""`

	Put struct {
		GroupID dao.PrimaryKey
	} `method:"put"` //更新
}

func (m *ListField) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return nil, nil
}
func (m *ListField) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {
	d := dao.Find(db.Orm(), &model.CustomizeField{})
	if m.Put.GroupID > 0 {
		d.Where(`"GroupID"=?`, m.Put.GroupID)
	}
	list := d.Order(`"Sort","UpdatedAt" desc`).List()
	return result.NewData(map[string]any{"List": list}), nil
}
