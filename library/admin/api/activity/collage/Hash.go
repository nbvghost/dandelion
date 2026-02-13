package collage

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Hash struct {
	Get struct {
		Hash string `uri:"Hash"`
	} `method:"Get"`
}

func (m *Hash) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *Hash) HandleGet(ctx constrain.IContext) (r constrain.IResult, err error) {
	//Orm := db.GetDB(ctx)
	//Hash := context.PathParams["Hash"]
	item := service.Activity.Collage.GetItemByHash(ctx, m.Get.Hash)
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "OK", item)}, err
}
