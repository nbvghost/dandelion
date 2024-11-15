package cart

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/viewmodel"
	"github.com/nbvghost/dandelion/service"
)

type Delete struct {
	User *model.User `mapping:""`
	Post struct {
		List []viewmodel.GoodsSpecification
	} `method:"post"`
}

func (m *Delete) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	//context.Request.ParseForm()

	//user := context.Session.Attributes.Get(play.SessionUser).(*entity.User)

	//_ShoppingCartIDs := context.Request.FormValue("GSIDs")
	//ShoppingCartIDs := strings.Split(m.Post.GSIDs, ",")

	tx := db.Orm().Begin()
	for _, goodsSpecification := range m.Post.List {
		err := service.Order.ShoppingCart.DeleteListByIDs(tx, m.User.ID, goodsSpecification.GoodsID, goodsSpecification.SpecificationID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return result.NewSuccess("删除成功"), nil
}

func (m *Delete) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	//TODO implement me
	panic("implement me")

}
