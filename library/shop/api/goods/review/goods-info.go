package review

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type GoodsInfo struct {
	User            *model.User `mapping:""`
	Get             struct {
		GoodsID  dao.PrimaryKey `form:"goods-id"`
		Index    int            `form:"index"`
		PageSize int            `form:"page-size"`
	} `method:"get"`
	Post struct {
		Rating   uint
		NickName string
		Title    string
		Content  string
		GoodsID  dao.PrimaryKey
	} `method:"Post"`
}

func (m *GoodsInfo) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	db := dao.Find(db.Orm(), &model.GoodsReview{}).Where(`"GoodsID"=?`, m.Get.GoodsID).Order(`"CreatedAt" desc`)
	total := db.Limit(m.Get.Index, m.Get.PageSize)
	list := db.List()
	return result.NewData(result.NewPagination(m.Get.Index, m.Get.PageSize, int(total), list)), nil
}

func (m *GoodsInfo) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	review := &model.GoodsReview{
		GoodsID:  m.Post.GoodsID,
		UserID:   m.User.ID,
		Title:    m.Post.Title,
		Content:  m.Post.Content,
		NickName: m.Post.NickName,
		Rating:   m.Post.Rating,
	}
	/*has := dao.GetBy(singleton.Orm(), &model.GoodsReview{}, map[string]any{"GoodsID": review.GoodsID, "UserID": review.UserID, "Title": review.Title}).(*model.GoodsReview)
	if has.IsZero() == false {

		return nil, errors.New("same title")
	}*/
	err := dao.Create(db.Orm(), review)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
