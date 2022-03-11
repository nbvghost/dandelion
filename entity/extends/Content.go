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
