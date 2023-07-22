package goods

import (
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/tool/object"
	"strings"
)

type OptionsType string

func (m OptionsType) String() string {
	return string(m)
}
func NewOptionsType(v string) OptionsType {
	return OptionsType(strings.ToUpper(v))
}

const (
	OptionsTypeAttribute     OptionsType = "ATTRIBUTE"
	OptionsTypeSpecification OptionsType = "SPECIFICATION"
	OptionsTypePackageNum    OptionsType = "PACKAGE_NUM"
	OptionsTypeWeight        OptionsType = "WEIGHT"
	OptionsTypePrice         OptionsType = "PRICE"
)

type OptionValue struct {
	ID    dao.PrimaryKey
	Value string
	Count int
}

func (m OptionValue) Key(Type OptionsType) string {
	return fmt.Sprintf("%s-%d", Type, m.ID)
}

type Option struct {
	Type  OptionsType
	Key   string
	Value []OptionValue
}

type Options struct {
	Attributes []Option
}

func (m *Options) AddAttributes(optionsType OptionsType, id dao.PrimaryKey, key, value string) {
	var has bool
	for i := 0; i < len(m.Attributes); i++ {
		item := m.Attributes[i]
		if strings.EqualFold(item.Key, key) {
			var hasOptionValue bool
			for ii := range item.Value {
				optionValue := item.Value[ii]
				if strings.EqualFold(optionValue.Value, value) {
					m.Attributes[i].Value[ii].Count = m.Attributes[i].Value[ii].Count + 1
					hasOptionValue = true
					break
				}
			}
			if !hasOptionValue {
				m.Attributes[i].Value = append(m.Attributes[i].Value, OptionValue{ID: id, Value: value, Count: 1})
			}
			has = true
			break
		}
	}
	if !has {
		m.Attributes = append(m.Attributes, Option{
			Type:  optionsType,
			Key:   key,
			Value: []OptionValue{{ID: id, Value: value, Count: 1}},
		})
	}
}
func ProductOptions(ctx constrain.IContext, oid dao.PrimaryKey) (*Options, error) {
	var skuService SKUService

	var options = &Options{}
	{
		goodsList := dao.Find(db.Orm(), &model.GoodsAttributes{}).Where(`"OID"=?`, oid).List()
		for i := range goodsList {
			item := goodsList[i].(*model.GoodsAttributes)
			options.AddAttributes(OptionsTypeAttribute, item.ID, item.Name, item.Value)
		}
	}

	{
		skuList := skuService.SkuLabelByOID(db.Orm(), oid)
		for i := range skuList {
			item := skuList[i]
			for i2 := range item.Data {
				itemData := item.Data[i2]
				options.AddAttributes(OptionsTypeSpecification, itemData.ID, item.Label.Label, itemData.Label)
			}
		}
	}

	{
		//select "Num" from "Specification" group by "Num" order by "Num" desc;
		specificationList, err := dao.Find(db.Orm(), &model.Specification{}).Where(`"OID"=?`, oid).Order(`"Num" desc`).Group("Num")
		if err != nil {
			return nil, err
		}
		sList, ok := specificationList.([]uint)
		if !ok {
			return nil, errors.New("error data type")
		}
		for _, u := range sList {
			options.AddAttributes(OptionsTypePackageNum, 0, "packing number", object.ParseString(u))
		}
	}

	{
		specificationWeightList, err := dao.Find(db.Orm(), &model.Specification{}).Where(`"OID"=?`, oid).Order(`"Weight"`).Group("Weight")
		if err != nil {
			return nil, err
		}
		weightList, ok := specificationWeightList.([]uint)
		if !ok {
			return nil, errors.New("error data type")
		}
		for _, u := range weightList {
			options.AddAttributes(OptionsTypeWeight, 0, "weight", object.ParseString(u))
		}
	}
	{
		specificationMarketPriceList, err := dao.Find(db.Orm(), &model.Specification{}).Where(`"OID"=?`, oid).Order(`"MarketPrice"`).Group("MarketPrice")
		if err != nil {
			return nil, err
		}
		priceList, ok := specificationMarketPriceList.([]uint)
		if !ok {
			return nil, errors.New("error data type")
		}
		for _, u := range priceList {
			options.AddAttributes(OptionsTypePrice, 0, "price", object.ParseString(u))
		}
	}
	var attributes []Option
	for i := 0; i < len(options.Attributes); i++ {
		if len(options.Attributes[i].Value) > 1 {
			attributes = append(attributes, options.Attributes[i])
		}
	}
	options.Attributes = attributes
	return options, nil
}
