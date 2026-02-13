package goods

import (
	"errors"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"github.com/shopspring/decimal"
)

type Specification struct {
	OIDMapping *entity.SessionMappingData `mapping:""`
	Get        struct {
		GoodsID uint `form:"goods-id"`
	} `method:"get"`
	Delete struct {
		ID dao.PrimaryKey `form:"id"`
	} `method:"Delete"`
	Post struct {
		GoodsID dao.PrimaryKey
		List    []struct {
			Label string
			//Name        string
			LabelIndex  sqltype.Array[dao.PrimaryKey]
			Num         uint
			Weight      decimal.Decimal
			Stock       uint
			CostPrice   decimal.Decimal
			MarketPrice decimal.Decimal
			Brokerage   decimal.Decimal
			Language    model.SpecificationLanguage
		}
	} `method:"Post"`
	Put model.Specification `method:"Put"`
}

func (g *Specification) Handle(ctx constrain.IContext) (r constrain.IResult, err error) {

	return result.NewData(map[string]any{"Specifications": dao.Find(db.GetDB(ctx), &model.Specification{}).Where(`"GoodsID"=?`, g.Get.GoodsID).Order(`"LabelIndex"::text asc`).List()}), err
}
func (g *Specification) HandlePut(ctx constrain.IContext) (r constrain.IResult, err error) {
	if g.Put.ID == 0 {
		return nil, errors.New("数据错误")
	}

	s := make(map[string]any)
	if len(g.Put.Label) > 0 {
		s["Label"] = g.Put.Label
	}
	/*if len(g.Put.Name) > 0 {
		s["Name"] = g.Put.Name
	}*/
	if len(g.Put.LabelIndex) > 0 {
		s["LabelIndex"] = g.Put.LabelIndex
	}
	if g.Put.Num > 0 {
		s["Num"] = g.Put.Num
	}
	if g.Put.Weight.GreaterThan(decimal.NewFromInt(0)) {
		s["Weight"] = g.Put.Weight
	}
	if g.Put.Stock > 0 {
		s["Stock"] = g.Put.Stock
	}
	if g.Put.CostPrice.GreaterThan(decimal.NewFromInt(0)) {
		s["CostPrice"] = g.Put.CostPrice
	}
	if g.Put.MarketPrice.GreaterThan(decimal.NewFromInt(0)) {
		s["MarketPrice"] = g.Put.MarketPrice
	}
	if g.Put.Brokerage.GreaterThan(decimal.NewFromInt(0)) {
		s["Brokerage"] = g.Put.Brokerage
	}
	if len(g.Put.Language.Label) > 0 {
		s["Language"] = g.Put.Language
	}

	err = dao.UpdateByPrimaryKey(db.GetDB(ctx), &model.Specification{}, g.Put.ID, s)
	if err != nil {
		return nil, err
	}
	return result.NewData(map[string]any{"Specifications": dao.Find(db.GetDB(ctx), &model.Specification{}).Where(`"GoodsID"=?`, g.Put.GoodsID).Order(`"LabelIndex"::text asc`).List()}), err
}
func (g *Specification) HandlePost(ctx constrain.IContext) (r constrain.IResult, err error) {
	if g.Post.GoodsID == 0 {
		return nil, errors.New("数据错误")
	}

	tx := db.GetDB(ctx).Begin()

	specificationList := dao.Find(tx, &model.Specification{}).Where(map[string]any{"GoodsID": g.Post.GoodsID}).List()

	for i := range specificationList {
		item := specificationList[i].(*model.Specification)
		var hasIndex = false
		for ii := range g.Post.List {
			items := g.Post.List[ii]
			if item.Label == items.Label {
				hasIndex = true
				break
			}
		}
		if !hasIndex {
			err = dao.DeleteByPrimaryKey(tx, &model.Specification{}, item.ID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	for i := range g.Post.List {
		item := g.Post.List[i]
		var hasIndex = -1
		for ii := range specificationList {
			items := specificationList[ii].(*model.Specification)
			if item.Label == items.Label {
				hasIndex = ii
				break
			}
		}
		if hasIndex == -1 {
			err = dao.Create(tx, &model.Specification{
				OID:     g.OIDMapping.OID,
				GoodsID: g.Post.GoodsID,
				//Name:        item.Name,
				Label:       item.Label,
				LabelIndex:  item.LabelIndex,
				Num:         item.Num,
				Unit:        "",
				Weight:      item.Weight,
				Stock:       item.Stock,
				CostPrice:   item.CostPrice,
				MarketPrice: item.MarketPrice,
				Brokerage:   item.Brokerage,
				Pictures:    nil,
				Language:    item.Language,
			})
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	tx.Commit()
	return result.NewData(map[string]any{"Specifications": dao.Find(db.GetDB(ctx), &model.Specification{}).Where(`"GoodsID"=?`, g.Post.GoodsID).Order(`"LabelIndex"::text asc`).List()}), err
}
func (g *Specification) HandleDelete(ctx constrain.IContext) (r constrain.IResult, err error) {
	if g.Delete.ID == 0 {
		return nil, errors.New("没找到记录")
	}
	has := dao.GetByPrimaryKey(db.GetDB(ctx), &model.Specification{}, g.Delete.ID).(*model.Specification)
	err = service.Goods.Specification.DeleteSpecification(ctx, has.ID)
	if err != nil {
		return nil, err
	}
	return result.NewData(map[string]any{"Specifications": dao.Find(db.GetDB(ctx), &model.Specification{}).Where(`"GoodsID"=?`, has.GoodsID).Order(`"LabelIndex"::text asc`).List()}), err
}
