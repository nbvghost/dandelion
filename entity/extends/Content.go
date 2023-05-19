package extends

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/gpa/types"
)

type Menus struct {
	ID           types.PrimaryKey
	Uri          string
	Name         string
	TemplateName string
	Type         model.ContentTypeType
	Introduction string //主类介绍
	Image        string //主类图片
	List         []Menus
}

func NewMenusByContentSubType(contentItem *model.ContentItem, contentSubType *model.ContentSubType) Menus {
	return Menus{
		ID:           contentSubType.ID,
		Uri:          contentSubType.Uri,
		Name:         contentSubType.Name,
		TemplateName: contentItem.TemplateName,
		Type:         contentItem.Type,
		Introduction: contentItem.Introduction,
		Image:        contentItem.Image,
		List:         nil,
	}
}
func NewMenusByContentItem(contentItem *model.ContentItem) Menus {
	return Menus{
		ID:           contentItem.ID,
		Uri:          contentItem.Uri,
		Name:         contentItem.Name,
		TemplateName: contentItem.TemplateName,
		Type:         contentItem.Type,
		Introduction: contentItem.Introduction,
		Image:        contentItem.Image,
		List:         nil,
	}
}

type MenusData struct {
	List []Menus
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
