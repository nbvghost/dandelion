package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Index struct {
	Get struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"get"`
	Delete struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"Delete"`
}

func (m *Index) Handle(context constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}

func (m *Index) HandleDelete(context constrain.IContext) (constrain.IResult, error) {
	return &result.JsonResult{Data: service.Goods.Goods.DeleteGoods(dao.PrimaryKey(m.Delete.ID))}, nil
}
