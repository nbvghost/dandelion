package content

import (
	"encoding/json"
	"fmt"
	"github.com/nbvghost/dandelion/library/db"
	"go.uber.org/zap"
	"log"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/domain/tag"
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/play"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/server/redis"
	"github.com/pkg/errors"

	"github.com/nbvghost/tool/object"
)

func (service ContentService) HotViewList(OID, ContentItemID dao.PrimaryKey, count uint) []model.Content {
	Orm := db.Orm()
	var result []model.Content
	db := Orm.Model(&model.Content{}).Where(map[string]interface{}{"OID": OID}).Where(`"ContentItemID"=?`, ContentItemID).Order(`"CountView" desc`).Limit(int(count))
	db.Find(&result)
	return result
}
func (service ContentService) HotLikeList(OID, ContentItemID dao.PrimaryKey, count uint) []model.Content {
	Orm := db.Orm()
	var result []model.Content
	db := Orm.Model(&model.Content{}).Where(map[string]interface{}{"OID": OID}).Where(`"ContentItemID"=?`, ContentItemID).Order(`"CountLike" desc`).Limit(int(count))
	db.Find(&result)
	return result
}
func (service ContentService) SortList(OID, ContentItemID dao.PrimaryKey, sort string, sortMethod int, count uint) []model.Content {
	Orm := db.Orm()
	var result []model.Content
	db := Orm.Model(&model.Content{}).Where(map[string]interface{}{"OID": OID}).Where(`"ContentItemID"=?`, ContentItemID)
	if sortMethod >= 0 {
		db = db.Order(fmt.Sprintf(`"%s" asc`, sort))
	} else {
		db = db.Order(fmt.Sprintf(`"%s" desc`, sort))
	}
	db.Limit(int(count)).Find(&result)
	return result
}

