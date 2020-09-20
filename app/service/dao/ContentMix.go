package dao

type ContentItemContentSubType struct {
	ContentItem    `gorm:"column:ContentItem"`
	ContentSubType `gorm:"column:ContentSubType"`
}

type MenusSub struct {
	Item    ContentSubType
	SubType []MenusSub
}

func (m MenusSub) Get(ID uint64) MenusSub {
	for index := range m.SubType {
		if m.SubType[index].Item.ID == ID {
			return m.SubType[index]
		}
	}
	return MenusSub{}
}

type Menus struct {
	Item    ContentItem
	SubType []MenusSub
}

func (m Menus) Get(ID uint64) MenusSub {
	for index := range m.SubType {
		if m.SubType[index].Item.ID == ID {
			return m.SubType[index]
		}
	}
	return MenusSub{}
}

type MenusData struct {
	List       []Menus
	ID         uint64
	SubID      uint64
	SubChildID uint64

	Top      Menus
	Sub      MenusSub
	SubChild MenusSub
}

func (m *MenusData) SetCurrentMenus(ID, SubID, SubChildID uint64) {
	for index := range m.List {
		if m.List[index].Item.ID == ID {
			m.Top = m.List[index]
			m.Sub = m.Top.Get(SubID)
			m.SubChild = m.Sub.Get(SubChildID)

			m.ID = ID
			m.SubID = SubID
			m.SubChildID = SubChildID
			break
		}
	}

}
func (m MenusData) Get(ID uint64) Menus {
	for index := range m.List {
		if m.List[index].Item.ID == ID {
			return m.List[index]
		}
	}
	return Menus{}
}
