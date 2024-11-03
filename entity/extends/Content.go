package extends

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type Menus struct {
	ParentMenus  *Menus
	ID           dao.PrimaryKey
	Name         string
	TemplateName string
	Type         model.ContentTypeType
	Introduction string //主类介绍
	Image        string //主类图片
	Badge        string //主类图片
	ShowAtMenu   bool
	UrlPath      string
	List         []Menus
}

func NewMenusByContentSubType(contentItem *model.ContentItem, contentSubType *model.ContentSubType) Menus {
	return Menus{
		ID: contentSubType.ID,
		//Uri:          contentSubType.Uri,
		Name:         contentSubType.Name,
		TemplateName: contentItem.TemplateName,
		Type:         contentItem.Type,
		Introduction: contentItem.Introduction,
		Image:        contentItem.Image,
		Badge:        contentItem.Badge,
		List:         nil,
	}
}
func NewMenusByContentItem(contentItem *model.ContentItem) Menus {
	return Menus{
		ID: contentItem.ID,
		//Uri:          contentItem.Uri,
		Name:         contentItem.Name,
		TemplateName: contentItem.TemplateName,
		Type:         contentItem.Type,
		Introduction: contentItem.Introduction,
		Image:        contentItem.Image,
		Badge:        contentItem.Badge,
		List:         nil,
	}
}

type MenusData struct {
	List []Menus
}

func (m MenusData) readCurrentMenus(list []Menus, id dao.PrimaryKey) Menus {
	if len(list) == 0 {
		return Menus{}
	}
	for i := range list {
		if list[i].ID == id {
			return list[i]
		} else {
			menus := m.readCurrentMenus(list[i].List, id)
			if menus.ID == 0 {
				continue
			}
			return menus
		}
	}
	return Menus{}
}
func (m MenusData) GetCurrentMenus(id dao.PrimaryKey) Menus {
	return m.readCurrentMenus(m.List, id)
}
func (m MenusData) ListMenusByType(t model.ContentTypeType) []Menus {
	var menus []Menus
	for index, value := range m.List {
		if value.Type == t {
			menus = append(menus, m.List[index])
		}
	}
	return menus
}
func (m MenusData) GetMenusByType(t model.ContentTypeType) Menus {
	for index, value := range m.List {
		if value.Type == t {
			return m.List[index]
		}
	}
	return Menus{}
}
