package content

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/dandelion/service/journal"

	"gorm.io/gorm"

	"github.com/nbvghost/gpa/types"
	"github.com/nbvghost/tool/object"
	"strings"
	"time"

	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
	"strconv"
)

type ContentService struct {
	model.BaseDao
	Journal journal.JournalService
}

func (service ContentService) AddContentItem(company *model.Organization, item *model.ContentItem) *result.ActionResult {
	Orm := singleton.Orm()

	have := service.GetContentItemByNameAndOID(item.Name, company.ID)
	if have.ID != 0 || strings.EqualFold(item.Name, "") {
		return &result.ActionResult{
			Code:    result.Fail,
			Message: "这个名字已经被使用了",
			Data:    nil,
		} //&gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("这个名字已经被使用了"), "", nil)}
	}

	var mt model.ContentType
	Orm.Where(map[string]interface{}{"ID": item.ContentTypeID}).First(&mt)
	if mt.ID == 0 {
		return &result.ActionResult{
			Code:    result.Fail,
			Message: "没有找到类型",
			Data:    nil,
		} //&gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("没有找到类型"), "", nil)}
	}

	if strings.EqualFold(string(mt.Type), string(model.ContentTypeBlog)) || strings.EqualFold(string(mt.Type), string(model.ContentTypeProducts)) {
		haveList := service.FindContentItemByType(mt.Type, company.ID)
		if len(haveList) != 0 {
			return &result.ActionResult{
				Code:    result.Fail,
				Message: fmt.Sprintf("这个类型（%v）只能创建一个", have.Name),
				Data:    nil,
			} // &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New(fmt.Sprintf("这个类型（%v）只能创建一个", content_item.Type)), "", nil)}
		}
	}

	item.OID = company.ID
	item.Type = mt.Type

	{

		contentItemList := service.ListContentItemByOID(company.ID)
		if len(contentItemList) > 0 {
			item.Sort = contentItemList[len(contentItemList)-1].Sort + 1
		}
	}

	err := service.Add(Orm, item)
	if glog.Error(err) {
		return &result.ActionResult{
			Code:    result.SQLError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return &result.ActionResult{
		Code:    result.Success,
		Message: "添加成功",
		Data:    nil,
	} //(&result.ActionResult{}).SmartError(err, "添加成功", nil)
}
func (service ContentService) ChangeContentConfig(OID types.PrimaryKey, fieldName, fieldValue string) error {

	changeMap := make(map[string]interface{})

	switch fieldName {
	case "Name":
		changeMap["Name"] = fieldValue
	case "Logo":
		changeMap["Logo"] = fieldValue
	case "SocialAccount":
		var socialAccount sqltype.SocialAccountList
		json.Unmarshal([]byte(fieldValue), &socialAccount)
		changeMap["SocialAccount"] = socialAccount
	case "CustomerService":
		var customerService sqltype.CustomerServiceList
		json.Unmarshal([]byte(fieldValue), &customerService)
		changeMap["CustomerService"] = customerService
	case "EnableHTMLCache":
		EnableHTMLCache, _ := strconv.ParseBool(fieldValue)
		changeMap["EnableHTMLCache"] = EnableHTMLCache
	case "FocusPicture":
		var focusPicture sqltype.FocusPictureList
		json.Unmarshal([]byte(fieldValue), &focusPicture)
		changeMap["FocusPicture"] = focusPicture

	}
	Orm := singleton.Orm()
	err := Orm.Model(&model.ContentConfig{}).Where(map[string]interface{}{"OID": OID}).Updates(changeMap).Error
	return err
}

func (service ContentService) AddContentConfig(db *gorm.DB, company *model.Organization) error {
	Orm := db
	item := service.GetContentConfig(db, company.ID)
	if (&item).IsZero() {
		err := Orm.Create(&model.ContentConfig{OID: company.ID, Name: company.Name}).Error
		return err
	}
	return nil
}

func (service ContentService) GetContentConfig(db *gorm.DB, OID types.PrimaryKey) model.ContentConfig {
	Orm := db
	var contentConfig model.ContentConfig
	Orm.Model(&model.ContentConfig{}).Where(map[string]interface{}{"OID": OID}).First(&contentConfig)
	return contentConfig
}

func (service ContentService) GetContentItemDefault(db *gorm.DB, OID types.PrimaryKey) *model.ContentItem {
	Orm := db
	var contentItem model.ContentItem
	Orm.Model(&model.ContentItem{}).Where("OID=? And Type=?", OID, model.ContentTypeProducts).First(&contentItem)
	if contentItem.ID == 0 {
		return nil
	}
	return &contentItem
}

func (service ContentService) GetContentItemIDs(OID types.PrimaryKey) []uint {
	Orm := singleton.Orm()
	var levea []uint
	if OID <= 0 {
		return levea
	}
	Orm.Model(&model.ContentItem{}).Where(map[string]interface{}{"OID": OID}).Pluck(`"ID"`, &levea)
	return levea
}

func (service ContentService) GetClassifyByName(Name string, ContentItemID, ParentContentSubTypeID types.PrimaryKey) model.ContentSubType {
	Orm := singleton.Orm()
	var menus model.ContentSubType

	Orm.Where(map[string]interface{}{
		"Name":                   Name,
		"ContentItemID":          ContentItemID,
		"ParentContentSubTypeID": ParentContentSubTypeID,
	}).First(&menus)

	return menus

}

func (service ContentService) GetContentSubTypeByNameContentItemIDParentContentSubTypeID(Name string, ContentItemID, ParentContentSubTypeID uint) model.ContentSubType {
	Orm := singleton.Orm()
	var menus model.ContentSubType

	Orm.Where("Name=?", Name).Where("ContentItemID=? and ParentContentSubTypeID=?", ContentItemID, ParentContentSubTypeID).First(&menus)

	return menus

}

func (service ContentService) FindMenusData(OID types.PrimaryKey) extends.MenusData {
	Orm := singleton.Orm()

	var contentItemList []model.ContentItem
	Orm.Model(model.ContentItem{}).Where(map[string]interface{}{
		"Hide": false,
		"OID":  OID,
	}).Order(`"Sort"`).Find(&contentItemList)

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

	var menus extends.MenusData

	list := []extends.Menus{}
	for i := 0; i < len(contentItemList); i++ {
		contentItem := contentItemList[i]
		menus := extends.Menus{
			ID:           contentItem.ID,
			Name:         contentItem.Name,
			TemplateName: contentItem.TemplateName,
			Type:         contentItem.Type,
			List:         nil,
		}
		if contentItem.Type == model.ContentTypeProducts {
			menus.ID = 0
			for ii := 0; ii < len(goodsTypeList); ii++ {
				goodsType := goodsTypeList[ii]
				subMenus := extends.Menus{
					ID:           goodsType.ID,
					Name:         goodsType.Name,
					TemplateName: contentItem.TemplateName,
					Type:         contentItem.Type,
					List:         nil,
				}
				for iii := 0; iii < len(goodsTypeChildList); iii++ {
					goodsTypeChild := goodsTypeChildList[iii]
					if goodsType.ID == goodsTypeChild.GoodsTypeID {
						subMenus.List = append(subMenus.List, extends.Menus{
							ID:           goodsTypeChild.ID,
							Name:         goodsTypeChild.Name,
							TemplateName: contentItem.TemplateName,
							Type:         contentItem.Type,
							List:         nil,
						})
					}
				}
				menus.List = append(menus.List, subMenus)
			}
		} else {
			for ii := 0; ii < len(contentSubTypeList); ii++ {
				contentSubType := contentSubTypeList[ii]
				if menus.ID == contentSubType.ContentItemID && contentSubType.ParentContentSubTypeID == 0 {
					subMenus := extends.Menus{
						ID:           contentSubType.ID,
						Name:         contentSubType.Name,
						TemplateName: contentItem.TemplateName,
						Type:         contentItem.Type,
						List:         nil,
					}
					menus.List = append(menus.List, subMenus)
				}
			}

		}
		list = append(list, menus)

	}

	for i := 0; i < len(list); i++ {
		menus := list[i]
		if menus.Type == model.ContentTypeProducts {
			continue
		}
		for ii := 0; ii < len(menus.List); ii++ {
			subMenus := menus.List[ii]

			for iii := 0; iii < len(contentSubTypeList); iii++ {
				contentSubType := contentSubTypeList[iii]
				if contentSubType.ParentContentSubTypeID != 0 && contentSubType.ParentContentSubTypeID == subMenus.ID {
					subSubMenus := extends.Menus{
						ID:           contentSubType.ID,
						Name:         contentSubType.Name,
						TemplateName: menus.TemplateName,
						Type:         menus.Type,
						List:         nil,
					}
					subMenus.List = append(subMenus.List[:], subSubMenus)
				}
			}
			menus.List[ii] = subMenus
		}

	}
	menus.List = list
	return menus

}
func (service ContentService) FindContentSubTypesByContentItemIDs(ContentItemIDs []uint) []model.ContentSubType {
	Orm := singleton.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where("ContentItemID in (?)", ContentItemIDs).Order("Sort asc").Find(&menus)
	return menus
}
func (service ContentService) FindContentSubTypesByContentItemID(ContentItemID uint) []model.ContentSubType {
	Orm := singleton.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"ContentItemID":          ContentItemID,
		"ParentContentSubTypeID": 0,
	}).Order(`"Sort" asc`).Find(&menus)
	return menus
}
func (service ContentService) FindContentSubTypesByParentContentSubTypeID(ParentContentSubTypeID uint) []model.ContentSubType {
	Orm := singleton.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where("ParentContentSubTypeID=?", ParentContentSubTypeID).Order("Sort asc").Find(&menus)
	return menus
}
func (service ContentService) FindContentSubTypesByContentItemIDAndParentContentSubTypeID(ContentItemID, ParentContentSubTypeID types.PrimaryKey) []model.ContentSubType {
	Orm := singleton.Orm()
	var menus []model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"ContentItemID":          ContentItemID,
		"ParentContentSubTypeID": ParentContentSubTypeID,
	}).Order(`"Sort" asc`).Find(&menus)
	return menus
}

