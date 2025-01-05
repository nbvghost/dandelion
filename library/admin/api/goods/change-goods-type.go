package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/library/result"
)

type ChangeGoodsType struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		*model.GoodsType
	} `method:"post"`
}

func (g *ChangeGoodsType) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	return nil, err
}
func (g *ChangeGoodsType) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//item := &model.GoodsType{}
	//err = util.RequestBodyToJSON(context.Request.Body, item)
	//log.Println(err)
	err = service.Goods.GoodsType.ChangeGoodsType(g.Organization.ID, g.Post.GoodsType)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}