func (service ContentService) FindContentByTag(OID dao.PrimaryKey, tag extends.Tag, _pageIndex int, orders ...extends.Order) (pageIndex, pageSize int, total int64, list []*model.Content, err error) {
	//select * from "Content" where array_length("Tags",1) is null;
	db := db.Orm().Model(model.Content{}).Where(`"OID"=?`, OID).
		Where(`array_length("Tags",1) is not null`).
		Where(`"Tags" @> array[?]`, tag.Name)

	db.Count(&total)

	for _, v := range orders {
		db.Order(fmt.Sprintf(`"%s" %s`, v.ColumnName, v.Method))
	}

	pageSize = 20

	err = db.Limit(pageSize).Offset(_pageIndex * pageSize).Find(&list).Error
	pageIndex = _pageIndex

	return
}
func (service ContentService) FindContentTags(OID dao.PrimaryKey) ([]extends.Tag, error) {
	//SELECT unnest("Tags") as Tag,count("Tags") as Count FROM "Content" where  group by unnest("Tags");
	var tags []extends.Tag
	err := db.Orm().Model(model.Content{}).Select(`unnest("Tags") as "Name",count("Tags") as "Count"`).Where(map[string]interface{}{
		"OID": OID,
	}).Where(`array_length("Tags",1)>0`).Group(`unnest("Tags")`).Find(&tags).Error
	tags = tag.CreateUri(tags)
	return tags, err
}
func (service ContentService) FindContentTagsByContentItemID(OID, ContentItemID dao.PrimaryKey) []extends.Tag {
	//SELECT unnest("Tags") as Tag,count("Tags") as Count FROM "Content" where  group by unnest("Tags");
	var tags []extends.Tag
	db.Orm().Model(model.Content{}).Select(`unnest("Tags") as "Name",count("Tags") as "Count"`).Where(map[string]interface{}{
		"OID":           OID,
		"ContentItemID": ContentItemID,
	}).Where(`array_length("Tags",1)>0`).Group(`unnest("Tags")`).Find(&tags)
	tags = tag.CreateUri(tags)
	return tags
}
func (service ContentService) PaginationContent(OID, ContentItemID, ContentSubTypeID dao.PrimaryKey, pageIndex int, pageSize int) (int, int, int, []*model.Content) {
	if ContentItemID == 0 {
		db := dao.Find(db.Orm(), &model.Content{}).Where(`"OID"=?`, OID)
		total := db.Limit(pageIndex, pageSize)
		mlist := db.List()
		list := make([]*model.Content, 0)
		for i := range mlist {
			list = append(list, mlist[i].(*model.Content))
		}
		return pageIndex, pageSize, int(total), list
	}
	if ContentSubTypeID == 0 {
		db := dao.Find(db.Orm(), &model.Content{}).Where(`"OID"=? and "ContentItemID"=?`, OID, ContentItemID)
		total := db.Limit(pageIndex, pageSize)
		mlist := db.List()
		list := make([]*model.Content, 0)
		for i := range mlist {
			list = append(list, mlist[i].(*model.Content))
		}
		return pageIndex, pageSize, int(total), list
		//return repository.Content.FindByOIDAndContentItemIDLimit(OID, ContentItemID, params.NewLimit(pageIndex, 20))
	}
	contentSubTypeIDs := service.GetContentSubTypeAllIDByID(ContentItemID, ContentSubTypeID)

	var total int64
	db := db.Orm().Model(model.Content{}).Where(map[string]interface{}{
		"OID":           OID,
		"ContentItemID": ContentItemID,
	}).Where(`"ContentSubTypeID" in (?)`, contentSubTypeIDs)
	db.Count(&total)

	var list []*model.Content
	db.Limit(pageSize).Offset(pageIndex * pageSize).Find(&list)

	return pageIndex, pageSize, int(total), list
}
func (service ContentService) GetContentTypeByID(OID dao.PrimaryKey, ContentItemID, ContentSubTypeID dao.PrimaryKey) (model.ContentItem, model.ContentSubType) {
	Orm := db.Orm()
	var item model.ContentItem
	var itemSub model.ContentSubType

	itemMap := map[string]interface{}{"OID": OID, "ID": ContentItemID}
	Orm.Model(model.ContentItem{}).Where(itemMap).First(&item)

	itemSubMap := map[string]interface{}{
		"OID":           OID,
		"ContentItemID": item.ID,
		"ID":            ContentSubTypeID,
	}
	Orm.Model(model.ContentSubType{}).Where(itemSubMap).First(&itemSub)
	return item, itemSub
}

