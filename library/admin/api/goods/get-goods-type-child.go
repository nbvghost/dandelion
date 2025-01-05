package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type GetGoodsTypeChild struct {
	Organization *model.Organization `mapping:""`
	Get          struct {
		ID          dao.PrimaryKey `form:"ID"`
		Name        string         `form:"Name"`
		GoodsTypeID dao.PrimaryKey `form:"GoodsTypeID"`
	} `method:"get"`
}

func (g *GetGoodsTypeChild) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
	//ID := object.ParseUint(context.Request.URL.Query().Get("ID"))
	//var goods model.GoodsTypeChild
	where := make(map[string]any)
	where["OID"] = g.Organization.ID
	if g.Get.ID > 0 {
		where["ID"] = g.Get.ID
	}
	if len(g.Get.Name) > 0 {
		where["Name"] = g.Get.Name
	}
	if g.Get.GoodsTypeID > 0 {
		where["GoodsTypeID"] = g.Get.GoodsTypeID
	}
	goods := dao.GetBy(db.Orm(), entity.GoodsTypeChild, where)
	return result.NewData(goods), err
}
