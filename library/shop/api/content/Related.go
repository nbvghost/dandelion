package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Related struct {
	Get struct {
		Offset           int            `form:"Offset"`
		ContentItemID    dao.PrimaryKey `form:"ContentItemID"`
		ContentSubTypeID dao.PrimaryKey `form:"ContentSubTypeID"`
	} `method:"Get"`
}

func (g *Related) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//TODO implement me
	panic("implement me")
}

func (g *Related) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	//Offset, _ := strconv.Atoi(context.Request.URL.Query().Get("Offset"))
	//ContentItemID, _ := strconv.ParseUint(context.PathParams["ContentItemID"], 10, 64)
	//ContentSubTypeID, _ := strconv.ParseUint(context.PathParams["ContentSubTypeID"], 10, 64)

	//var articles []entity.Content
	//controller.Content.FindOrderWhereLength(entity.Orm(),"Look desc",&articles,)
	pagin := service.Content.FindSelectWherePaging(db.GetDB(ctx), "ID,Title,Picture,ContentItemID,ContentSubTypeID,Author,Look,FromUrl", "CreatedAt desc", model.Content{}, g.Get.Offset, "ContentItemID=? and ContentSubTypeID=?", g.Get.ContentItemID, g.Get.ContentSubTypeID)
	return &result.JsonResult{Data: &pagin}, nil
}
