package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type ReviewDetails struct {
	User *model.User `mapping:""`
	Get  struct {
		ID dao.PrimaryKey `form:"id"`
	} `method:"get"`
}
type GoodsReviewDetails struct {
	model.GoodsReview `json:"GoodsReview"`
	model.Goods       `json:"Goods"`
	ReviewCount       int64                `json:"ReviewCount"`
	GoodsRating       *extends.GoodsRating `gorm:"-"`
}

func (m *ReviewDetails) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	var item GoodsReviewDetails
	db.Orm().Table(`"GoodsReview"`).
		Select(`"GoodsReview".*,"Goods".*`).
		Joins(`JOIN "Goods" on "Goods"."ID" = "GoodsReview"."GoodsID"`).
		Where(`"GoodsReview"."UserID"=? and "GoodsReview"."ID"=?`, context.UID(), m.Get.ID).First(&item)

	item.ReviewCount = dao.Find(db.Orm(), &model.GoodsReview{}).Where(`"GoodsID"=?`, item.GoodsReview.GoodsID).Count()

	item.GoodsRating = service.Goods.Goods.Rating(item.Goods.ID)

	return result.NewData(item), nil
}
