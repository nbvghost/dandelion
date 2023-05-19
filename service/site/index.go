package site

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/company"
	"github.com/nbvghost/dandelion/service/content"
	"github.com/nbvghost/dandelion/service/goods"
	"github.com/nbvghost/dandelion/service/site/module"
	"github.com/nbvghost/gpa/types"
)

type Service struct {
	GoodsService        goods.GoodsService
	OrganizationService company.OrganizationService
	ContentService      content.ContentService
}

func (service Service) FindShowMenus(OID types.PrimaryKey) extends.MenusData {
	return service.menus(OID, 2)
}
func (service Service) FindAllMenus(OID types.PrimaryKey) extends.MenusData {
	return service.menus(OID, 0)
}
func newRedisContentDataKey(OID types.PrimaryKey, ContentItemUri, ContentSubTypeUri string, pageIndex int) string {
	return fmt.Sprintf("content:%d:%s:%s:%d", OID, ContentItemUri, ContentSubTypeUri, pageIndex)
}
func newRedisGoodsDataKey(OID types.PrimaryKey, ContentItemUri, ContentSubTypeUri string, pageIndex int) string {
	return fmt.Sprintf("goods:%d:%s:%s:%d", OID, ContentItemUri, ContentSubTypeUri, pageIndex)
}
func (service Service) menus(OID types.PrimaryKey, hide uint) extends.MenusData {
	Orm := singleton.Orm()

	var contentItemList []model.ContentItem

	switch hide {
	case 0: //all
		Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
			"OID": OID,
		}).Order(`"Sort"`).Find(&contentItemList)
	case 1: //hide
		Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
			"Hide": true,
			"OID":  OID,
		}).Order(`"Sort"`).Find(&contentItemList)
	case 2: //show
		Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
			"Hide": false,
			"OID":  OID,
		}).Order(`"Sort"`).Find(&contentItemList)
	default:
		Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
			"OID": OID,
		}).Order(`"Sort"`).Find(&contentItemList)

	}

	var contentItemIDs []types.PrimaryKey
	for i := 0; i < len(contentItemList); i++ {
		contentItem := contentItemList[i]
		var have bool
		for ii := 0; ii < len(contentItemIDs); ii++ {
			if contentItem.ID == contentItemIDs[ii] {
				have = true
				break
			}
		}
		if !have {
			contentItemIDs = append(contentItemIDs, contentItem.ID)
		}
	}

	var contentSubTypeList []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(`"ContentItemID" in ?`, contentItemIDs).Order(`"Sort"`).Order(`"ID"`).Find(&contentSubTypeList)

	var goodsTypeList []model.GoodsType
	Orm.Model(model.GoodsType{}).Where(`"OID"=?`, OID).Order(`"ID"`).Find(&goodsTypeList)

	var goodsTypeChildList []model.GoodsTypeChild
	Orm.Model(model.GoodsTypeChild{}).Where(`"OID" = ?`, OID).Order(`"ID"`).Find(&goodsTypeChildList)

	var menusData extends.MenusData

	list := []extends.Menus{}
	for i := 0; i < len(contentItemList); i++ {
		contentItem := contentItemList[i]
		menussddddd := extends.Menus{
			ID:           contentItem.ID,
			Uri:          contentItem.Uri,
			Name:         contentItem.Name,
			TemplateName: contentItem.TemplateName,
			Type:         contentItem.Type,
			Introduction: contentItem.Introduction,
			Image:        contentItem.Image,
			List:         nil,
		}
		if contentItem.Type == model.ContentTypeProducts {
			//menussddddd.ID = 0
			for ii := 0; ii < len(goodsTypeList); ii++ {
				goodsType := goodsTypeList[ii]
				subMenus := extends.Menus{
					ID:           goodsType.ID,
					Uri:          goodsType.Uri,
					Name:         goodsType.Name,
					TemplateName: contentItem.TemplateName,
					Type:         contentItem.Type,
					Introduction: contentItem.Introduction,
					Image:        contentItem.Image,
					List:         nil,
				}
				for iii := 0; iii < len(goodsTypeChildList); iii++ {
					goodsTypeChild := goodsTypeChildList[iii]
					if goodsType.ID == goodsTypeChild.GoodsTypeID {
						subMenus.List = append(subMenus.List, extends.Menus{
							ID:           goodsTypeChild.ID,
							Uri:          goodsTypeChild.Uri,
							Name:         goodsTypeChild.Name,
							TemplateName: contentItem.TemplateName,
							Type:         contentItem.Type,
							List:         nil,
						})
					}
				}
				menussddddd.List = append(menussddddd.List, subMenus)
			}
		} else {
			for ii := 0; ii < len(contentSubTypeList); ii++ {
				contentSubType := contentSubTypeList[ii]
				if menussddddd.ID == contentSubType.ContentItemID && contentSubType.ParentContentSubTypeID == 0 {
					subMenus := extends.Menus{
						ID:           contentSubType.ID,
						Uri:          contentSubType.Uri,
						Name:         contentSubType.Name,
						TemplateName: contentItem.TemplateName,
						Type:         contentItem.Type,
						List:         nil,
					}
					menussddddd.List = append(menussddddd.List, subMenus)
				}
			}

		}
		list = append(list, menussddddd)

	}

	for i := 0; i < len(list); i++ {
		menussddddd := list[i]
		if menussddddd.Type == model.ContentTypeProducts {
			continue
		}
		for ii := 0; ii < len(menussddddd.List); ii++ {
			subMenus := menussddddd.List[ii]

			for iii := 0; iii < len(contentSubTypeList); iii++ {
				contentSubType := contentSubTypeList[iii]
				if contentSubType.ParentContentSubTypeID != 0 && contentSubType.ParentContentSubTypeID == subMenus.ID {
					subSubMenus := extends.Menus{
						ID:           contentSubType.ID,
						Uri:          contentSubType.Uri,
						Name:         contentSubType.Name,
						TemplateName: menussddddd.TemplateName,
						Type:         menussddddd.Type,
						List:         nil,
					}
					subMenus.List = append(subMenus.List[:], subSubMenus)
				}
			}
			menussddddd.List[ii] = subMenus
		}

	}
	menusData.List = list
	return menusData

}
func (service Service) GetGoodsTypeByUri(context constrain.IContext, OID types.PrimaryKey, GoodsTypeUri, GoodsTypeChildUri string, pageIndex int) module.SiteData[*model.Goods] {
	var moduleContentData module.SiteData[*model.Goods]

	Orm := singleton.Orm()
	var item model.GoodsType
	var itemSub model.GoodsTypeChild

	itemMap := map[string]interface{}{"OID": OID, "Uri": GoodsTypeUri}
	Orm.Model(model.GoodsType{}).Where(itemMap).First(&item)

	itemSubMap := map[string]interface{}{
		"OID":         OID,
		"GoodsTypeID": item.ID,
		"Uri":         GoodsTypeChildUri,
	}
	Orm.Model(model.GoodsTypeChild{}).Where(itemSubMap).First(&itemSub)
	if itemSub.IsZero() {
		itemSub.Uri = "all"
	}

	allMenusData := service.FindAllMenus(OID)

	menusData := service.FindShowMenus(OID)

	currentMenuData := module.NewProductMenusData(item, itemSub)
	for _, v := range menusData.List {
		if v.Type == model.ContentTypeProducts {
			currentMenuData.Menus = v
			break
		}
	}

	menusPage := allMenusData.ListMenusByType(model.ContentTypePage)

	pageIndex, pageSize, total, list := service.GoodsService.PaginationGoods(OID, currentMenuData.TypeID, currentMenuData.SubTypeID, pageIndex)

	pagination := module.NewPagination[*model.Goods](pageIndex, pageSize, total, list)

	var navigations []extends.Menus

	for index, v := range menusData.List {
		if v.Type == model.ContentTypeProducts {
			navigations = append(navigations, menusData.List[index])
			for si, sv := range v.List {
				if sv.ID == currentMenuData.TypeID {
					navigations = append(navigations, menusData.List[index].List[si])
					for ssi, ssv := range sv.List {
						if ssv.ID == currentMenuData.SubTypeID {
							navigations = append(navigations, menusData.List[index].List[si].List[ssi])
							break
						}
					}
					break
				}
			}
			break
		}
	}

	organization := service.OrganizationService.GetOrganization(OID).(*model.Organization)
	contentConfig := service.ContentService.GetContentConfig(singleton.Orm(), OID)

	moduleContentData = module.SiteData[*model.Goods]{
		AllMenusData:    allMenusData,
		MenusData:       menusData,
		PageMenus:       menusPage,
		CurrentMenuData: currentMenuData,
		ContentItem:     model.ContentItem{},
		ContentSubType:  model.ContentSubType{},
		Pagination:      pagination,
		Tags:            []extends.Tag{},
		Navigations:     navigations,
		Organization:    *organization,
		ContentConfig:   contentConfig,
		SiteAuthor:      "",
		LeftRight:       [2]*model.Goods{},
	}

	companyName := contentConfig.Name
	if len(companyName) == 0 {
		companyName = organization.Name
	}
	moduleContentData.SiteAuthor = companyName
	return moduleContentData
}
func (service Service) GetContentTypeByUri(context constrain.IContext, OID types.PrimaryKey, ContentItemUri, ContentSubTypeUri string, pageIndex int) module.SiteData[*model.Content] {
	var moduleContentData module.SiteData[*model.Content]
	Orm := singleton.Orm()
	var item model.ContentItem
	var itemSub model.ContentSubType

	itemMap := map[string]interface{}{"OID": OID, "Uri": ContentItemUri}
	Orm.Model(model.ContentItem{}).Where(itemMap).First(&item)

	itemSubMap := map[string]interface{}{
		"OID":           OID,
		"ContentItemID": item.ID,
		"Uri":           ContentSubTypeUri,
	}
	Orm.Model(model.ContentSubType{}).Where(itemSubMap).First(&itemSub)
	if itemSub.IsZero() {
		itemSub.Uri = "all"
	}

	currentMenuData := module.NewMenusData(item, itemSub)

	menusData := service.FindShowMenus(OID)
	for _, v := range menusData.List {
		if v.ID == currentMenuData.TypeID {
			currentMenuData.Menus = v
			break
		}
	}

	allMenusData := service.FindAllMenus(OID)

	pageIndex, pageSize, total, list := service.ContentService.PaginationContent(OID, currentMenuData.TypeID, currentMenuData.SubTypeID, pageIndex)

	pagination := module.NewPagination(pageIndex, pageSize, total, list)

	tags := service.ContentService.FindContentTagsByContentItemID(OID, currentMenuData.TypeID)

	var navigations []extends.Menus

	var typeNameMap = make(map[types.PrimaryKey]extends.Menus)

	for index, v := range menusData.List {
		for sv := range v.List {
			typeNameMap[v.List[sv].ID] = v.List[sv]
			for ssv := range v.List[sv].List {
				typeNameMap[v.List[sv].List[ssv].ID] = v.List[sv].List[ssv]
			}
		}
		if v.ID == currentMenuData.TypeID {
			navigations = append(navigations, menusData.List[index])
			for si, sv := range v.List {
				typeNameMap[sv.ID] = sv
				if sv.ID == currentMenuData.SubTypeID {
					navigations = append(navigations, menusData.List[index].List[si])
				} else {
					for _, ssv := range sv.List {
						typeNameMap[ssv.ID] = ssv
					}
					for ssi, ssv := range sv.List {
						if ssv.ID == currentMenuData.SubTypeID {
							navigations = append(navigations, menusData.List[index].List[si])
							navigations = append(navigations, menusData.List[index].List[si].List[ssi])
							break
						}
					}
				}
			}
			break
		}
	}

	organization := service.OrganizationService.GetOrganization(OID).(*model.Organization)
	contentConfig := service.ContentService.GetContentConfig(singleton.Orm(), OID)

	menusPage := allMenusData.ListMenusByType(model.ContentTypePage)
	moduleContentData = module.SiteData[*model.Content]{
		AllMenusData:    allMenusData,
		MenusData:       menusData,
		PageMenus:       menusPage,
		CurrentMenuData: currentMenuData,
		ContentItem:     item,
		ContentSubType:  itemSub,
		Pagination:      pagination,
		Tags:            tags,
		Navigations:     navigations,
		Organization:    *organization,
		ContentConfig:   contentConfig,
		TypeNameMap:     typeNameMap,
	}

	companyName := contentConfig.Name
	if len(companyName) == 0 {
		companyName = organization.Name
	}
	moduleContentData.SiteAuthor = companyName

	if len(list) > 0 {
		c := list[0]
		moduleContentData.LeftRight = service.ContentService.FindContentListForLeftRight(c.ContentItemID, c.ContentSubTypeID, c.ID, c.CreatedAt)
	}
	return moduleContentData
}

