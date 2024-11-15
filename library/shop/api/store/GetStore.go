package store

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type GetStore struct {
	Store *model.Store `mapping:""`
	Get   struct {
		StoreID dao.PrimaryKey `uri:"StoreID"`
	} `method:"Get"`
}

func (g *GetStore) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//StoreID, _ := strconv.ParseUint(context.PathParams["StoreID"], 10, 64)

	//var Store model.Store
	Store := dao.GetByPrimaryKey(db.Orm(), entity.Store, g.Get.StoreID)
	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: Store}}, nil

}
