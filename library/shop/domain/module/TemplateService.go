package module

type TemplateService struct {
	//Content content.ContentService
}

/*func (service TemplateService) IndexTemplate(context *gweb.Context) (map[string]interface{}, *template.Template) {
	siteName := context.PathParams["siteName"]
	allTemplates := template.Must(template.ParseGlob(conf.Config.ViewDir + "/sites/" + siteName + "/template/Menus"))

	return nil, allTemplates
}*/

/*func (service TemplateService) MenusTemplate(context *gweb.Context) entity.MenusData {
	//siteName := context.PathParams["siteName"]

	org := context.Session.Attributes.Get(play.SessionOrganization).(*entity.Organization)

	subTypes := service.Content.FindAllContentSubType(org.ID)

	menus := make([]entity.Menus, 0)

	for index := range subTypes {

		item := subTypes[index]

		//var topMenus entity.Menus
		var topIndex = -1
		for index := range menus {
			sitem := menus[index].Item
			if sitem.ID == item.ContentItem.ID {
				topIndex = index
				break
			}
		}

		if topIndex == -1 {
			topMenus := entity.Menus{Item: item.ContentItem, SubType: make([]entity.MenusSub, 0)}
			menus = append(menus, topMenus)
			topIndex = len(menus) - 1
		}



		hasIndex := -1
		for index := range menus[topIndex].SubType {

			mItem := menus[topIndex].SubType[index].Item

			ContentSubTypeID := dao.PrimaryKey(0)
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
				menus[topIndex].SubType = append(menus[topIndex].SubType, entity.MenusSub{
					Item:    entity.ContentSubType{BaseModel: entity.BaseModel{ID: item.ContentSubType.ParentContentSubTypeID}},
					SubType: make([]entity.MenusSub, 0),
				})
				hasIndex = len(menus[topIndex].SubType) - 1
			} else if item.ContentSubType.ID > 0 {
				menus[topIndex].SubType = append(menus[topIndex].SubType, entity.MenusSub{
					Item:    item.ContentSubType,
					SubType: make([]entity.MenusSub, 0),
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

				menus[topIndex].SubType[hasIndex].SubType = append(menus[topIndex].SubType[hasIndex].SubType, entity.MenusSub{Item: item.ContentSubType, SubType: make([]entity.MenusSub, 0)})
			}

		}



	}

	menusData := entity.MenusData{}
	//menusData.MenusSubIndex = menusSubIndex
	menusData.List = menus
	//menusData.Item = contentItem
	//menusData.CurrentSubType = currentSubType

	return menusData
}*/
