package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type VerificationGet struct {
	Get struct {
		VerificationNo string `uri:"VerificationNo"`
	} `method:"get"`
}

func (m *VerificationGet) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//VerificationNo := context.PathParams["VerificationNo"]
	verification := service.Order.Verification.GetVerificationByVerificationNo(m.Get.VerificationNo)

	if verification.StoreID > 0 && verification.Quantity > 0 {
		//var store model.Store
		store := dao.GetByPrimaryKey(db.Orm(), entity.Store, verification.StoreID)
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: store}}, nil
	} else {
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: "", Data: nil}}, nil
	}

}
