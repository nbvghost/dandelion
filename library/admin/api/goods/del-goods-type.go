package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service"
)

type DelGoodsType struct {
	Get struct {
		ID uint `form:"ID"`
	} `method:"get"`
}

func (m *DelGoodsType) Handle(context constrain.IContext) (constrain.IResult, error) {
	//ID, _ := strconv.ParseUint(context.Request.URL.Query().Get("ID"), 10, 64)
	//ID := object.ParseUint(context.Request.URL.Query().Get("ID"))
	err := service.Goods.GoodsType.DeleteGoodsType(dao.PrimaryKey(m.Get.ID))
	return nil, err
}
