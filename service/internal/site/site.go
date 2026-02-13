package site

import (
	"fmt"
	"path"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"github.com/nbvghost/dandelion/repository"
	"github.com/nbvghost/dandelion/service/internal/cache"
	"github.com/nbvghost/dandelion/service/internal/company"
	"github.com/nbvghost/dandelion/service/internal/content"
	"github.com/nbvghost/dandelion/service/internal/goods"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"github.com/samber/lo"
)

type Service struct {
	GoodsService        goods.GoodsService
	OrganizationService company.OrganizationService
	ContentService      content.ContentService
	GoodsTypeService    goods.GoodsTypeService
}

type MenuShowType int

const (
	MenuShowTypeAll  MenuShowType = 0
	MenuShowTypeHide MenuShowType = 1
	MenuShowTypeShow MenuShowType = 2
)

func (m Service) FindShowMenus(ctx constrain.IContext, OID dao.PrimaryKey) extends.MenusData {
	return m.menus(ctx, OID, MenuShowTypeShow)
}
func (m Service) FindAllMenus(ctx constrain.IContext, OID dao.PrimaryKey) extends.MenusData {
	return m.menus(ctx, OID, MenuShowTypeAll)
}
func newRedisContentDataKey(OID dao.PrimaryKey, ContentItemUri, ContentSubTypeUri string, pageIndex int) string {
	return fmt.Sprintf("content:%d:%s:%s:%d", OID, ContentItemUri, ContentSubTypeUri, pageIndex)
}
func newRedisGoodsDataKey(OID dao.PrimaryKey, ContentItemUri, ContentSubTypeUri string, pageIndex int) string {
	return fmt.Sprintf("goods:%d:%s:%s:%d", OID, ContentItemUri, ContentSubTypeUri, pageIndex)
}
func (m Service) menus(ctx constrain.IContext, OID dao.PrimaryKey, showType MenuShowType) extends.MenusData {
	//Orm := db.GetDB(ctx)

	var contentItemList []model.ContentItem

	{
		var list = cache.GetCacheContentItem(ctx, OID)

		switch showType {
		case 0: //all
			/*Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
				"OID": OID,
			}).Order(`"Sort"`).Find(&contentItemList)*/
			contentItemList = list
		case 1: //hide
			/*Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
				"ShowAtMenu": false,
				"OID":        OID,
			}).Order(`"Sort"`).Find(&contentItemList)*/
			contentItemList = lo.Filter[model.ContentItem](list, func(item model.ContentItem, index int) bool {
				return item.ShowAtMenu == false
			})
		case 2: //show
			/*Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
				"ShowAtMenu": true,
				"OID":        OID,
			}).Order(`"Sort"`).Find(&contentItemList)*/
			contentItemList = lo.Filter[model.ContentItem](list, func(item model.ContentItem, index int) bool {
				return item.ShowAtMenu
			})
		default:
			/*Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
				"OID": OID,
			}).Order(`"Sort"`).Find(&contentItemList)*/
			contentItemList = list

		}
	}

	/*var contentItemIDs []dao.PrimaryKey
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
	}*/

	var contentSubTypeList []model.ContentSubType
	{
		var list = cache.GetCacheContentSubType(ctx, OID)
		contentSubTypeList = list
		//Orm.Model(model.ContentSubType{}).Where(`"OID"=?`, OID).Where(`"ContentItemID" in ?`, contentItemIDs).Order(`"Sort"`).Order(`"ID"`).Find(&contentSubTypeList)
	}

	var menusData extends.MenusData

	var rootMenusList []extends.Menus
	for i := 0; i < len(contentItemList); i++ {
		contentItem := contentItemList[i]
		urlPath := path.Join(fmt.Sprintf("/%s", contentItem.Type), contentItem.Uri)
		rootMenus := extends.Menus{
			ParentMenus:  &extends.Menus{},
			ID:           contentItem.ID,
			Name:         contentItem.Name,
			TemplateName: contentItem.TemplateName,
			Type:         contentItem.Type,
			Introduction: contentItem.Introduction,
			Image:        contentItem.Image,
			Badge:        contentItem.Badge,
			ShowAtMenu:   contentItem.ShowAtMenu,
			UrlPath:      urlPath,
			List:         nil,
		}

		if contentItem.Type == model.ContentTypeIndex {
			rootMenus.UrlPath = "/"
		} else if contentItem.Type == model.ContentTypeProducts {
			var goodsTypeList []model.GoodsType
			{
				var list = cache.GetCacheGoodsType(ctx, OID)
				//Orm.Model(model.GoodsType{}).Where(`"OID"=?`, OID).Order(`"ID"`).Find(&goodsTypeList)
				goodsTypeList = list
			}
			var goodsTypeChildList []model.GoodsTypeChild
			{
				var list = cache.GetCacheGoodsTypeChild(ctx, OID)
				goodsTypeChildList = list
				//Orm.Model(model.GoodsTypeChild{}).Where(`"OID" = ?`, OID).Order(`"ID"`).Find(&goodsTypeChildList)
			}
			rootMenus.List = lo.FilterMap[model.GoodsType, extends.Menus](goodsTypeList, func(goodsType model.GoodsType, index int) (extends.Menus, bool) {
				subMenus := extends.Menus{
					ParentMenus:  &rootMenus,
					ID:           goodsType.ID,
					Name:         goodsType.Name,
					TemplateName: contentItem.TemplateName,
					Type:         contentItem.Type,
					Introduction: contentItem.Introduction,
					Image:        contentItem.Image,
					Badge:        contentItem.Badge,
					ShowAtMenu:   goodsType.ShowAtMenu,
					UrlPath:      path.Join(urlPath, goodsType.Uri),
					List:         nil,
				}

				subMenus.List = lo.FilterMap[model.GoodsTypeChild, extends.Menus](goodsTypeChildList, func(goodsTypeChild model.GoodsTypeChild, index int) (extends.Menus, bool) {
					if goodsType.ID == goodsTypeChild.GoodsTypeID {
						return extends.Menus{
							ParentMenus:  &subMenus,
							ID:           goodsTypeChild.ID,
							Name:         goodsTypeChild.Name,
							Image:        goodsTypeChild.Image,
							TemplateName: contentItem.TemplateName,
							Type:         contentItem.Type,
							ShowAtMenu:   false,
							UrlPath:      path.Join(urlPath, goodsTypeChild.Uri),
							List:         nil,
						}, true
					} else {
						return extends.Menus{}, false
					}
				})
				return subMenus, true
			})
		} else {
			rootMenus.List = lo.FilterMap[model.ContentSubType, extends.Menus](contentSubTypeList, func(contentSubType model.ContentSubType, index int) (extends.Menus, bool) {
				if rootMenus.ID == contentSubType.ContentItemID && contentSubType.ParentContentSubTypeID == 0 {
					menus := extends.Menus{
						ParentMenus:  &rootMenus,
						ID:           contentSubType.ID,
						Name:         contentSubType.Name,
						TemplateName: contentItem.TemplateName,
						Type:         contentItem.Type,
						UrlPath:      path.Join(urlPath, contentSubType.Uri),
						ShowAtMenu:   false,
						List:         nil,
					}

					menus.List = lo.FilterMap[model.ContentSubType, extends.Menus](contentSubTypeList, func(contentSubType model.ContentSubType, index int) (extends.Menus, bool) {
						if contentSubType.ParentContentSubTypeID != 0 && contentSubType.ParentContentSubTypeID == menus.ID {
							subMenus := extends.Menus{
								ParentMenus:  &menus,
								ID:           contentSubType.ID,
								Name:         contentSubType.Name,
								TemplateName: contentItem.TemplateName,
								Type:         contentItem.Type,
								UrlPath:      path.Join(urlPath, contentSubType.Uri),
								ShowAtMenu:   false,
								List:         nil,
							}
							return subMenus, true
						} else {
							return extends.Menus{}, false
						}
					})
					return menus, true
				} else {
					return extends.Menus{}, false
				}
			})
		}
		rootMenusList = append(rootMenusList, rootMenus)

	}

	showAtMenu := make([]extends.Menus, 0)
	lo.ForEach(rootMenusList, func(item extends.Menus, index int) {
		showAtMenu = append(showAtMenu, lo.FilterMap[extends.Menus, extends.Menus](item.List, func(item extends.Menus, index int) (extends.Menus, bool) {
			return item, item.ShowAtMenu
		})...)
	})
	menusData.List = append(rootMenusList, showAtMenu...)
	return menusData

}
func (m Service) GoodsList(ctx constrain.IContext, OID dao.PrimaryKey, ContentItemUri, GoodsTypeUri string, filterOption []serviceargument.Option, sortMethod *serviceargument.SortMethod, pageIndex, pageSize int) serviceargument.SiteData[*extends.GoodsDetail] {
	var moduleContentData serviceargument.SiteData[*extends.GoodsDetail]

	Orm := db.GetDB(ctx)

	contentItem := dao.GetBy(Orm, &model.ContentItem{}, map[string]interface{}{"OID": OID, "Uri": ContentItemUri}).(*model.ContentItem)

	var item model.GoodsType
	var itemSub model.GoodsTypeChild

	{
		list := cache.GetCacheGoodsType(ctx, OID)

		item, _ = lo.Find[model.GoodsType](list, func(a model.GoodsType) bool {
			return a.Uri == GoodsTypeUri
		})
		if item.IsZero() {
			clist := cache.GetCacheGoodsTypeChild(ctx, OID)
			itemSub, _ = lo.Find[model.GoodsTypeChild](clist, func(a model.GoodsTypeChild) bool {
				return a.Uri == GoodsTypeUri
			})
			if !itemSub.IsZero() {
				item, _ = lo.Find[model.GoodsType](list, func(a model.GoodsType) bool {
					return a.ID == itemSub.GoodsTypeID
				})
			}
		}

		/*Orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID, "Uri": GoodsTypeUri}).First(&item)
		if item.IsZero() {
			Orm.Model(model.GoodsTypeChild{}).Where(map[string]interface{}{"OID": OID, "Uri": GoodsTypeUri}).First(&itemSub)
			if !itemSub.IsZero() {
				Orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID, "ID": itemSub.GoodsTypeID}).First(&item)
			}
		}*/
	}

	allMenusData := m.FindAllMenus(ctx, OID)
	menusData := m.FindShowMenus(ctx, OID)

	var currentMenus extends.Menus
	if itemSub.IsZero() == false {
		currentMenus = menusData.GetCurrentMenus(itemSub.ID)
	} else if item.IsZero() == false {
		currentMenus = menusData.GetCurrentMenus(item.ID)
	} else {
		currentMenus = menusData.GetCurrentMenus(contentItem.ID)
	}

	contentItemMap := repository.ContentItemDao.ListContentItemByOIDMap(ctx, OID)

	currentMenuData := serviceargument.NewProductMenusData(item, itemSub)
	for _, v := range menusData.List {
		if v.Type == model.ContentTypeProducts {
			currentMenuData.Menus = v
			break
		}
	}

	menusPage := allMenusData.ListMenusByType(model.ContentTypePage)

	pagination, options := m.GoodsService.PaginationGoodsDetail(ctx, OID, currentMenuData.TypeID, currentMenuData.SubTypeID, filterOption, sortMethod, pageIndex, pageSize)

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

	organization := m.OrganizationService.GetOrganization(Orm, OID).(*model.Organization)
	contentConfig := repository.ContentConfigDao.GetContentConfig(db.GetDB(ctx), OID)

	moduleContentData = serviceargument.SiteData[*extends.GoodsDetail]{
		AllMenusData:    allMenusData,
		MenusData:       menusData,
		PageMenus:       menusPage,
		CurrentMenuData: currentMenuData,
		CurrentMenu:     currentMenus,
		ContentItem:     contentItem,
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
func (m Service) GoodsDetail(ctx constrain.IContext, OID dao.PrimaryKey, GoodsTypeUri, GoodsTypeChildUri string) serviceargument.SiteData[*extends.GoodsDetail] {
	var moduleContentData serviceargument.SiteData[*extends.GoodsDetail]

	Orm := db.GetDB(ctx)
	var item model.GoodsType
	var itemSub model.GoodsTypeChild

	Orm.Model(model.GoodsType{}).Where(map[string]interface{}{"OID": OID, "Uri": GoodsTypeUri}).First(&item)

	Orm.Model(model.GoodsTypeChild{}).Where(map[string]interface{}{"OID": OID, "GoodsTypeID": item.ID, "Uri": GoodsTypeChildUri}).First(&itemSub)
	if itemSub.IsZero() {
		//itemSub.Uri = "all"
	}

	contentItemMap := repository.ContentItemDao.ListContentItemByOIDMap(ctx, OID)

	allMenusData := m.FindAllMenus(ctx, OID)

	menusData := m.FindShowMenus(ctx, OID)

	currentMenuData := serviceargument.NewProductMenusData(item, itemSub)
	for _, v := range menusData.List {
		if v.Type == model.ContentTypeProducts {
			currentMenuData.Menus = v
			break
		}
	}

	menusPage := allMenusData.ListMenusByType(model.ContentTypePage)

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

	organization := m.OrganizationService.GetOrganization(Orm, OID).(*model.Organization)
	contentConfig := repository.ContentConfigDao.GetContentConfig(db.GetDB(ctx), OID)

	leftRightArr := [2]*extends.GoodsDetail{}

	moduleContentData = serviceargument.SiteData[*extends.GoodsDetail]{
		AllMenusData:    allMenusData,
		MenusData:       menusData,
		PageMenus:       menusPage,
		CurrentMenuData: currentMenuData,
		ContentItem:     &model.ContentItem{},
		ContentSubType:  model.ContentSubType{},
		Pagination:      serviceargument.Pagination[*extends.GoodsDetail]{}, //pagination,
		Tags:            []extends.Tag{},
		Navigations:     navigations,
		Organization:    *organization,
		ContentConfig:   contentConfig,
		SiteAuthor:      "",
		LeftRight:       leftRightArr,
		ContentItemMap:  contentItemMap,
		Options:         serviceargument.Options{Attributes: make([]serviceargument.Option, 0)},
	}

	companyName := contentConfig.Name
	if len(companyName) == 0 {
		companyName = organization.Name
	}
	moduleContentData.SiteAuthor = companyName
	return moduleContentData
}
func (m Service) GetContentTypeByUri(ctx constrain.IContext, OID dao.PrimaryKey, ContentItemUri, ContentSubTypeUri string, pageIndex int) serviceargument.SiteData[*model.Content] {
	var moduleContentData serviceargument.SiteData[*model.Content]
	Orm := db.GetDB(ctx)

	var itemSub model.ContentSubType

	//itemMap := map[string]interface{}{"OID": OID, "Uri": ContentItemUri}
	//Orm.Model(model.ContentItem{}).Where(itemMap).First(&item)

	contentItem := dao.GetBy(Orm, &model.ContentItem{}, map[string]interface{}{"OID": OID, "Uri": ContentItemUri}).(*model.ContentItem)

	itemSubMap := map[string]interface{}{
		"OID":           OID,
		"ContentItemID": contentItem.ID,
		"Uri":           ContentSubTypeUri,
	}
	Orm.Model(model.ContentSubType{}).Where(itemSubMap).First(&itemSub)
	if itemSub.IsZero() {
		itemSub.Uri = "all"
	}

	contentItemMap := repository.ContentItemDao.ListContentItemByOIDMap(ctx, OID)

	currentMenuData := serviceargument.NewMenusData(contentItem, itemSub)

	menusData := m.FindShowMenus(ctx, OID)

	allMenusData := m.FindAllMenus(ctx, OID)
	for _, v := range allMenusData.List {
		if v.ID == currentMenuData.TypeID {
			currentMenuData.Menus = v
			break
		}
	}

	pageIndex, pageSize, total, list := m.ContentService.PaginationContent(ctx, OID, currentMenuData.TypeID, currentMenuData.SubTypeID, pageIndex, 20)

	pagination := serviceargument.NewPagination(pageIndex, pageSize, total, list)

	tags := m.ContentService.FindContentTagsByContentItemID(ctx, OID, currentMenuData.TypeID)

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

	organization := m.OrganizationService.GetOrganization(Orm, OID).(*model.Organization)
	contentConfig := repository.ContentConfigDao.GetContentConfig(db.GetDB(ctx), OID)

	menusPage := allMenusData.ListMenusByType(model.ContentTypePage)
	moduleContentData = serviceargument.SiteData[*model.Content]{
		AllMenusData:    allMenusData,
		MenusData:       menusData,
		PageMenus:       menusPage,
		CurrentMenuData: currentMenuData,
		ContentItem:     contentItem,
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
		moduleContentData.LeftRight = m.ContentService.FindContentListForLeftRight(ctx, c.ContentItemID, c.ContentSubTypeID, c.ID, c.CreatedAt)
	}
	return moduleContentData
}