//-----------------------------------------Content----------------------------------------------------------

func (service ContentService) GetContentItemByIDAndOID(ID, OID uint) model.ContentItem {
	Orm := singleton.Orm()
	var menus model.ContentItem

	Orm.Where(`"ID"=? and "OID"=?`, ID, OID).First(&menus)

	return menus
}
func (service ContentService) GetContentItemByID(ID types.PrimaryKey) model.ContentItem {
	Orm := singleton.Orm()
	var menus model.ContentItem
	Orm.Where(`"ID"=?`, ID).First(&menus)
	return menus
}
func (service ContentService) GetContentSubTypeByID(ID types.PrimaryKey) model.ContentSubType {
	Orm := singleton.Orm()
	var menus model.ContentSubType
	Orm.Where(`"ID"=?`, ID).First(&menus)
	return menus
}
func (service ContentService) ListContentItemByOID(OID types.PrimaryKey) []model.ContentItem {
	Orm := singleton.Orm()
	var menus []model.ContentItem
	Orm.Model(model.ContentItem{}).Where(map[string]interface{}{"OID": OID}).Order(`"Sort"`).Order(`"UpdatedAt" desc`).Find(&menus)
	return menus
}
func (service ContentService) GetContentItemByNameAndOID(Name string, OID types.PrimaryKey) model.ContentItem {
	Orm := singleton.Orm()
	var menus model.ContentItem

	Orm.Where(map[string]interface{}{"Name": Name, "OID": OID}).First(&menus)

	return menus
}

