package content

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
)

type ISiteData interface {
	SiteData(context constrain.IContext, organization *model.Organization, contentConfig model.ContentConfig, menus extends.Menus, content *model.Content) any
}
