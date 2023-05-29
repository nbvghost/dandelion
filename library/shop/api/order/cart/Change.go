package cart

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/viewmodel"
	"github.com/nbvghost/dandelion/service/order"
)

type Change struct {
	ShoppingCartService order.ShoppingCartService
	User                *model.User `mapping:""`
	Post                struct {
		List []viewmodel.GoodsSpecification
	} `method:"post"`
}

func (m *Change) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//context.Request.ParseForm()
	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)

	//GSID, _ := strconv.ParseUint(context.Request.FormValue("GSID"), 10, 64)
	//Quantity, _ := strconv.ParseUint(context.Request.FormValue("Quantity"), 10, 64)
	//GSID := object.ParseUint(context.Request.FormValue("GSID"))
	//Quantity := object.ParseUint(context.Request.FormValue("Quantity"))

	tx := db.Orm().Begin()
	for _, goodsSpecification := range m.Post.List {
		err := m.ShoppingCartService.UpdateByUserIDAndID(tx, m.User.ID, goodsSpecification.GoodsID, goodsSpecification.SpecificationID, uint(goodsSpecification.Quantity))
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return result.NewSuccess("修改成功"), nil
}

func (m *Change) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
