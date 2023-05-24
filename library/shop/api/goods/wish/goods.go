package wish

import (
	"errors"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/activity"
	"github.com/nbvghost/gpa/types"
)

type Goods struct {
	CardItemService activity.CardItemService
	User            *model.User `mapping:""`
	Get             struct {
		Index    int `form:"index"`
		PageSize int `form:"page-size"`
	} `method:"get"`
	Post struct {
		GoodsID         types.PrimaryKey
		SpecificationID types.PrimaryKey
		Quantity        uint
		Comment         string
	} `method:"Post"`
	Delete struct {
		ID types.PrimaryKey `form:"id"`
	} `method:"Delete"`
}

type GoodsWish struct {
	model.GoodsWish `json:"GoodsWish"`
	model.Goods     `json:"Goods"`
}

func (m *Goods) HandleDelete(ctx constrain.IContext) (constrain.IResult, error) {
	err := dao.DeleteBy(singleton.Orm(), &model.GoodsWish{}, map[string]any{"UserID": ctx.UID(), "ID": m.Delete.ID})
	if err != nil {
		return nil, err
	}
	return result.NewData(nil), nil
}
func (m *Goods) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	var total int64
	var list []GoodsWish
	singleton.Orm().Table("GoodsWish").Select(`"GoodsWish".*,"Goods".*`).Joins(`JOIN "Goods" on "Goods"."ID" = "GoodsWish"."GoodsID"`).
		Order(`"GoodsWish"."CreatedAt" DESC`).Where(`"GoodsWish"."UserID"=?`, ctx.UID()).
		Count(&total).Offset(m.Get.Index * m.Get.PageSize).
		Limit(m.Get.PageSize).Find(&list)

	return result.NewData(result.NewPagination(m.Get.Index, m.Get.PageSize, int(total), list)), nil
}

func (m *Goods) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	if m.Post.GoodsID == 0 || m.Post.SpecificationID == 0 {
		return nil, errors.New("params error")
	}
	if m.Post.Quantity == 0 {
		m.Post.Quantity = 1
	}

	wish := &model.GoodsWish{
		OID:             m.User.OID,
		UserID:          m.User.ID,
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
