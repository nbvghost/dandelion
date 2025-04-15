package specification

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/oss"
	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/tool/object"
	"github.com/samber/lo"
	"strings"
)

type Index struct {
	OIDMapping *entity.SessionMappingData `mapping:""`
	Get        struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"get"`
	Delete struct {
		ID dao.PrimaryKey `form:"ID"`
	} `method:"Delete"`
	Post struct {
		Specification struct {
			ID          dao.PrimaryKey
			GoodsID     dao.PrimaryKey
			Label       string
			CodeNo      string
			CodeHS      string
			Num         uint
			Unit        string
			Weight      float64
			Stock       uint
			CostPrice   float64
			MarketPrice float64
			Currency    model.Currency
			Brokerage   float64
			Pictures    sqltype.Array[sqltype.Image]
			Language    model.SpecificationLanguage
			Remark      string
		}
	} `method:"Post"`
	Put struct {
		Specification struct {
			ID          dao.PrimaryKey
			GoodsID     dao.PrimaryKey
			Label       string
			CodeNo      string
			CodeHS      string
			Num         uint
			Unit        string
			Weight      float64
			Stock       uint
			CostPrice   float64
			MarketPrice float64
			Brokerage   float64
			Pictures    sqltype.Array[sqltype.Image]
			Language    model.SpecificationLanguage
			Remark      string
		}
	} `method:"Put"`
}

func (g *Index) Handle(context constrain.IContext) (r constrain.IResult, err error) {

	return result.NewData(dao.GetByPrimaryKey(db.Orm(), &model.Specification{}, g.Get.ID)), err
}
func (g *Index) HandlePut(context constrain.IContext) (r constrain.IResult, err error) {
	has := dao.GetByPrimaryKey(db.Orm(), &model.Specification{}, g.Put.Specification.ID).(*model.Specification)
	if has.IsZero() {
		return nil, errors.New("数据错误")
	}

	s := make(map[string]any)
	if g.Put.Specification.Label != has.Label {
		s["Label"] = g.Put.Specification.Label
	}
	/*if len(g.Put.Name) > 0 {
		s["Name"] = g.Put.Name
	}*/
	/*if len(g.Put.Specification.LabelIndex) > 0 {
		s["LabelIndex"] = g.Put.Specification.LabelIndex
	}*/
	if g.Put.Specification.Num != has.Num {
		s["Num"] = g.Put.Specification.Num
	}
	if g.Put.Specification.CodeNo != has.CodeNo {
		s["CodeNo"] = g.Put.Specification.CodeNo
	}
	if g.Put.Specification.Unit != has.Unit {
		s["Unit"] = g.Put.Specification.Unit
	}

	s["Weight"] = object.ParseUint(object.Decimal((g.Put.Specification.Weight)*1000.0, 0))

	if g.Put.Specification.Stock != has.Stock {
		s["Stock"] = g.Put.Specification.Stock
	}

	s["CodeHS"] = g.Put.Specification.CodeHS

	s["CostPrice"] = object.ParseUint(object.Decimal((g.Put.Specification.CostPrice)*100.0, 0))

	s["MarketPrice"] = object.ParseUint(object.Decimal(g.Put.Specification.MarketPrice*100.0, 0))

	s["Brokerage"] = object.ParseUint(object.Decimal(g.Put.Specification.Brokerage*100.0, 0))

	if g.Put.Specification.Remark != has.Remark {
		s["Remark"] = g.Put.Specification.Remark
	}
	if g.Put.Specification.Language.Label != has.Language.Label {
		s["Language"] = g.Put.Specification.Language
	}

	if len(g.Put.Specification.Pictures) > 0 {
		for i := 0; i < len(g.Put.Specification.Pictures); i++ {
			item := g.Put.Specification.Pictures[i]
			if strings.HasPrefix(item.Src, oss.TempFilePrefix) {
				fileBytes, err := oss.GetTempFile(item.Src)
				if err != nil {
					//tx.Rollback()
					return nil, err
				}
				file, err := oss.UploadFile(context, fileBytes, fmt.Sprintf("goods/specification/%d", g.Put.Specification.ID), "", true, item.Alt)
				if err != nil {
					//tx.Rollback()
					return nil, err
				}
				if file.Code != 0 {
					return nil, errors.New(file.Message)
				}
				g.Put.Specification.Pictures[i].Src = file.Data.Path
			}
		}
	}
	s["Pictures"] = g.Put.Specification.Pictures

	err = dao.UpdateByPrimaryKey(db.Orm(), &model.Specification{}, g.Put.Specification.ID, s)
	if err != nil {
		return nil, err
	}
	return result.NewDataMessage(map[string]any{"Specifications": dao.Find(db.Orm(), &model.Specification{}).Where(`"GoodsID"=?`, g.Put.Specification.GoodsID).Order(`"LabelIndex"::text asc`).List()}, "修改成功"), err
}
func (g *Index) HandlePost(context constrain.IContext) (r constrain.IResult, err error) {
	if g.Post.Specification.GoodsID == 0 {
		return nil, errors.New("数据错误")
	}

	tx := db.Orm().Begin()

	specificationList := dao.Find(tx, &model.Specification{}).Where(map[string]any{"GoodsID": g.Post.Specification.GoodsID}).List()

	/*for i := range specificationList {
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
	}*/

	specification := g.Post.Specification

	hasItem, has := lo.Find(specificationList, func(item dao.IEntity) bool {
		s := item.(*model.Specification)
		if strings.EqualFold(s.Label, specification.Label) {
			return true
		} else {
			return false
		}
	})
	if has {
		return nil, errors.New(fmt.Sprintf("存在相同的规格:%s", hasItem.(*model.Specification).Label))
	}

	currency := g.Post.Specification.Currency
	if len(currency) == 0 {
		currency = "CNY"
	}

	newSpecification := &model.Specification{
		OID:     g.OIDMapping.OID,
		GoodsID: g.Post.Specification.GoodsID,
		//Name:        item.Name,
		Label: specification.Label,
		//LabelIndex:  specification.LabelIndex,
		Num:         specification.Num,
		Unit:        specification.Unit,
		CodeNo:      specification.CodeNo,
		CodeHS:      specification.CodeHS,
		Weight:      object.ParseUint(specification.Weight * 1000.0),
		Stock:       specification.Stock,
		CostPrice:   object.ParseUint(specification.CostPrice * 100.0),
		MarketPrice: object.ParseUint(specification.MarketPrice * 100.0),
		Brokerage:   object.ParseUint(specification.Brokerage * 100.0),
		Language:    specification.Language,
		Remark:      specification.Remark,
		Currency:    currency,
	}
	err = dao.Create(tx, newSpecification)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for i := 0; i < len(specification.Pictures); i++ {
		item := specification.Pictures[i]
		if strings.HasPrefix(item.Src, oss.TempFilePrefix) {
			fileBytes, err := oss.GetTempFile(item.Src)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			file, err := oss.UploadFile(context, fileBytes, fmt.Sprintf("goods/specification/%d", newSpecification.ID), "", true, item.Alt)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			specification.Pictures[i].Src = file.Data.Path
		}
	}
	err = dao.UpdateByPrimaryKey(tx, &model.Specification{}, newSpecification.ID, map[string]any{"Pictures": specification.Pictures})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return result.NewDataMessage(map[string]any{"Specifications": dao.Find(db.Orm(), &model.Specification{}).Where(`"GoodsID"=?`, g.Post.Specification.GoodsID).Order(`"LabelIndex"::text asc`).List()}, "添加成功"), err
}
func (m *Index) HandleDelete(context constrain.IContext) (r constrain.IResult, err error) {
	if m.Delete.ID == 0 {
		return nil, errors.New("没找到记录")
	}
	has := dao.GetByPrimaryKey(db.Orm(), &model.Specification{}, m.Delete.ID).(*model.Specification)
	err = service.Goods.Specification.DeleteSpecification(has.ID)
	if err != nil {
		return nil, err
	}
	return result.NewDataMessage(map[string]any{"Specifications": dao.Find(db.Orm(), &model.Specification{}).Where(`"GoodsID"=?`, has.GoodsID).Order(`"LabelIndex"::text asc`).List()}, "删除成功"), err
}
