package viewmodel

import "github.com/nbvghost/gpa/types"

type GoodsSpecification struct {
	GoodsID         types.PrimaryKey
	SpecificationID types.PrimaryKey
	Quantity        uint
}
