package store

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/gpa/types"
)

type AddStar struct {
	StoreService company.StoreService
	Store        *model.Store `mapping:""`
	Post         struct {
		StoreID types.PrimaryKey `form:"StoreID"`
		Num     uint             `form:"Num"`
	} `method:"Post"`
}

func (g *AddStar) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {

	//StoreID, _ := strconv.ParseUint(context.Request.FormValue("StoreID"), 10, 64)
	//StoreID := object.ParseUint(context.Request.FormValue("StoreID"))
	//Num, _ := strconv.ParseUint(context.Request.FormValue("Num"), 10, 64)
	//Num := object.ParseUint(context.Request.FormValue("Num"))

	//var store model.Store
	store := dao.GetByPrimaryKey(singleton.Orm(), entity.Store, g.Post.StoreID).(*model.Store)
	if g.Post.Num > 5 {
		g.Post.Num = 5
	}
	store.Stars = store.Stars + g.Post.Num

	store.StarsCount = store.StarsCount + 1
	err := dao.UpdateByPrimaryKey(singleton.Orm(), entity.Store, store.ID, &model.Store{Stars: store.Stars, StarsCount: store.StarsCount})
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "评价成功", nil)}, nil
}

func (g *AddStar) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
