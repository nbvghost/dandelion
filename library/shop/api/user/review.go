package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type Review struct {
	User *model.User `mapping:""`
	Get  struct {
		Index    int `form:"index"`
		PageSize int `form:"page-size"`
	} `method:"get"`
}

type GoodsReview struct {
	model.GoodsReview `json:"GoodsReview"`
	model.Goods       `json:"Goods"`
}

func (m *Review) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {
	var list []GoodsReview
	var total int64
	db.GetDB(ctx).Table(`"GoodsReview"`).
		Select(`"GoodsReview".*,"Goods".*`).
		Joins(`JOIN "Goods" on "Goods"."ID" = "GoodsReview"."GoodsID"`).
		Order(`"GoodsReview"."CreatedAt" DESC`).
		Where(`"GoodsReview"."UserID"=?`, ctx.UID()).
		Count(&total).Offset(m.Get.Index * m.Get.PageSize).
		Limit(m.Get.PageSize).Find(&list)

	return result.NewData(result.NewPagination(m.Get.Index, m.Get.PageSize, total, list)), nil
}
