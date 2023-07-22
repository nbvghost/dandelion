package extends

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type GoodsAttribute struct {
	ID    dao.PrimaryKey
	Name  string
	Value string
}
type GoodsAttributes struct {
	GroupID   dao.PrimaryKey
	GroupName string
	Attrs     []GoodsAttribute
}
