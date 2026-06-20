package goods

import (
	"errors"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type GoodsTypeLevel struct {
	OIDMapping entity.SessionMappingData `mapping:""`
	Put        struct {
		ID          dao.PrimaryKey
		NewParentID dao.PrimaryKey
	} `method:"put"`
}

func (g *GoodsTypeLevel) Handle(context constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}
func (g *GoodsTypeLevel) HandlePut(ctx constrain.IContext) (constrain.IResult, error) {
	goodsType := service.Goods.GoodsType.GetGoodsType(ctx, g.Put.ID)
	if goodsType.IsZero() {
		return nil, errors.New("找不到原数据")
	}
	newGoodsType := service.Goods.GoodsType.GetGoodsType(ctx, g.Put.NewParentID)
	if newGoodsType.IsZero() {
		return nil, errors.New("找不到父级数据")
	}

	if goodsType.OID != g.OIDMapping.OID || newGoodsType.OID != g.OIDMapping.OID {
		return nil, errors.New("无法修改")
	}

	err := dao.UpdateByPrimaryKey(ctx.GetDB(), &model.GoodsType{}, goodsType.ID, map[string]any{"ParentID": newGoodsType.ID})
	if err != nil {
		return nil, err
	}

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}
