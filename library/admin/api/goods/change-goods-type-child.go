package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/library/result"
)

type ChangeGoodsTypeChild struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		*model.GoodsTypeChild
	} `method:"Post"`
}

func (m *ChangeGoodsTypeChild) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (m *ChangeGoodsTypeChild) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	err = service.Goods.GoodsType.ChangeGoodsTypeChild(ctx, m.Organization.ID, m.Post.ID, m.Post.GoodsTypeChild.Name, m.Post.GoodsTypeChild.Image)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}
