package goods

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
)

type Attributes struct {
	Organization *model.Organization `mapping:""`
	POST         struct {
		GoodsID dao.PrimaryKey `form:"GoodsID"`
		GroupID dao.PrimaryKey `form:"GroupID"`
		Name    string         `form:"Name"`
		Value   string         `form:"Value"`
	} `method:"POST"`

	DELETE struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"DELETE"`
}

func (g *Attributes) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	err = service.Goods.Attributes.AddGoodsAttributes(g.Organization.ID, g.POST.GoodsID, g.POST.GroupID, g.POST.Name, g.POST.Value)
	return &result.JsonResult{Data: result.New(err, "属性添加成功")}, err
}

func (g *Attributes) Handle(ctx constrain.IContext) (constrain.IResult, error) {
	panic("implement me")
}

func (g *Attributes) HandleDelete(context constrain.IContext) (r constrain.IResult, err error) {
	err = service.Goods.Attributes.DeleteGoodsAttributes(g.DELETE.ID)
	return &result.JsonResult{Data: result.New(err, "属性删除成功")}, err
}
