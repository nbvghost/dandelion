package extends

import "github.com/nbvghost/dandelion/entity/model"

type AllGoodsType struct {
	GoodsType model.GoodsType
	Sub       []model.GoodsTypeChild
	MaxPrice  uint
	MinPrice  uint
}
