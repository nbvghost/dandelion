package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type AddGoodsTypeChild struct {
	Organization *model.Organization `mapping:""`
	Get          struct {
		*model.GoodsTypeChild
	} `method:"Get"`
	Post struct {
		*model.GoodsTypeChild
	} `method:"Post"`
}

func (m *AddGoodsTypeChild) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	//item := &model.GoodsTypeChild{}
	//err = util.RequestBodyToJSON(context.Request.Body, item)
	//log.Println(err)
	err = dao.Create(db.GetDB(ctx), m.Get.GoodsTypeChild)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "添加成功", nil)}, err
}
func (m *AddGoodsTypeChild) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	goodsTypeChild, err := service.Goods.Goods.GoodsTypeService.AddGoodsTypeChild(ctx, m.Organization.ID, m.Post.GoodsTypeID, m.Post.GoodsTypeChild.Name, m.Post.GoodsTypeChild.Image)
	return result.NewData(goodsTypeChild), err
}