// uri 和 name 在 ContentItemID 下面唯一
func (service ContentService) GetContentSubTypeByUri(OID, ContentItemID, ID dao.PrimaryKey, uri string) model.ContentSubType {
	Orm := db.Orm()
	var item model.ContentSubType
	Orm.Model(model.ContentSubType{}).Where(map[string]interface{}{
		"OID":           OID,
		"ContentItemID": ContentItemID,
		"Uri":           uri,
	}).Where(`"ID"<>?`, ID).First(&item)
	return item
}
func (service ContentService) SaveContentSubType(OID dao.PrimaryKey, item *model.ContentSubType) error {
	Orm := db.Orm()
	mm := service.GetContentSubTypeByName(OID, item.ContentItemID, item.ID, item.Name)
	if !mm.IsZero() {
		return errors.Errorf("名字重复")
	}

	uri := service.PinyinService.AutoDetectUri(item.Name)
	g := service.GetContentSubTypeByUri(OID, item.ContentItemID, item.ID, uri)
	if !g.IsZero() {
		uri = fmt.Sprintf("%s-%d", uri, time.Now().Unix())
	}
	item.Uri = uri
	item.OID = OID

	if item.IsZero() {
		item.Sort = time.Now().Unix()
		contentItem := service.GetContentItemByID(item.ContentItemID)
		if contentItem.IsZero() {
			return errors.Errorf("类型不存在:%d", item.ContentItemID)
		}
		if !item.ParentContentSubTypeID.IsZero() {
			cst := service.GetContentSubTypeByID(item.ParentContentSubTypeID)
			if cst.IsZero() {
				return errors.Errorf("父类不存在:%d", item.ContentItemID)
			}
		}
		err := dao.Create(Orm, item)
		if err != nil {
			return &result.ActionResult{
				Code:    result.SQLError,
				Message: err.Error(),
				Data:    nil,
			}
		}
	} else {
		return dao.UpdateByPrimaryKey(Orm, &model.ContentSubType{}, dao.PrimaryKey(item.ID), &model.ContentSubType{Name: item.Name, Uri: item.Uri})
	}

	return nil
}
func (service ContentService) SaveContentItem(OID dao.PrimaryKey, item *model.ContentItem) error {
	Orm := db.Orm()

	if len(item.Name) == 0 {
		return errors.Errorf("请指定名称")
	}

	have := service.ExistContentItemByNameAndOID(OID, item.ID, item.Name)
	if !have.IsZero() {
		return &result.ActionResult{
			Code:    result.Fail,
			Message: "这个名字已经被使用了",
			Data:    nil,
		} //&gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("这个名字已经被使用了"), "", nil)}
	}

	item.OID = OID

	uri := service.PinyinService.AutoDetectUri(item.Name)
	g := service.getContentItemByUri(OID, item.ID, uri)
	if !g.IsZero() {
		uri = fmt.Sprintf("%s-%d", uri, time.Now().Unix())
	}
	item.Uri = uri

	if item.ID == 0 {
		contentItemList := service.ListContentItemByOID(OID)
		if len(contentItemList) > 0 {
			item.Sort = contentItemList[len(contentItemList)-1].Sort + 1
		}

		var mt model.ContentType
		Orm.Where(map[string]interface{}{"ID": item.ContentTypeID}).First(&mt)
		if mt.ID == 0 {
			return &result.ActionResult{
				Code:    result.Fail,
				Message: "没有找到类型",
				Data:    nil,
			}
		}
		if strings.EqualFold(string(mt.Type), string(model.ContentTypeBlog)) || strings.EqualFold(string(mt.Type), string(model.ContentTypeProducts)) {
			haveList := service.FindContentItemByType(mt.Type, OID)
			if len(haveList) != 0 {
				return &result.ActionResult{
					Code:    result.Fail,
					Message: fmt.Sprintf("这个类型（%v）只能创建一个", have.Name),
					Data:    nil,
				}
			}
		}
		item.Type = mt.Type
		err := dao.Create(Orm, item)
		if err != nil {
			return &result.ActionResult{
				Code:    result.SQLError,
				Message: err.Error(),
				Data:    nil,
			}
		}
	} else {
		err := Orm.Model(model.ContentItem{}).Where(`"ID"=?`, item.ID).Updates(map[string]interface{}{
			"Name":         item.Name,
			"Uri":          item.Uri,
			"Image":        item.Image,
			"Badge":        item.Badge,
			"Introduction": item.Introduction,
		}).Error
		if err != nil {
			return &result.ActionResult{
				Code:    result.SQLError,
				Message: err.Error(),
				Data:    nil,
			}
		}
	}

	return &result.ActionResult{
		Code:    result.Success,
		Message: "保存成功",
		Data:    nil,
	}
}
func (service ContentService) ChangeContentConfig(OID dao.PrimaryKey, fieldName, fieldValue string) error {

	changeMap := make(map[string]interface{})

	switch fieldName {
	case "Name":
		changeMap["Name"] = fieldValue
	case "Logo":
		changeMap["Logo"] = fieldValue
	case "FaviconIco":
		changeMap["FaviconIco"] = fieldValue
	case "SocialAccount":
		var socialAccount sqltype.SocialAccountList
		err := json.Unmarshal([]byte(fieldValue), &socialAccount)
		if err != nil {
			return err
		}
		changeMap["SocialAccount"] = socialAccount
	case "CustomerService":
		var customerService sqltype.CustomerServiceList
		err := json.Unmarshal([]byte(fieldValue), &customerService)
		if err != nil {
			return err
		}
		changeMap["CustomerService"] = customerService
	case "EnableHTMLCache":
		EnableHTMLCache, _ := strconv.ParseBool(fieldValue)
		changeMap["EnableHTMLCache"] = EnableHTMLCache
	case "FocusPicture":
		var focusPicture sqltype.FocusPictureList
		err := json.Unmarshal([]byte(fieldValue), &focusPicture)
		if err != nil {
			return err
		}
		changeMap["FocusPicture"] = focusPicture

	}
	Orm := db.Orm()
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

func (service ContentService) GetContentConfig(orm *gorm.DB, OID dao.PrimaryKey) model.ContentConfig {
	var contentConfig model.ContentConfig
	orm.Model(&model.ContentConfig{}).Where(map[string]interface{}{"OID": OID}).First(&contentConfig)
	return contentConfig
}

//-----------------------------------------Content----------------------------------------------------------

func (service ContentService) GetContentItemByIDAndOID(ID, OID uint) model.ContentItem {
	Orm := db.Orm()
	var menus model.ContentItem

	Orm.Where(`"ID"=? and "OID"=?`, ID, OID).First(&menus)

	return menus
}
func (service ContentService) GetContentItemByID(ID dao.PrimaryKey) model.ContentItem {
	Orm := db.Orm()
	var menus model.ContentItem
	Orm.Where(`"ID"=?`, ID).First(&menus)
	return menus
}
func (service ContentService) GetContentSubTypeByID(ID dao.PrimaryKey) model.ContentSubType {
	Orm := db.Orm()
	var menus model.ContentSubType
	Orm.Where(`"ID"=?`, ID).First(&menus)
	return menus
}

func (service ContentService) ExistContentItemByNameAndOID(OID, ID dao.PrimaryKey, Name string) model.ContentItem {
	Orm := db.Orm()
	var menus model.ContentItem
	Orm.Where(`"OID"=?`, OID).Where(map[string]interface{}{"Name": Name}).Where(`"ID"<>?`, ID).First(&menus)
	return menus
}

func (service ContentService) FindContentItemByType(Type model.ContentTypeType, OID dao.PrimaryKey) []model.ContentItem {
	Orm := db.Orm()
	menus := make([]model.ContentItem, 0)
	Orm.Where(map[string]interface{}{
		"Type": Type,
		"OID":  OID,
	}).Find(&menus)
	return menus
}

func (service ContentService) ListContentType() []dao.IEntity {
	Orm := db.Orm()
	return dao.Find(Orm, &model.ContentType{}).List()
}
func (service ContentService) ListContentTypeByType(Type string) *model.ContentType {
	//Orm := singleton.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//var list model.ContentType
	item := dao.GetBy(db.Orm(), &model.ContentType{}, map[string]any{"Type": Type}).(*model.ContentType) //service.FindWhere(Orm, &list, "Type=?", Type)
	return item
}
func (service ContentService) FindContentSubTypesByNameAndContentItemID(Name string, ContentItemID dao.PrimaryKey) model.ContentSubType {
	Orm := db.Orm()
	var cst model.ContentSubType
	Orm.Where("ContentItemID=? and Name=?", ContentItemID, Name).First(&cst)
	return cst
}

// AddSpiderContent -----------------------
func (service ContentService) AddSpiderContent(OID dao.PrimaryKey, ContentName string, ContentSubTypeName string, Author, Title string, FromUrl string, Introduce string, Picture string, Content string, CreatedAt time.Time) error {
	var article model.Content
	article.Title = Title
	article.FromUrl = FromUrl
	article.CreatedAt = CreatedAt
	article.UpdatedAt = CreatedAt

	IntroduceRune := []rune(Introduce)
	if len(IntroduceRune) > 255 {
		article.Summary = string(IntroduceRune[:255])
	} else {
		article.Summary = Introduce
	}

	//Picture=tool.DownloadInternetImage(Picture,"Mozilla/5.0 (Linux; Android 7.0; SLA-AL00 Build/HUAWEISLA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/6.2 TBS/044109 Mobile Safari/537.36 MicroMessenger/6.6.7.1321(0x26060739) NetType/WIFI Language/zh_CN",weixin_tmp_url)
	article.Picture = Picture
	article.Content = Content

	contentType := service.ListContentTypeByType(play.ContentTypeArticles)

	content := service.ExistContentItemByNameAndOID(OID, 0, ContentName)
	if content.ID == 0 {
		content.OID = OID
		content.Type = contentType.Type
		content.Name = ContentName
		content.ContentTypeID = contentType.ID
		err := dao.Save(db.Orm(), &content)
		if err != nil {
			return err
		}

	}

	article.ContentItemID = content.ID
	contentSubType := service.FindContentSubTypesByNameAndContentItemID(ContentSubTypeName, content.ID)
	if contentSubType.ID == 0 {
		contentSubType.Name = ContentSubTypeName
		contentSubType.ContentItemID = content.ID
		err := dao.Save(db.Orm(), &contentSubType)
		if err != nil {
			return err
		}
	}

	article.Author = Author
	article.ContentSubTypeID = contentSubType.ID
	err := service.SaveContent(OID, &article)
	if err != nil {
		return err
	}
	return nil
}

func (service ContentService) ChangeContent(article *model.Content) error {

	return dao.Save(db.Orm(), article)
}

func (service ContentService) GetContentByTitle(Orm *gorm.DB, OID dao.PrimaryKey, Title string) *model.Content {
	article := &model.Content{}
	err := Orm.Where(`"OID"=?`, OID).Where(`"Title"=?`, Title).First(article).Error //SelectOne(user, "select * from User where Email=?", Email)
	if err != nil {
		log.Println(err)
	}
	return article
}
func (service ContentService) DelContent(ID dao.PrimaryKey) error {
	err := dao.DeleteByPrimaryKey(db.Orm(), &model.Content{}, ID)
	return err
}

func (service ContentService) FindContentByContentSubTypeID(ContentSubTypeID dao.PrimaryKey) []dao.IEntity {
	//var contentList []model.Content
	//err := service.FindWhere(singleton.Orm(), &contentList, "ContentSubTypeID=?", ContentSubTypeID) //SelectOne(user, "select * from User where Email=?", Email)

	contentList := dao.Find(db.Orm(), &model.Content{}).Where(`"ContentSubTypeID"=?`, ContentSubTypeID).List()

	return contentList
}
func (service ContentService) FindContentByContentItemIDAndContentSubTypeID(ContentItemID uint, ContentSubTypeID uint) *model.Content {
	//service.FindWhere(singleton.Orm(), &content, "ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeID) //SelectOne(user, "select * from User where Email=?", Email)
	content := dao.GetBy(db.Orm(), &model.Content{}, map[string]any{"ContentItemID": ContentItemID, "ContentSubTypeID": ContentSubTypeID}).(*model.Content)
	return content
}

func (service ContentService) FindContentByTypeID(menusData *extends.MenusData, ContentItemID, ContentSubTypeID, ContentSubTypeChildID dao.PrimaryKey) model.Content {

	var content model.Content

	/*if ContentItemID == 0 {
		log.Println("参数ContentItemID为0")
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

		}else if ContentSubTypeID > 0 {
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
func (service ContentService) FindContentListByTypeID(menusData *extends.MenusData, ContentItemID, ContentSubTypeID, ContentSubTypeChildID dao.PrimaryKey, _Page int, _Limit int) result.Pager {

	var pager result.Pager

	if ContentItemID == 0 {
		log.Println("参数ContentItemID为0")
		return pager
	}

	if ContentSubTypeID == 0 && ContentSubTypeChildID == 0 {
		if len(menusData.List) > 0 {

		}
		db := db.Orm().Model(&model.Content{}).Where(`"ContentItemID"=?`, ContentItemID).
			Order(`"CreatedAt" desc`).Order(`"ID" desc`)
		return model.Paging(db, _Page, _Limit, model.Content{})
	} else {

		if ContentSubTypeChildID > 0 {
			db := db.Orm().Model(&model.Content{}).Where(`"ContentItemID"=? and "ContentSubTypeID"=?`, ContentItemID, ContentSubTypeChildID).
				Order(`"CreatedAt" desc`).Order(`"ID" desc`)
			return model.Paging(db, _Page, _Limit, model.Content{})
		} else {

			ContentSubTypeIDList := service.GetContentSubTypeAllIDByID(ContentItemID, ContentSubTypeID)

			db := db.Orm().Model(&model.Content{}).
				Where(`"ContentItemID"=? and "ContentSubTypeID" in (?)`, ContentItemID, ContentSubTypeIDList).
				Order(`"CreatedAt" desc`).Order(`"ID" desc`)
			return model.Paging(db, _Page, _Limit, model.Content{})
		}
	}

}

func (service ContentService) FindContentListForLeftRight(ContentItemID, ContentSubTypeID dao.PrimaryKey, ContentID dao.PrimaryKey, ContentCreatedAt time.Time) [2]*model.Content {
	var contentList [2]*model.Content
	if ContentItemID == 0 {
		log.Println("参数ContentItemID为0")
		return contentList
	}

	var ContentSubTypeIDList []dao.PrimaryKey
	if ContentSubTypeID > 0 {
		ContentSubTypeIDList = service.GetContentSubTypeAllIDByID(ContentItemID, ContentSubTypeID)
	}

	ContentSubTypeIDListStr := make([]string, 0)
	for index := range ContentSubTypeIDList {
		ContentSubTypeIDListStr = append(ContentSubTypeIDListStr, object.ParseString(ContentSubTypeIDList[index]))
	}

	var whereSql = ""

	if len(ContentSubTypeIDList) > 0 {
		whereSql = fmt.Sprintf(`"ContentItemID"=%v and "ContentSubTypeID" in (%v)`, ContentItemID, strings.Join(ContentSubTypeIDListStr, ","))
	} else {
		whereSql = fmt.Sprintf(`"ContentItemID"=%v`, ContentItemID)
	}

	var left model.Content
	var right model.Content
	err := db.Orm().Raw(`SELECT * FROM "Content" WHERE `+whereSql+` and "ID"<>? and "CreatedAt">=? ORDER BY "CreatedAt","ID" limit 1`, ContentID, ContentCreatedAt).Scan(&left).Error
	if err != nil {
		log.Println(err)
	}
	err = db.Orm().Raw(`SELECT * FROM "Content" WHERE `+whereSql+` and "ID"<>? and "CreatedAt"<=? ORDER BY "CreatedAt" desc,"ID" desc limit 1`, ContentID, ContentCreatedAt).Scan(&right).Error
	if err != nil {
		log.Println(err)
	}

	return [2]*model.Content{&left, &right}
}

func (service ContentService) GetContentByContentItemID(ContentItemID uint) *model.Content {
	article := &model.Content{}
	db.Orm().Where(map[string]interface{}{
		"ContentItemID":    ContentItemID,
		"ContentSubTypeID": 0,
	}).First(article)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentByContentItemIDAndTitle(ContentItemID uint, Title string) *model.Content {
	article := &model.Content{}
	db.Orm().Where(map[string]interface{}{
		"ContentItemID": ContentItemID,
		"Title":         Title,
	}).First(article)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentByContentItemIDAndContentSubTypeID(ContentItemID, ContentSubTypeID dao.PrimaryKey) model.Content {
	article := model.Content{}
	db.Orm().Where(map[string]interface{}{
		"ContentItemID":    ContentItemID,
		"ContentSubTypeID": ContentSubTypeID,
	}).First(&article)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentByUri(OID dao.PrimaryKey, Uri string) *model.Content {
	article := &model.Content{}
	//err := service.Get(singleton.Orm(), ID, article) //SelectOne(user, "select * from User where Email=?", Email)
	db.Orm().Model(model.Content{}).Where(`"OID"=?`, OID).Where(`"Uri"=?`, Uri).First(article)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentByID(ID dao.PrimaryKey) *model.Content {

	article := dao.GetByPrimaryKey(db.Orm(), &model.Content{}, ID).(*model.Content) //SelectOne(user, "select * from User where Email=?", Email)

	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentAndAddLook(ctx constrain.IContext, ArticleID dao.PrimaryKey) *model.Content {

	article := dao.GetByPrimaryKey(db.Orm(), &model.Content{}, ArticleID).(*model.Content) //SelectOne(user, "select * from User where Email=?", Email)

	lc, _ := ctx.Redis().Get(ctx, redis.NewArticleLookCount(ctx.UID(), ArticleID))

	if len(lc) == 0 {
		now := time.Now()
		tomorrowTime := now.Add(24 * time.Hour)
		endDayTime := time.Date(tomorrowTime.Year(), tomorrowTime.Month(), tomorrowTime.Day(), 0, 0, 0, 0, tomorrowTime.Location())
		//context.Session.Attributes.Put(gweb.AttributesKey(strconv.Itoa(int(ArticleID))), "CountView")
		ctx.Redis().Set(ctx, redis.NewArticleLookCount(ctx.UID(), ArticleID), "true", endDayTime.Sub(now))
		err := dao.UpdateByPrimaryKey(db.Orm(), &model.Content{}, ArticleID, map[string]interface{}{"CountView": article.CountView + 1})
		if err != nil {
			ctx.Logger().Error("GetContentAndAddLook", zap.Error(err))
		}

		LookArticle := 0 //todo config.Config.LookArticle

		//if context.Session.Attributes.Get(play.SessionUser) != nil {
		//user := context.Session.Attributes.Get(play.SessionUser).(*model.User)
		//err = service.Journal.AddScoreJournal(db.Orm(), ctx.UID(), "看文章送积分", "看文章/"+strconv.Itoa(int(article.ID)), play.ScoreJournal_Type_Look_Article, int64(LookArticle), extends.KV{Key: "ArticleID", Value: article.ID})
		err = service.Journal.AddScoreJournal(db.Orm(), ctx.UID(), "看文章送积分", "看文章/"+strconv.Itoa(int(article.ID)), model.ScoreJournal_Type_Look_Article, int64(LookArticle))
		if err != nil {
			ctx.Logger().Error("GetContentAndAddLook", zap.Error(err))
		}
		//}

	}
	return article
}

func (service ContentService) HaveContentByTitle(ContentItemID, ContentSubTypeID uint, Title string) bool {
	Orm := db.Orm()
	_article := &model.Content{}
	Orm.Where("ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeID).Where("Title=?", Title).First(_article)
	if _article.ID == 0 {
		return false
	} else {
		return true
	}

}
func (service ContentService) FindContentByIDAndNum(contentItemIDList []dao.PrimaryKey, num int) []model.Content {
	Orm := db.Orm()
	_articleList := make([]model.Content, 0)
	Orm.Where(`"ContentItemID" in ?`, contentItemIDList).Order(`"CreatedAt" desc`).Limit(num).Find(&_articleList)
	return _articleList
}
func (service ContentService) getContentByUri(OID dao.PrimaryKey, uri string) model.Content {
	Orm := db.Orm()
	var item model.Content
	Orm.Model(model.Content{}).Where(map[string]interface{}{"OID": OID, "Uri": uri}).First(&item)
	return item
}
func (service ContentService) getContentItemByUri(OID, ID dao.PrimaryKey, uri string) model.ContentItem {
	Orm := db.Orm()
	var item model.ContentItem
	item.OID = OID
	item.Uri = uri
	Orm.Model(model.ContentItem{}).Where(map[string]interface{}{"OID": item.OID, "Uri": item.Uri}).Where(`"ID"<>?`, ID).First(&item)
	return item
}
func (service ContentService) SaveContent(OID dao.PrimaryKey, article *model.Content) error {
	Orm := db.Orm()
	if article.ContentItemID == 0 {
		return errors.Errorf("必须指定ContentItemID")
	}

	contentItem := service.GetContentItemByID(article.ContentItemID)
	switch contentItem.Type {

	case model.ContentTypeContent:
		if article.ContentSubTypeID == 0 {
			//return &result.ActionResult{Code: result.Fail, Message: fmt.Sprintf("%v内容请指定类型", contentItem.Type)}
		} else {
			contentSubType := service.GetContentSubTypeByID(article.ContentSubTypeID)

			if contentSubType.ID == 0 {
				return errors.Errorf("无效的类别%v", contentSubType.ID)
			}

			content := service.GetContentByContentItemIDAndContentSubTypeID(article.ContentItemID, article.ContentSubTypeID)
			if content.ID > 0 && article.ID != content.ID {
				return errors.Errorf("添加的内容与原内容冲突")
			}
		}
	}

	g := service.GetContentByTitle(Orm, OID, article.Title)
	if !g.IsZero() && g.ID != article.ID {
		return errors.Errorf("添加失败,存在相同的标题")
	}
	if g.IsZero() {
		g = service.GetContentByID(article.ID)
	}

	article.OID = OID
	if len(g.Uri) == 0 {
		uri := service.PinyinService.AutoDetectUri(article.Title)
		hasContent := service.getContentByUri(OID, uri)
		if !hasContent.IsZero() && hasContent.ID != g.ID {
			return errors.Errorf("添加失败,存在相同的标题")
		}
		article.Uri = uri
	} else {
		article.Uri = g.Uri
	}
	articleID := article.ID
	var err error
	if articleID == 0 {
		err = dao.Create(Orm, article) //self.model.AddArticle(Orm, article)
	} else {
		err = dao.UpdateByPrimaryKey(Orm, &model.Content{}, articleID, map[string]interface{}{
			"Author":           article.Author,
			"Content":          article.Content,
			"ContentSubTypeID": article.ContentSubTypeID,
			"FromUrl":          article.FromUrl,
			"Summary":          article.Summary,
			"Picture":          article.Picture,
			"Title":            article.Title,
			"Tags":             article.Tags,
			"Uri":              article.Uri,
			"Images":           article.Images,
			"FieldGroupID":     article.FieldGroupID,
			"FieldData":        article.FieldData,
			"Keywords":         article.Keywords,
			"Description":      article.Description,
		})
	}
	return err
}

func (service ContentService) GalleryBlock(OID dao.PrimaryKey, num int) ([]model.ContentItem, []model.Content) {
	contentItemList := service.FindContentItemByType(model.ContentTypeGallery, OID)
	contentItemIDList := make([]dao.PrimaryKey, 0)
	for _, item := range contentItemList {
		contentItemIDList = append(contentItemIDList, item.ID)
	}
	contentList := service.FindContentByIDAndNum(contentItemIDList, num)
	return contentItemList, contentList
}

func (service ContentService) FindContentByTypeTemplate(oid dao.PrimaryKey, contentType string, templateName string, pageIndex int) (int64, []*model.Content) {
	var list []*model.Content
	var total int64

	d := db.Orm().Model(model.Content{}).Select(`"Content".*`).
		Joins(`left join "ContentItem" on "Content"."ContentItemID"="ContentItem"."ID"`).Order(`"Content"."CreatedAt" desc`).
		Where(`"Content"."OID"=? and "ContentItem"."Type"=? and "ContentItem"."TemplateName"=?`, oid, contentType, templateName)

	d.Count(&total)
	d.Offset(pageIndex * 20).Limit(20).Find(&list)
	return total, list
}

/*func (service ContentService) FindByOIDLimit(oid dao.PrimaryKey, pageIndex int, pageSize int) (int, int, int, []*model.Content) {
	return service.PaginationContent(oid, 0, 0, pageIndex, pageSize)
}
func (service ContentService) FindByOIDAndContentItemIDLimit(oid, contentItemID dao.PrimaryKey, pageIndex int, pageSize int) (int, int, int, []*model.Content) {
	return service.PaginationContent(oid, contentItemID, 0, pageIndex, pageSize)
}*/

//FindByOIDLimit                 func(OID dao.PrimaryKey, pagination *params.Limit) (pageIndex int, pageSize int, total int, list []*model.Content, err error)                `gpa:"AutoCreate"`
//FindByOIDAndContentItemIDLimit func(OID, ContentItemID dao.PrimaryKey, pagination *params.Limit) (pageIndex int, pageSize int, total int, list []*model.Content, err error) `gpa:"AutoCreate"`
