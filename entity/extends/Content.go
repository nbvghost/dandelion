package extends

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/gpa/types"
)

type Menus struct {
	ID           types.PrimaryKey
	Name         string
	TemplateName string
	Type         model.ContentTypeType
	List         []Menus
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
