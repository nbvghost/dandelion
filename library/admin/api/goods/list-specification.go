package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type ListSpecification struct {
	Organization *model.Organization `mapping:""`
	Get          struct {
		GoodsID dao.PrimaryKey `form:"GoodsID"`
	} `method:"get"`
}

func (g *ListSpecification) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *ListSpecification) HandleGet(ctx constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//GoodsID, _ := strconv.ParseUint(context.Request.URL.Query().Get("GoodsID"), 10, 64)
	//GoodsID := object.ParseUint(context.Request.URL.Query().Get("GoodsID"))
	//var gts []model.Specification
	gts := dao.Find(db.GetDB(ctx), &model.Specification{}).Where(`"GoodsID"=?`, g.Get.GoodsID).List()
	//err = g.Goods.SpecificationService.FindWhere(db.GetDB(ctx), &gts, g.Organization.ID, model.Specification{GoodsID: dao.PrimaryKey(g.Get.GoodsID)})
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "OK", gts)}, err
}
