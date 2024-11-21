package viewmodel

import "github.com/nbvghost/dandelion/library/dao"

type GoodsSpecification struct {
	GoodsID         dao.PrimaryKey
	SpecificationID dao.PrimaryKey
	Quantity        uint
}
