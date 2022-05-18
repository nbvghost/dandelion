package extends

import "github.com/nbvghost/gpa/types"

type GoodsAttribute struct {
	ID    types.PrimaryKey
	Name  string
	Value string
}
type GoodsAttributes struct {
	GroupID   types.PrimaryKey
	GroupName string
	Attrs     []GoodsAttribute
}
