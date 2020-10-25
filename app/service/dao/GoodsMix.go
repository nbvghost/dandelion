package dao

type GoodsTypeGoodsTypeChild struct {
	GoodsType      `gorm:"column:GoodsType"`
	GoodsTypeChild `gorm:"column:GoodsTypeChild"`
}

type GoodsTypeItemSub struct {
	Item    *GoodsTypeChild
	SubType []*GoodsTypeItemSub
}

func (m *GoodsTypeItemSub) Get(ID uint64) *GoodsTypeItemSub {
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

func (m GoodsTypeItem) Get(ID uint64) *GoodsTypeItemSub {
	for index := range m.SubType {
		if m.SubType[index].Item.ID == ID {
			return m.SubType[index]
		}
	}
	return &GoodsTypeItemSub{Item: &GoodsTypeChild{}, SubType: []*GoodsTypeItemSub{}}
}

type GoodsTypeData struct {
	List  []*GoodsTypeItem
	ID    uint64
	SubID uint64

	Top *GoodsTypeItem
	Sub *GoodsTypeItemSub
}

func (m *GoodsTypeData) SetCurrentMenus(ID, SubID uint64) {
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
func (m *GoodsTypeData) Get(ID uint64) *GoodsTypeItem {
	for index := range m.List {
		if m.List[index].Item.ID == ID {
			return m.List[index]
		}
	}
	return &GoodsTypeItem{Item: &GoodsType{}, SubType: []*GoodsTypeItemSub{}}
}
