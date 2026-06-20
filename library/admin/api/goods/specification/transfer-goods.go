package specification

import (
	"errors"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type TransferGoods struct {
	OIDMapping *entity.SessionMappingData `mapping:""`
	Put        struct {
		ID         dao.PrimaryKey
		NewGoodsID dao.PrimaryKey
	} `method:"Put"`
}

func (g *TransferGoods) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	return nil, nil
}

func (g *TransferGoods) HandlePut(ctx constrain.IContext) (constrain.IResult, error) {
	specification := service.Goods.Specification.GetSpecification(ctx, g.Put.ID)
	if specification.IsZero() {
		return nil, errors.New("找不到原数据")
	}
	newGoods := dao.GetByPrimaryKey(ctx.GetDB(), &model.Goods{}, g.Put.NewGoodsID).(*model.Goods)
	if newGoods.IsZero() {
		return nil, errors.New("找不到产品数据")
	}

	if specification.OID != g.OIDMapping.OID || newGoods.OID != g.OIDMapping.OID {
		return nil, errors.New("无法修改")
	}
	err := dao.UpdateByPrimaryKey(ctx.GetDB(), &model.Specification{}, specification.ID, map[string]any{"GoodsID": newGoods.ID})
	if err != nil {
		return nil, err
	}
	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "修改成功", nil)}, err
}
