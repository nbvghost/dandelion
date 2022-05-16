package model

import (
	"github.com/nbvghost/gpa/types"
)

type GoodsTypeGoodsTypeChild struct {
	GoodsType      `gorm:"column:GoodsType"`
	GoodsTypeChild `gorm:"column:GoodsTypeChild"`
}

type GoodsTypeItemSub struct {
	Item    *GoodsTypeChild
	SubType []*GoodsTypeItemSub
}

func (m *GoodsTypeItemSub) Get(ID types.PrimaryKey) *GoodsTypeItemSub {
	for index := range m.SubType {
		if m.SubType[index].Item.ID == ID {
			return m.SubType[index]
		}
	}
	return &GoodsTypeItemSub{Item: &GoodsTypeChild{}, SubType: []*GoodsTypeItemSub{}}
}

type GoodsTypeItem struct {
	Item    *GoodsType
	SubType []*GoodsTypeItemSub
}

func (m GoodsTypeItem) Get(ID types.PrimaryKey) *GoodsTypeItemSub {
	for index := range m.SubType {
		if m.SubType[index].Item.ID == ID {
			return m.SubType[index]
		}
	}
	return &GoodsTypeItemSub{Item: &GoodsTypeChild{}, SubType: []*GoodsTypeItemSub{}}
}

type GoodsTypeData struct {
	List  []*GoodsTypeItem
	ID    types.PrimaryKey
	SubID types.PrimaryKey

	Top *GoodsTypeItem
	Sub *GoodsTypeItemSub
}

func (m *GoodsTypeData) SetCurrentMenus(ID, SubID types.PrimaryKey) {
	for index := range m.List {
		if m.List[index].Item.ID == ID {
			m.Top = m.List[index]
			m.Sub = m.Top.Get(SubID)

			m.ID = ID
			m.SubID = SubID
			return
		}
	}

	m.Top = &GoodsTypeItem{Item: &GoodsType{}, SubType: []*GoodsTypeItemSub{}}
	m.Sub = &GoodsTypeItemSub{Item: &GoodsTypeChild{}, SubType: []*GoodsTypeItemSub{}}
}
func (m *GoodsTypeData) Get(ID types.PrimaryKey) *GoodsTypeItem {
	for index := range m.List {
		if m.List[index].Item.ID == ID {
			return m.List[index]
		}
	}
	return &GoodsTypeItem{Item: &GoodsType{}, SubType: []*GoodsTypeItemSub{}}
}
