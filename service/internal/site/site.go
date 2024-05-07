package site

import (
	"fmt"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service/internal/company"
	"github.com/nbvghost/dandelion/service/internal/content"
	"github.com/nbvghost/dandelion/service/internal/goods"
	"github.com/nbvghost/dandelion/service/serviceargument"
)

type Service struct {
	GoodsService        goods.GoodsService
	OrganizationService company.OrganizationService
	ContentService      content.ContentService
}

func (m Service) FindShowMenus(OID dao.PrimaryKey) extends.MenusData {
	return m.menus(OID, 2)
}
func (m Service) FindAllMenus(OID dao.PrimaryKey) extends.MenusData {
	return m.menus(OID, 0)
}
func newRedisContentDataKey(OID dao.PrimaryKey, ContentItemUri, ContentSubTypeUri string, pageIndex int) string {
	return fmt.Sprintf("content:%d:%s:%s:%d", OID, ContentItemUri, ContentSubTypeUri, pageIndex)
}
func newRedisGoodsDataKey(OID dao.PrimaryKey, ContentItemUri, ContentSubTypeUri string, pageIndex int) string {
	return fmt.Sprintf("goods:%d:%s:%s:%d", OID, ContentItemUri, ContentSubTypeUri, pageIndex)
}
func (m Service) menus(OID dao.PrimaryKey, hide uint) extends.MenusData {
	Orm := db.Orm()

	var contentItemList []model.ContentItem

	switch hide {
	case 0: //all
		Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
			"OID": OID,
		}).Order(`"Sort"`).Find(&contentItemList)
	case 1: //hide
		Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
			"ShowAtMenu": false,
			"OID":        OID,
		}).Order(`"Sort"`).Find(&contentItemList)
	case 2: //show
		Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
			"ShowAtMenu": true,
			"OID":        OID,
		}).Order(`"Sort"`).Find(&contentItemList)
	default:
		Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
			"OID": OID,
		}).Order(`"Sort"`).Find(&contentItemList)

	}

	var contentItemIDs []dao.PrimaryKey
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
			Badge:        contentItem.Badge,
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
					Badge:        contentItem.Badge,
					List:         nil,
				}
				for iii := 0; iii < len(goodsTypeChildList); iii++ {
					goodsTypeChild := goodsTypeChildList[iii]
					if goodsType.ID == goodsTypeChild.GoodsTypeID {
						subMenus.List = append(subMenus.List, extends.Menus{
							ID:           goodsTypeChild.ID,
							Uri:          goodsTypeChild.Uri,
							Name:         goodsTypeChild.Name,
							Image:        goodsTypeChild.Image,
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

func (m Service) GoodsDetail(context constrain.IContext, OID dao.PrimaryKey, GoodsTypeUri, GoodsTypeChildUri string, filterOption []serviceargument.Option, pageIndex int) serviceargument.SiteData[*extends.GoodsDetail] {
	var moduleContentData serviceargument.SiteData[*extends.GoodsDetail]

	Orm := db.Orm()
	var item model.GoodsType
	var itemSub model.GoodsTypeChild

	Orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID, "Uri": GoodsTypeUri}).First(&item)

	Orm.Model(model.GoodsTypeChild{}).Where(map[string]interface{}{"OID": OID, "GoodsTypeID": item.ID, "Uri": GoodsTypeChildUri}).First(&itemSub)
	if itemSub.IsZero() {
		//itemSub.Uri = "all"
	}

	contentItemMap := repository.ContentItemDao.ListContentItemByOIDMap(OID)

	allMenusData := m.FindAllMenus(OID)

	menusData := m.FindShowMenus(OID)

	currentMenuData := serviceargument.NewProductMenusData(item, itemSub)
	for _, v := range menusData.List {
		if v.Type == model.ContentTypeProducts {
			currentMenuData.Menus = v
			break
		}
	}

	menusPage := allMenusData.ListMenusByType(model.ContentTypePage)

	pagination, options := m.GoodsService.PaginationGoodsDetail(OID, currentMenuData.TypeID, currentMenuData.SubTypeID, filterOption, pageIndex)

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

	organization := m.OrganizationService.GetOrganization(OID).(*model.Organization)
	contentConfig := repository.ContentConfigDao.GetContentConfig(db.Orm(), OID)

	moduleContentData = serviceargument.SiteData[*extends.GoodsDetail]{
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
		LeftRight:       [2]*extends.GoodsDetail{},
		ContentItemMap:  contentItemMap,
		Options:         *options,
	}

	companyName := contentConfig.Name
	if len(companyName) == 0 {
		companyName = organization.Name
	}
	moduleContentData.SiteAuthor = companyName
	return moduleContentData
}
func (m Service) GetContentTypeByUri(context constrain.IContext, OID dao.PrimaryKey, ContentItemUri, ContentSubTypeUri string, pageIndex int) serviceargument.SiteData[*model.Content] {
	var moduleContentData serviceargument.SiteData[*model.Content]
	Orm := db.Orm()
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

	contentItemMap := repository.ContentItemDao.ListContentItemByOIDMap(OID)

	currentMenuData := serviceargument.NewMenusData(item, itemSub)

	menusData := m.FindShowMenus(OID)

	allMenusData := m.FindAllMenus(OID)
	for _, v := range allMenusData.List {
		if v.ID == currentMenuData.TypeID {
			currentMenuData.Menus = v
			break
		}
	}

	pageIndex, pageSize, total, list := m.ContentService.PaginationContent(OID, currentMenuData.TypeID, currentMenuData.SubTypeID, pageIndex, 20)

	pagination := serviceargument.NewPagination(pageIndex, pageSize, total, list)

	tags := m.ContentService.FindContentTagsByContentItemID(OID, currentMenuData.TypeID)

	var navigations []extends.Menus

	var typeNameMap = make(map[dao.PrimaryKey]extends.Menus)

	for index, v := range allMenusData.List {
		for sv := range v.List {
			typeNameMap[v.List[sv].ID] = v.List[sv]
			for ssv := range v.List[sv].List {
				typeNameMap[v.List[sv].List[ssv].ID] = v.List[sv].List[ssv]
			}
		}
		if v.ID == currentMenuData.TypeID {
			navigations = append(navigations, allMenusData.List[index])
			for si, sv := range v.List {
				typeNameMap[sv.ID] = sv
				if sv.ID == currentMenuData.SubTypeID {
					navigations = append(navigations, allMenusData.List[index].List[si])
				} else {
					for _, ssv := range sv.List {
						typeNameMap[ssv.ID] = ssv
					}
					for ssi, ssv := range sv.List {
						if ssv.ID == currentMenuData.SubTypeID {
							navigations = append(navigations, allMenusData.List[index].List[si])
							navigations = append(navigations, allMenusData.List[index].List[si].List[ssi])
							break
						}
					}
				}
			}
			break
		}
	}

	organization := m.OrganizationService.GetOrganization(OID).(*model.Organization)
	contentConfig := repository.ContentConfigDao.GetContentConfig(db.Orm(), OID)

	menusPage := allMenusData.ListMenusByType(model.ContentTypePage)
	moduleContentData = serviceargument.SiteData[*model.Content]{
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
		ContentItemMap:  contentItemMap,
	}

	companyName := contentConfig.Name
	if len(companyName) == 0 {
		companyName = organization.Name
	}
	moduleContentData.SiteAuthor = companyName

	if len(list) > 0 {
		c := list[0]
		moduleContentData.LeftRight = m.ContentService.FindContentListForLeftRight(c.ContentItemID, c.ContentSubTypeID, c.ID, c.CreatedAt)
	}
	return moduleContentData
}
