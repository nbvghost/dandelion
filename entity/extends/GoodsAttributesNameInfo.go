package extends

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type GoodsAttributesNameInfo struct {
	Name string
	Num  uint
}
type GoodsAttributesValueInfo struct {
	Value string
	Num   uint
}

func (m *GoodsAttributesValueInfo) TableName() string {
	return "GoodsAttributesValueInfo"
}

func (m *GoodsAttributesValueInfo) IsZero() bool {
	return len(m.Value) == 0
}

func (m *GoodsAttributesValueInfo) Primary() dao.PrimaryKey {
	return 0
}
