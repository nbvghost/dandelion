package index

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/dandelion/service/order"
)

type VerificationGet struct {
	VerificationService order.VerificationService
	StoreService        company.StoreService
	Get                 struct {
		VerificationNo string `uri:"VerificationNo"`
	} `method:"get"`
}

func (m *VerificationGet) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//VerificationNo := context.PathParams["VerificationNo"]
	verification := m.VerificationService.GetVerificationByVerificationNo(m.Get.VerificationNo)

	if verification.StoreID > 0 && verification.Quantity > 0 {
		//var store model.Store
		store := dao.GetByPrimaryKey(singleton.Orm(), entity.Store, verification.StoreID)
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Success, Message: "", Data: store}}, nil
	} else {
		return &result.JsonResult{Data: &result.ActionResult{Code: result.Fail, Message: "", Data: nil}}, nil
	}

}
