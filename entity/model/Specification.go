package model

import (
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/shopspring/decimal"
)

type SpecificationLanguage struct {
	Label string
}
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyCNY Currency = "CNY"
)

// Specification 商品规格
type Specification struct {
	dao.Entity
	OID     dao.PrimaryKey `gorm:"column:OID;index"`
	GoodsID dao.PrimaryKey `gorm:"column:GoodsID;index"` //
	//Name        string                        `gorm:"column:Name"`                               //
	Label       string                        `gorm:"column:Label"`                              //
	CodeNo      string                        `gorm:"column:CodeNo"`                             //
	CodeHS      string                        `gorm:"column:CodeHS"`                             //
	LabelIndex  sqltype.Array[dao.PrimaryKey] `gorm:"column:LabelIndex;type:JSON"`               //
	Num         uint                          `gorm:"column:Num"`                                //这个规格里面包含多少个小件
	Unit        string                        `gorm:"column:Unit"`                               //单位
	Weight      decimal.Decimal               `gorm:"column:Weight;type:numeric(24,6)"`          //单个规格多少重g,如果Num是多件的话，Weight=Num*(每小件的重量)
	Stock       uint                          `gorm:"column:Stock"`                              //
	CostPrice   decimal.Decimal               `gorm:"column:CostPrice;type:numeric(24,6)"`       //成本价,默认是本国货币
	MarketPrice decimal.Decimal               `gorm:"column:MarketPrice;type:numeric(24,6)"`     //市场价
	Currency    Currency                      `gorm:"column:Currency"`                           //市场价货币
	Brokerage   decimal.Decimal               `gorm:"column:Brokerage;type:numeric(24,6)"`       //总佣金
	Pictures    sqltype.Array[sqltype.Image]  `gorm:"column:Pictures;type:JSON"`                 //规格的多张图片
	Language    SpecificationLanguage         `gorm:"column:Language;serializer:json;type:json"` //其它语言信息
	Remark      string                        `gorm:"column:Remark"`                             //备注
}

func (m *Specification) GetMarketPrice(quantity uint) decimal.Decimal {
	return m.MarketPrice.Mul(decimal.NewFromInt(int64(quantity)))
}
func (Specification) TableName() string {
	return "Specification"
}
