package collage

import (
	"errors"
	"github.com/nbvghost/dandelion/service"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/dao"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
)

type GoodsAdd struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		GoodsID     uint   `form:"GoodsID"`
		CollageHash string `form:"CollageHash"`
	} `method:"Post"`
}

func (m *GoodsAdd) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	panic("implement me")
}

func (m *GoodsAdd) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	//organization := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)

	//context.Request.ParseForm()

	//GoodsID, _ := strconv.ParseUint(context.Request.FormValue("GoodsID"), 10, 64)
	//GoodsID := object.ParseUint(context.Request.FormValue("GoodsID"))
	//CollageHash := context.Request.FormValue("CollageHash")

	goods := service.Goods.Goods.FindGoodsByOrganizationIDAndGoodsID(m.Organization.ID, dao.PrimaryKey(m.Post.GoodsID))
	if goods.ID == 0 {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("没有找到商品"), "", nil)}, err
	}
	collage := service.Activity.Collage.GetCollageByHash(m.Post.CollageHash, m.Organization.ID)
	if collage.ID == 0 {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("没有找到限时抢购"), "", nil)}, err
	}

	have := service.Activity.Collage.GetCollageGoodsByGoodsID(goods.ID, m.Organization.ID)
	if have.ID > 0 {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("这个商品已被添加为限时抢购"), "", nil)}, err
	}

	//service.ChangeMap(db.Orm(), timeSell.ID, &model.TimeSell{}, map[string]interface{}{})
	err = dao.Create(db.Orm(), &model.CollageGoods{
		CollageHash: collage.Hash,
		GoodsID:     goods.ID,
		Disable:     false,
		OID:         m.Organization.ID,
	})

	return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(err, "Success", goods)}, err
}