func (service ContentService) FindContentItemByType(Type model.ContentTypeType, OID types.PrimaryKey) []model.ContentItem {
	Orm := singleton.Orm()
	menus := make([]model.ContentItem, 0)
	Orm.Where(map[string]interface{}{
		"Type": Type,
		"OID":  OID,
	}).Find(&menus)
	return menus
}

func (service ContentService) ListContentType() []model.ContentType {
	Orm := singleton.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	var list []model.ContentType
	err := service.FindAll(Orm, &list)
	glog.Trace(err)
	return list
}
func (service ContentService) ListContentTypeByType(Type string) model.ContentType {
	Orm := singleton.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	var list model.ContentType
	err := service.FindWhere(Orm, &list, "Type=?", Type)
	glog.Trace(err)
	return list
}
func (service ContentService) FindContentSubTypesByNameAndContentItemID(Name string, ContentItemID types.PrimaryKey) model.ContentSubType {
	Orm := singleton.Orm()
	var cst model.ContentSubType
	Orm.Where("ContentItemID=? and Name=?", ContentItemID, Name).First(&cst)
	return cst
}

//-----------------------
func (service ContentService) AddSpiderContent(OID types.PrimaryKey, ContentName string, ContentSubTypeName string, Author, Title string, FromUrl string, Introduce string, Picture string, Content string, CreatedAt time.Time) {
	var article model.Content
	article.Title = Title
	article.FromUrl = FromUrl
	article.CreatedAt = CreatedAt
	article.UpdatedAt = CreatedAt

	IntroduceRune := []rune(Introduce)
	if len(IntroduceRune) > 255 {
		article.Introduce = string(IntroduceRune[:255])
	} else {
		article.Introduce = Introduce
	}

	//Picture=tool.DownloadInternetImage(Picture,"Mozilla/5.0 (Linux; Android 7.0; SLA-AL00 Build/HUAWEISLA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/6.2 TBS/044109 Mobile Safari/537.36 MicroMessenger/6.6.7.1321(0x26060739) NetType/WIFI Language/zh_CN",weixin_tmp_url)
	article.Picture = Picture
	article.Content = Content

	contentType := service.ListContentTypeByType(play.ContentTypeArticles)

	content := service.GetContentItemByNameAndOID(ContentName, OID)
	if content.ID == 0 {
		content.OID = OID
		content.Type = contentType.Type
		content.Name = ContentName
		content.ContentTypeID = contentType.ID
		service.Save(singleton.Orm(), &content)

	}

	article.ContentItemID = content.ID
	contentSubType := service.FindContentSubTypesByNameAndContentItemID(ContentSubTypeName, content.ID)
	if contentSubType.ID == 0 {
		contentSubType.Name = ContentSubTypeName
		contentSubType.ContentItemID = content.ID
		service.Save(singleton.Orm(), &contentSubType)
	}

	article.Author = Author
	article.ContentSubTypeID = contentSubType.ID
	service.AddContent(&article)

}

