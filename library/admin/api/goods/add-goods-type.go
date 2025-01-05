package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type AddGoodsType struct {
	Organization *model.Organization `mapping:""`
	Get          struct {
		*model.GoodsType
	} `method:"Get"`
	Post struct {
		*model.GoodsType
	} `method:"Post"`
}

func (m *AddGoodsType) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return nil, err
}
func (m *AddGoodsType) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	gt, err := service.Goods.GoodsType.AddGoodsType(m.Organization.ID, m.Post.GoodsType)
	return result.NewData(gt), err
}
