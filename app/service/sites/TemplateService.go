package sites

import (
	"errors"
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/conf"
	"html/template"
)

type TemplateService struct {
	Content service.ContentService
}

func (service TemplateService) CommonTemplate(context *gweb.Context, params map[string]interface{}) string {
	siteName := context.PathParams["siteName"]

	return "/sites/" + siteName + "/template/common/*"
}
func (service TemplateService) IndexTemplate(context *gweb.Context) (map[string]interface{}, *template.Template) {
	siteName := context.PathParams["siteName"]
	allTemplates := template.Must(template.ParseGlob(conf.Config.ViewDir + "/sites/" + siteName + "/template/Menus"))

	return nil, allTemplates
}
func (service TemplateService) MenusTemplate(context *gweb.Context, params map[string]interface{}) string {
	siteName := context.PathParams["siteName"]

	org := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	subTypes := service.Content.FindAllContentSubType(org.ID)

	menus := make([]map[string]interface{}, 0)
	for index := range subTypes {

		item := subTypes[index]

		/*_, ok := menus[item.ContentItem.ID]
		if ok == false {
			menus[item.ContentItem.ID] = map[string]interface{}{
				"Item":    item.ContentItem,
				"SubType": make(map[uint64]interface{}),
			}
		}*/

		var topMenus map[string]interface{}

		var has = false
		for index := range menus {
			sitem := menus[index]["Item"].(dao.ContentItem)
			if sitem.ID == item.ContentItem.ID {
				topMenus = menus[index]
				has = true
				break
			}

		}

		if !has {
			topMenus = map[string]interface{}{
				"Item":    item.ContentItem,
				"SubType": make(map[uint64]interface{}),
			}

			menus = append(menus, topMenus)
		}

		if item.ContentSubType.ParentContentSubTypeID == 0 && item.ContentSubType.ID > 0 {
			_, ok := topMenus["SubType"].(map[uint64]interface{})[item.ContentSubType.ID]
			if ok == false {
				topMenus["SubType"].(map[uint64]interface{})[item.ContentSubType.ID] = map[string]interface{}{
					"Item":    item.ContentSubType,
					"SubType": make(map[uint64]interface{}),
				}
			}

			topMenus["SubType"].(map[uint64]interface{})[item.ContentSubType.ID].(map[string]interface{})["Item"] = item.ContentSubType

		} else if item.ContentSubType.ID > 0 {
			_, ok := topMenus["SubType"].(map[uint64]interface{})[item.ContentSubType.ParentContentSubTypeID]
			if ok == false {
				topMenus["SubType"].(map[uint64]interface{})[item.ContentSubType.ParentContentSubTypeID] = map[string]interface{}{
					"SubType": make(map[uint64]interface{}),
				}
			}

			topMenus["SubType"].(map[uint64]interface{})[item.ContentSubType.ParentContentSubTypeID].(map[string]interface{})["SubType"].(map[uint64]interface{})[item.ContentSubType.ID] = item.ContentSubType
		}

	}

	key := "Menus"
	if _, ok := params[key]; !ok {
		params[key] = menus
	} else {
		panic(errors.New("参数名冲突:" + key))
	}

	return "/sites/" + siteName + "/template/Menus.html"
}
