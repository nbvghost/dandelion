package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type AttributesGroup struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
		GoodsID dao.PrimaryKey `form:"GoodsID"`
		Name    string         `form:"Name"`
	} `method:"POST"`
	Put struct {
		ID   dao.PrimaryKey `form:"ID"`
		Name string         `form:"Name"`
	} `method:"Put"`
	Get struct {
		GoodsID dao.PrimaryKey `form:"GoodsID"`
	} `method:"Get"`
	Delete struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"Delete"`
}

func (g *AttributesGroup) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	list := service.Goods.Attributes.ListGoodsAttributesGroupByGoodsID(ctx, g.Get.GoodsID)
	return &result.JsonResult{Data: result.NewData(list)}, nil
}
func (g *AttributesGroup) HandlePost(ctx constrain.IContext) (constrain.IResult, error) {
	gr, err := service.Goods.Attributes.AddGoodsAttributesGroup(ctx, g.Organization.ID, g.POST.GoodsID, g.POST.Name)
	if err != nil {
		return nil, err
	}
	return &result.JsonResult{Data: result.NewData(gr)}, err
}
func (g *AttributesGroup) HandlePut(ctx constrain.IContext) (r constrain.IResult, err error) {
	err = service.Goods.Attributes.ChangeGoodsAttributesGroup(ctx, g.Put.ID, g.Put.Name)
	if err != nil {
		return nil, err
	}
	return &result.JsonResult{Data: result.NewData(nil)}, err
}
func (g *AttributesGroup) HandleDelete(ctx constrain.IContext) (r constrain.IResult, err error) {
	err = service.Goods.Attributes.DeleteGoodsAttributesGroup(ctx, g.Delete.ID)
	if err != nil {
		return nil, err
	}
	return &result.JsonResult{Data: result.NewData(nil)}, err
}
