package serviceargument

import (
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"reflect"
)

type SiteData[T ListType] struct {
	AllMenusData    extends.MenusData //包含被隐藏的菜单
	MenusData       extends.MenusData
	PageMenus       []extends.Menus
	CurrentMenuData CurrentMenuData
	CurrentMenu     extends.Menus
	ContentItem     *model.ContentItem
	ContentItemMap  map[dao.PrimaryKey]model.ContentItem
	ContentSubType  model.ContentSubType
	Pagination      Pagination[T]
	Tags            []extends.Tag
	Navigations     []extends.Menus
	LeftRight       [2]T
	SiteAuthor      string
	ContentConfig   model.ContentConfig
	Organization    model.Organization
	Item            T
	TypeNameMap     map[dao.PrimaryKey]extends.Menus
	Options         Options
}

func (m *SiteData[T]) ListFirst() T {
	var modelContent T
	if m.Pagination.Total > 0 {
		modelContent = m.Pagination.List[0]
	}
	if modelContent == nil {
		modelContent = reflect.New(reflect.TypeOf(modelContent).Elem()).Interface().(T)
	}
	return modelContent
}
