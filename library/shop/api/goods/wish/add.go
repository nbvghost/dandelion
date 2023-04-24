package wish

import (
	"errors"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/activity"
	"github.com/nbvghost/gpa/types"
)

type Add struct {
	CardItemService activity.CardItemService
	User            *model.User `mapping:""`
	Post            struct {
		GoodsID         types.PrimaryKey
		SpecificationID types.PrimaryKey
		Quantity        uint
		Comment         string
	} `method:"Post"`
}

func (m *Add) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	return nil, nil
}
func (m *Add) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	if m.Post.GoodsID == 0 || m.Post.SpecificationID == 0 {
		return nil, errors.New("params error")
	}
	if m.Post.Quantity == 0 {
		m.Post.Quantity = 1
	}

	wish := &model.GoodsWish{
		GoodsID:         m.Post.GoodsID,
		SpecificationID: m.Post.SpecificationID,
		Quantity:        m.Post.Quantity,
		Comment:         m.Post.Comment,
	}

	has := dao.GetBy(singleton.Orm(), &model.GoodsWish{}, map[string]any{"GoodsID": wish.GoodsID, "SpecificationID": wish.SpecificationID}).(*model.GoodsWish)
	if has.IsZero() == false {
		has.Quantity = has.Quantity + m.Post.Quantity
		err := dao.UpdateByPrimaryKey(singleton.Orm(), &model.GoodsWish{}, has.ID, map[string]any{"Quantity": has.Quantity})
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	err := dao.Create(singleton.Orm(), wish)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
