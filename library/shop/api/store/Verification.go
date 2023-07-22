package store

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/activity"
	"github.com/nbvghost/dandelion/service/order"
)

type Verification struct {
	VerificationService order.VerificationService
	CardItemService     activity.CardItemService
	Store               *model.Store `mapping:""`
	User                *model.User  `mapping:""`
	Get                 struct {
		VerificationNo string `uri:"VerificationNo"`
	} `method:"Get"`
	Post struct {
		Action         string         `form:"Action"`
		Quantity       uint           `form:"Quantity"`
		VerificationNo string         `form:"VerificationNo"`
		StoreStockID   dao.PrimaryKey `form:"StoreStockID"`
	} `method:"Post"`
}

func (g *Verification) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {

	switch g.Post.Action {
	case "User":
		//核销卡卷

		//verification := controller.Verification.GetVerificationByVerificationNo(VerificationNo)
		tx := db.Orm().Begin()
		err := g.VerificationService.VerificationCardItem(tx, g.Post.VerificationNo, g.Post.Quantity, g.User, g.Store)
		if err != nil {
			tx.Rollback()
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "", nil)}, nil
		} else {
			tx.Commit()
			return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(nil, "核销成功", nil)}, nil
		}
		//fmt.Println(verification)
	case "Self":

		//StoreStockID, _ := strconv.ParseUint(context.Request.FormValue("StoreStockID"), 10, 64)
		//Quantity, _ := strconv.ParseUint(context.Request.FormValue("Quantity"), 10, 64)

		as := g.VerificationService.VerificationSelf(g.Store.ID, g.Post.StoreStockID, g.Post.Quantity)
		return &result.JsonResult{Data: as}, nil

	}
	return &result.JsonResult{Data: &result.ActionResult{}}, nil
}

func (g *Verification) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	verification := g.VerificationService.GetVerificationByVerificationNo(g.Get.VerificationNo)

	//var cardItem model.CardItem
	cardItem := dao.GetByPrimaryKey(db.Orm(), entity.CardItem, verification.CardItemID)

	if verification.ID == 0 {
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: "", Data: nil}}, nil
	}

	results := make(map[string]interface{})
	results["CardItem"] = cardItem
	results["Verification"] = verification

	return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: results}}, nil

}
