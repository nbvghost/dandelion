package serviceargument

import (
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
)

type SpecificationInfo struct {
	Label       string
	LabelIndex  sqltype.Array[dao.PrimaryKey]
	Num         uint
	Weight      uint
	Stock       uint
	CostPrice   uint
	MarketPrice uint
	Brokerage   uint
}
