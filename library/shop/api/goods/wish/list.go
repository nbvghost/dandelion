package wish

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/goods"
	"github.com/nbvghost/gpa/types"
	"log"
)

type List struct {
	GoodsService goods.GoodsService
	User         *model.User `mapping:""`
	Get          struct {
		Index    int `form:"index"`
		PageSize int `form:"page-size"`
	} `method:"get"`
}

func (m *List) Handle(ctx constrain.IContext) (constrain.IResult, error) {

	db := dao.Find(singleton.Orm(), &model.GoodsWish{}).Order(`"CreatedAt" DESC`)
	total := db.Limit(m.Get.Index, m.Get.PageSize)

	var goodsIDs []types.PrimaryKey
	db.Pluck(`"ID"`, &goodsIDs)

	log.Println(goodsIDs, total)

	m.GoodsService.ListGoodsByIDs(goodsIDs)

	return result.NewData(result.Pagination{}), nil
}
