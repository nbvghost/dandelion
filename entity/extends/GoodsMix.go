package extends

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type GoodsTypeGoodsTypeChild struct {
	model.GoodsType      `gorm:"column:GoodsType"`
	model.GoodsTypeChild `gorm:"column:GoodsTypeChild"`
}

type GoodsTypeItemSub struct {
	Item    *model.GoodsTypeChild
	SubType []*GoodsTypeItemSub
}

func (m *GoodsTypeItemSub) Get(ID dao.PrimaryKey) *GoodsTypeItemSub {
	for index := range m.SubType {
		if m.SubType[index].Item.ID == ID {
			return m.SubType[index]
		}
	}
	return &GoodsTypeItemSub{Item: &model.GoodsTypeChild{}, SubType: []*GoodsTypeItemSub{}}
}

type GoodsTypeItem struct {
	Item    *model.GoodsType
	SubType []*GoodsTypeItemSub
}

func (m GoodsTypeItem) Get(ID dao.PrimaryKey) *GoodsTypeItemSub {
	for index := range m.SubType {
		if m.SubType[index].Item.ID == ID {
			return m.SubType[index]
		}
	}
	return &GoodsTypeItemSub{Item: &model.GoodsTypeChild{}, SubType: []*GoodsTypeItemSub{}}
}

type GoodsTypeData struct {
	List  []*GoodsTypeItem
	ID    dao.PrimaryKey
	SubID dao.PrimaryKey

	Top *GoodsTypeItem
	Sub *GoodsTypeItemSub
}

func (m *GoodsTypeData) SetCurrentMenus(ID, SubID dao.PrimaryKey) {
	for index := range m.List {
		if m.List[index].Item.ID == ID {
			m.Top = m.List[index]
			m.Sub = m.Top.Get(SubID)

			m.ID = ID
			m.SubID = SubID
			return
		}
	}

	m.Top = &GoodsTypeItem{Item: &model.GoodsType{}, SubType: []*GoodsTypeItemSub{}}
	m.Sub = &GoodsTypeItemSub{Item: &model.GoodsTypeChild{}, SubType: []*GoodsTypeItemSub{}}
}
func (m *GoodsTypeData) Get(ID dao.PrimaryKey) *GoodsTypeItem {
	for index := range m.List {
		if m.List[index].Item.ID == ID {
			return m.List[index]
		}
	}
	return &GoodsTypeItem{Item: &model.GoodsType{}, SubType: []*GoodsTypeItemSub{}}
}