func GetSiteData[T module.ListType](context constrain.IContext, OID types.PrimaryKey) module.SiteData[T] {
	var service = Service{}

	var moduleContentData module.SiteData[T]

	var item model.ContentItem
	var subItem = model.ContentSubType{Uri: "all"}

	currentMenuData := module.NewMenusData(item, subItem)

	menusData := service.FindShowMenus(OID)
	for _, v := range menusData.List {
		if v.ID == currentMenuData.TypeID {
			currentMenuData.Menus = v
			break
		}
	}

	allMenusData := service.FindAllMenus(OID)

	tags := service.ContentService.FindContentTagsByContentItemID(OID, currentMenuData.TypeID)

	var navigations []extends.Menus

	var typeNameMap = make(map[types.PrimaryKey]extends.Menus)

	for index, v := range menusData.List {
		if v.ID == currentMenuData.TypeID {
			navigations = append(navigations, menusData.List[index])
			for si, sv := range v.List {
				typeNameMap[sv.ID] = sv
				if sv.ID == currentMenuData.SubTypeID {
					navigations = append(navigations, menusData.List[index].List[si])
				} else {
					for _, ssv := range sv.List {
						typeNameMap[ssv.ID] = ssv
					}
					for ssi, ssv := range sv.List {
						if ssv.ID == currentMenuData.SubTypeID {
							navigations = append(navigations, menusData.List[index].List[si])
							navigations = append(navigations, menusData.List[index].List[si].List[ssi])
							break
						}
					}
				}
			}
			break
		}
	}

	organization := service.OrganizationService.GetOrganization(OID).(*model.Organization)
	contentConfig := service.ContentService.GetContentConfig(singleton.Orm(), OID)

	menusPage := allMenusData.ListMenusByType(model.ContentTypePage)
	moduleContentData = module.SiteData[T]{
		AllMenusData:    allMenusData,
		MenusData:       menusData,
		PageMenus:       menusPage,
		CurrentMenuData: currentMenuData,
		ContentItem:     item,
		ContentSubType:  subItem,
		Pagination:      module.Pagination[T]{},
		Tags:            tags,
		Navigations:     navigations,
		Organization:    *organization,
		ContentConfig:   contentConfig,
		TypeNameMap:     typeNameMap,
	}

	companyName := contentConfig.Name
	if len(companyName) == 0 {
		companyName = organization.Name
	}
	moduleContentData.SiteAuthor = companyName

	return moduleContentData
}
