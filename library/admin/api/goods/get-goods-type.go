package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type GetGoodsType struct {
	Organization *model.Organization `mapping:""`
	Get          struct {
		ID   uint   `form:"ID"`
		Name string `form:"Name"`
	} `method:"get"`
}

func (g *GetGoodsType) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	where := make(map[string]any)
	where["OID"] = g.Organization.ID
	if g.Get.ID > 0 {
		where["ID"] = g.Get.ID
	}
	if len(g.Get.Name) > 0 {
		where["Name"] = g.Get.Name
	}
	goods := dao.GetBy(db.Orm(), entity.GoodsType, where)
	return result.NewData(goods), err
}