func (service ContentService) ChangeContent(article *model.Content) error {

	return service.Save(singleton.Orm(), article)
}

func (service ContentService) GetContentByTitle(Orm *gorm.DB, Title string) *model.Content {
	article := &model.Content{}
	err := Orm.Where("Title=?", Title).First(article).Error //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)
	return article
}
func (service ContentService) DelContent(ID types.PrimaryKey) error {
	err := service.Delete(singleton.Orm(), &model.Content{}, ID)
	return err
}

//获取ID，返回子类ID,包括本身
func (service ContentService) GetContentSubTypeAllIDByID(ContentItemID, ContentSubTypeID types.PrimaryKey) []uint {
	var IDList []uint
	err := singleton.Orm().Model(&model.ContentSubType{}).Where("ContentItemID=? and (ID=? or ParentContentSubTypeID=?)", ContentItemID, ContentSubTypeID, ContentSubTypeID).Pluck("ID", &IDList).Error
	glog.Error(err)
	return IDList
}
func (service ContentService) FindContentByContentSubTypeID(ContentSubTypeID uint) []model.Content {
	var contentList []model.Content
	err := service.FindWhere(singleton.Orm(), &contentList, "ContentSubTypeID=?", ContentSubTypeID) //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)
	return contentList
}
func (service ContentService) FindContentByContentItemIDAndContentSubTypeID(ContentItemID uint, ContentSubTypeID uint) model.Content {

	var content model.Content
	if ContentItemID == 0 {
		glog.Trace("参数ContentItemID为0")
		return content
	}
	if ContentSubTypeID == 0 {
		glog.Trace("参数ContentSubTypeID为0")
		return content
	}

	service.FindWhere(singleton.Orm(), &content, "ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeID) //SelectOne(user, "select * from User where Email=?", Email)

	return content
}

//ContentItemID
//ContentSubTypeID
//ContentSubTypeChildID
func (service ContentService) FindContentByTypeID(menusData *extends.MenusData, ContentItemID, ContentSubTypeID, ContentSubTypeChildID types.PrimaryKey) model.Content {

	var content model.Content

	/*if ContentItemID == 0 {
		glog.Trace("参数ContentItemID为0")
		return content
	}

	if ContentSubTypeID == 0 && ContentSubTypeChildID == 0 {

		singleton.Orm().Model(&model.Content{}).
			Where("ContentItemID=?", ContentItemID).
			Where("ContentSubTypeID=?", 0).
			Order("CreatedAt desc").Order("ID desc").First(&content)

		if content.ID > 0 {
			return content
		}

		if len(menusData.Top.SubType) > 0 {

			singleton.Orm().Model(&model.Content{}).
				Where("ContentItemID=?", ContentItemID).
				Where("ContentSubTypeID=?", menusData.Top.SubType[0].Item.ID).
				Order("CreatedAt desc").Order("ID desc").First(&content)
			if content.ID > 0 {
				ContentSubTypeID = menusData.Top.SubType[0].Item.ID
				menusData.SetCurrentMenus(ContentItemID, ContentSubTypeID, ContentSubTypeChildID)
				return content
			}

			if len(menusData.Top.SubType[0].SubType) > 0 {

				singleton.Orm().Model(&model.Content{}).
					Where("ContentItemID=?", ContentItemID).
					Where("ContentSubTypeID=?", menusData.Top.SubType[0].SubType[0].Item.ID).
					Order("CreatedAt desc").Order("ID desc").First(&content)
				if content.ID > 0 {
					ContentSubTypeID = menusData.Top.SubType[0].Item.ID
					ContentSubTypeChildID = menusData.Top.SubType[0].SubType[0].Item.ID
					menusData.SetCurrentMenus(ContentItemID, ContentSubTypeID, ContentSubTypeChildID)
					return content
				}
			}
		}

		return content
	} else {

		if ContentSubTypeChildID > 0 {
			singleton.Orm().Model(&model.Content{}).
				Where("ContentItemID=?", ContentItemID).
				Where("ContentSubTypeID=?", ContentSubTypeChildID).
				Order("CreatedAt desc").Order("ID desc").First(&content)

			return content

		} else if ContentSubTypeID > 0 {
			singleton.Orm().Model(&model.Content{}).
				Where("ContentItemID=?", ContentItemID).
				Where("ContentSubTypeID=?", ContentSubTypeID).
				Order("CreatedAt desc").Order("ID desc").First(&content)
			if content.ID > 0 {
				//ContentSubTypeID =menusData.Top.SubType[0].Item.ID
				//ContentSubTypeChildID=menusData.Top.SubType[0].SubType[0].Item.ID
				return content
			}

			if len(menusData.Sub.SubType) > 0 {
				singleton.Orm().Model(&model.Content{}).
					Where("ContentItemID=?", ContentItemID).
					Where("ContentSubTypeID=?", menusData.Sub.SubType[0].Item.ID).
					Order("CreatedAt desc").Order("ID desc").First(&content)
				if content.ID > 0 {
					//ContentSubTypeID =menusData.Top.SubType[0].Item.ID
					ContentSubTypeChildID = menusData.Sub.SubType[0].Item.ID
					menusData.SetCurrentMenus(ContentItemID, ContentSubTypeID, ContentSubTypeChildID)
					return content
				}
			}
		}

	}*/

	return content
}
func (service ContentService) FindContentListByTypeID(menusData *extends.MenusData, ContentItemID, ContentSubTypeID, ContentSubTypeChildID types.PrimaryKey, _Page int, _Limit int) result.Pager {

	var pager result.Pager

	if ContentItemID == 0 {
		glog.Trace("参数ContentItemID为0")
		return pager
	}

	if ContentSubTypeID == 0 && ContentSubTypeChildID == 0 {
		if len(menusData.List) > 0 {

		}
		db := singleton.Orm().Model(&model.Content{}).Where("ContentItemID=?", ContentItemID).
			Order("CreatedAt desc").Order("ID desc")
		return model.Paging(db, _Page, _Limit, model.Content{})
	} else {

		if ContentSubTypeChildID > 0 {
			db := singleton.Orm().Model(&model.Content{}).Where("ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeChildID).
				Order("CreatedAt desc").Order("ID desc")
			return model.Paging(db, _Page, _Limit, model.Content{})
		} else {

			ContentSubTypeIDList := service.GetContentSubTypeAllIDByID(ContentItemID, ContentSubTypeID)

			db := singleton.Orm().Model(&model.Content{}).
				Where("ContentItemID=? and ContentSubTypeID in (?)", ContentItemID, ContentSubTypeIDList).
				Order("CreatedAt desc").Order("ID desc")
			return model.Paging(db, _Page, _Limit, model.Content{})
		}
	}

}

func (service ContentService) FindContentListForLeftRight(ContentItemID, ContentSubTypeID, ContentSubTypeChildID types.PrimaryKey, ContentID types.PrimaryKey, ContentCreatedAt time.Time) [2]model.Content {

	var contentList [2]model.Content
	if ContentItemID == 0 {
		glog.Trace("参数ContentItemID为0")
		return contentList
	}

	var ContentSubTypeIDList []uint
	if ContentSubTypeID == 0 && ContentSubTypeChildID == 0 {
		ContentSubTypeIDList = []uint{}
	} else {
		if ContentSubTypeChildID > 0 {
			ContentSubTypeIDList = service.GetContentSubTypeAllIDByID(ContentItemID, ContentSubTypeChildID)
		} else {
			ContentSubTypeIDList = service.GetContentSubTypeAllIDByID(ContentItemID, ContentSubTypeID)
		}

	}

	ContentSubTypeIDListStr := make([]string, 0)
	for index := range ContentSubTypeIDList {
		ContentSubTypeIDListStr = append(ContentSubTypeIDListStr, object.ParseString(ContentSubTypeIDList[index]))
	}

	var whereSql = ""

	if len(ContentSubTypeIDList) > 0 {
		whereSql = fmt.Sprintf("'ContentItemID'=%v and 'ContentSubTypeID' in (%v)", ContentItemID, strings.Join(ContentSubTypeIDListStr, ","))
	} else {
		whereSql = fmt.Sprintf("'ContentItemID'=%v", ContentItemID)
	}

	var left model.Content
	var right model.Content
	err := singleton.Orm().Raw(`SELECT * FROM 'Content'  WHERE `+whereSql+` and ID<>? and CreatedAt>=? ORDER BY CreatedAt,ID limit 1`, ContentID, ContentCreatedAt).Scan(&left).Error
	glog.Error(err)
	err = singleton.Orm().Raw(`SELECT * FROM 'Content'  WHERE `+whereSql+` and ID<>? and CreatedAt<=? ORDER BY CreatedAt desc,ID desc limit 1`, ContentID, ContentCreatedAt).Scan(&right).Error
	glog.Error(err)

	return [2]model.Content{left, right}
}

func (service ContentService) GetContentByContentItemID(ContentItemID uint) *model.Content {
	article := &model.Content{}
	singleton.Orm().Where(map[string]interface{}{
		"ContentItemID":    ContentItemID,
		"ContentSubTypeID": 0,
	}).First(article)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentByContentItemIDAndTitle(ContentItemID uint, Title string) *model.Content {
	article := &model.Content{}
	singleton.Orm().Where(map[string]interface{}{
		"ContentItemID": ContentItemID,
		"Title":         Title,
	}).First(article)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentByContentItemIDAndContentSubTypeID(ContentItemID, ContentSubTypeID types.PrimaryKey) model.Content {
	article := model.Content{}
	singleton.Orm().Where(map[string]interface{}{
		"ContentItemID":    ContentItemID,
		"ContentSubTypeID": ContentSubTypeID,
	}).First(&article)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentByID(ID types.PrimaryKey) *model.Content {
	article := &model.Content{}
	err := service.Get(singleton.Orm(), ID, article) //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentAndAddLook(context *gweb.Context, ArticleID types.PrimaryKey) *model.Content {

	article := &model.Content{}
	err := service.Get(singleton.Orm(), ArticleID, article) //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)

	if context.Session.Attributes.Get(gweb.AttributesKey(strconv.Itoa(int(ArticleID)))) == nil {
		context.Session.Attributes.Put(gweb.AttributesKey(strconv.Itoa(int(ArticleID))), "Look")
		service.ChangeMap(singleton.Orm(), ArticleID, &model.Content{}, map[string]interface{}{"Look": article.Look + 1})

		LookArticle := 0 //todo config.Config.LookArticle

		if context.Session.Attributes.Get(play.SessionUser) != nil {
			user := context.Session.Attributes.Get(play.SessionUser).(*model.User)
			err := service.Journal.AddScoreJournal(singleton.Orm(),
				user.ID,
				"看文章送积分", "看文章/"+strconv.Itoa(int(article.ID)),
				play.ScoreJournal_Type_Look_Article, int64(LookArticle), extends.KV{Key: "ArticleID", Value: article.ID})
			glog.Error(err)
		}

	}
	return article
}

func (service ContentService) HaveContentByTitle(ContentItemID, ContentSubTypeID uint, Title string) bool {
	Orm := singleton.Orm()
	_article := &model.Content{}
	Orm.Where("ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeID).Where("Title=?", Title).First(_article)
	if _article.ID == 0 {
		return false
	} else {
		return true
	}

}
func (service ContentService) FindContentByIDAndNum(contentItemIDList []types.PrimaryKey, num int) []model.Content {
	Orm := singleton.Orm()
	_articleList := make([]model.Content, 0)
	Orm.Where(`"ContentItemID" in ?`, contentItemIDList).Order(`"CreatedAt" desc`).Limit(num).Find(&_articleList)
	return _articleList
}
func (service ContentService) AddContent(article *model.Content) *result.ActionResult {

	as := &result.ActionResult{}
	Orm := singleton.Orm()

	if article.ContentItemID == 0 {
		as.Code = result.Fail
		as.Message = "必须指定ContentItemID"
		return as
	}

	contentItem := service.GetContentItemByID(article.ContentItemID)
	switch contentItem.Type {

	case model.ContentTypeContent:
		if article.ContentSubTypeID == 0 {
			//return &result.ActionResult{Code: result.Fail, Message: fmt.Sprintf("%v内容请指定类型", contentItem.Type)}
		} else {
			contentSubType := service.GetContentSubTypeByID(article.ContentSubTypeID)

			if contentSubType.ID == 0 {
				return &result.ActionResult{Code: result.Fail, Message: fmt.Sprintf("无效的类别%v", contentSubType.ID)}
			}

			content := service.GetContentByContentItemIDAndContentSubTypeID(article.ContentItemID, article.ContentSubTypeID)
			if content.ID > 0 && article.ID != content.ID {
				return &result.ActionResult{Code: result.Fail, Message: fmt.Sprintf("添加的内容与原内容冲突")}
			}
		}

	}

	//_article := &model.Article{}
	//err := Orm.Where("ContentItemID=? and ContentSubTypeID=?", article.ContentItemID, article.ContentSubTypeID).Where("Title=?", article.Title).First(_article).Error
	//if _article.ID != 0 && _article.ID != article.ID {
	if false {
		as.Code = result.Fail
		as.Message = "添加失败，存在相同的标题"
	} else {
		articleID := article.ID
		var err error
		if articleID == 0 {
			err = service.Save(Orm, article) //self.model.AddArticle(Orm, article)
		} else {
			err = service.ChangeMap(Orm, articleID, &model.Content{}, map[string]interface{}{
				"Author":           article.Author,
				"Content":          article.Content,
				"ContentSubTypeID": article.ContentSubTypeID,
				"Description":      article.Description,
				"FromUrl":          article.FromUrl,
				"Introduce":        article.Introduce,
				"Keywords":         article.Keywords,
				"Picture":          article.Picture,
				"Title":            article.Title,
			})

		}

		if err != nil {
			glog.Error(err)
			as.Code = result.Fail
			as.Message = err.Error()
		} else {
			as.Code = result.Success
			as.Data = article
			if articleID != 0 {
				as.Message = "修改成功"
			} else {
				as.Message = "添加成功"
			}

		}
	}
	return as
}

func (service ContentService) GalleryBlock(OID types.PrimaryKey, num int) ([]model.ContentItem, []model.Content) {
	contentItemList := service.FindContentItemByType(model.ContentTypeGallery, OID)

	contentItemIDList := make([]types.PrimaryKey, 0)
	for _, item := range contentItemList {
		contentItemIDList = append(contentItemIDList, item.ID)
	}
	contentList := service.FindContentByIDAndNum(contentItemIDList, num)
	return contentItemList, contentList
}
