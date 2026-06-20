package goods

import (
	"log"

	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"

	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/util"
)

type AddGoods struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		GoodsJSON          string `form:"goods"`
		SpecificationsJSON string `form:"specifications"`
		ParamsJSON         string `form:"params"`
	} `method:"Post"`
}

func (m *AddGoods) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (m *AddGoods) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)

	//context.Request.ParseForm()
	//goodsStr := context.Request.FormValue("goods")
	//specificationsStr := context.Request.FormValue("specifications")
	//paramsStr := context.Request.FormValue("params")

	//var item model.Goods
	item, err := util.JSONToStruct[*model.Goods](m.Post.GoodsJSON)
	if err != nil {
		log.Println(err)
	}

	//content_item.Params = util.StructToJSON(&gps)

	//var specifications []model.Specification
	specifications, err := util.JSONToStruct[[]model.Specification](m.Post.SpecificationsJSON)
	if err != nil {
		log.Println(err)
	}

	item.OID = m.Organization.ID
	tx := db.GetDB(ctx).Begin()
	hasGoods, err := service.Goods.Goods.SaveGoods(ctx, tx, m.Organization.ID, item, specifications)
	if err != nil {
		tx.Rollback()
		as := &result.ActionResult{}
		as.Code = -55
		as.Message = err.Error()
		as.Data = map[string]interface{}{"Goods": hasGoods}
		return &result.JsonResult{Data: as}, nil
	}
	tx.Commit()

	as := &result.ActionResult{}
	as.Message = "添加成功"
	as.Data = map[string]interface{}{"Goods": item}
	return &result.JsonResult{Data: as}, err
}
