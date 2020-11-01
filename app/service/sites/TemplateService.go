package sites

import (
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/service/content"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/gweb"
	"github.com/nbvghost/gweb/conf"
	"html/template"
)

type TemplateService struct {
	Content content.ContentService
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

func (service TemplateService) MenusTemplate(context *gweb.Context) (dao.MenusData, string) {
	siteName := context.PathParams["siteName"]

	org := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)

	subTypes := service.Content.FindAllContentSubType(org.ID)

	menus := make([]dao.Menus, 0)

	for index := range subTypes {

		item := subTypes[index]

		//var topMenus dao.Menus
		var topIndex = -1
		for index := range menus {
			sitem := menus[index].Item
			if sitem.ID == item.ContentItem.ID {
				topIndex = index
				break
			}
		}

		if topIndex == -1 {
			topMenus := dao.Menus{Item: item.ContentItem, SubType: make([]dao.MenusSub, 0)}
			menus = append(menus, topMenus)
			topIndex = len(menus) - 1
		}

		/*if item.ContentSubType.ID == ContentSubTypeID {
			currentSubType = item.ContentSubType
		}*/

		hasIndex := -1
		for index := range menus[topIndex].SubType {

			mItem := menus[topIndex].SubType[index].Item

			ContentSubTypeID := uint64(0)
			if item.ContentSubType.ParentContentSubTypeID == 0 {
				ContentSubTypeID = item.ContentSubType.ID
			} else {
				ContentSubTypeID = item.ContentSubType.ParentContentSubTypeID
			}

			if mItem.ID == ContentSubTypeID {
				hasIndex = index
				break
			}

		}

		if hasIndex == -1 {
			if item.ContentSubType.ParentContentSubTypeID > 0 && item.ContentSubType.ID > 0 {
				menus[topIndex].SubType = append(menus[topIndex].SubType, dao.MenusSub{
					Item:    dao.ContentSubType{BaseModel: dao.BaseModel{ID: item.ContentSubType.ParentContentSubTypeID}},
					SubType: make([]dao.MenusSub, 0),
				})
				hasIndex = len(menus[topIndex].SubType) - 1
			} else if item.ContentSubType.ID > 0 {
				menus[topIndex].SubType = append(menus[topIndex].SubType, dao.MenusSub{
					Item:    item.ContentSubType,
					SubType: make([]dao.MenusSub, 0),
				})
				hasIndex = len(menus[topIndex].SubType) - 1
			}

		}
		if hasIndex < 0 {
			continue
		}
		//第三级菜单
		if item.ContentSubType.ParentContentSubTypeID == 0 && item.ContentSubType.ID > 0 {

			menus[topIndex].SubType[hasIndex].Item = item.ContentSubType

		}
		if item.ContentSubType.ID > 0 && item.ContentSubType.ParentContentSubTypeID > 0 {

			mItem := menus[topIndex].SubType[hasIndex].Item
			if mItem.ID == item.ContentSubType.ParentContentSubTypeID {

				menus[topIndex].SubType[hasIndex].SubType = append(menus[topIndex].SubType[hasIndex].SubType, dao.MenusSub{Item: item.ContentSubType, SubType: make([]dao.MenusSub, 0)})
			}

		}

		/*if item.ContentItem.ID == ContentItemID {
			menusSubIndex = topIndex
			contentItem = item.ContentItem
		}*/

	}

	menusData := dao.MenusData{}
	//menusData.MenusSubIndex = menusSubIndex
	menusData.List = menus
	//menusData.Item = contentItem
	//menusData.CurrentSubType = currentSubType

	return menusData, "/sites/" + siteName + "/template/menus.html"
}
