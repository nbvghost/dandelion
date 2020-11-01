package content

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/nbvghost/dandelion/app/play"
	"github.com/nbvghost/dandelion/app/result"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/dandelion/app/service/journal"
	"strings"
	"time"

	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
	"strconv"
)

type ContentService struct {
	dao.BaseDao
	Journal journal.JournalService
}

func (service ContentService) GetContentItemIDs(OID uint64) []uint64 {
	Orm := dao.Orm()
	var levea []uint64
	if OID <= 0 {
		return levea
	}
	Orm.Model(&dao.ContentItem{}).Where("OID=?", OID).Pluck("ID", &levea)
	return levea
}

func (service ContentService) GetClassifyByName(Name string, ContentItemID, ParentContentSubTypeID uint64) dao.ContentSubType {
	Orm := dao.Orm()
	var menus dao.ContentSubType

	Orm.Where("Name=?", Name).Where("ContentItemID=? and ParentContentSubTypeID=?", ContentItemID, ParentContentSubTypeID).First(&menus)

	return menus

}
func (service ContentService) FindAllContentSubType(OID uint64) []dao.ContentItemContentSubType {
	Orm := dao.Orm()
	var menus []dao.ContentItemContentSubType

	rows, err := Orm.Raw("SELECT Item.*,SubType.* FROM ContentSubType AS SubType RIGHT JOIN ContentItem AS Item ON (Item.ID=SubType.ContentItemID) WHERE Item.OID=? and Item.Hide=0 order by Item.Sort,Item.UpdatedAt desc", OID).Rows()

	if glog.Error(err) {
		return menus
	}
	defer rows.Close()

	for rows.Next() {

		var cc dao.ContentItemContentSubType
		glog.Error(Orm.ScanRows(rows, &cc))
		menus = append(menus, cc)

	}

	return menus

}
func (service ContentService) FindContentSubTypesByContentItemIDs(ContentItemIDs []uint64) []dao.ContentSubType {
	Orm := dao.Orm()
	var menus []dao.ContentSubType
	Orm.Model(dao.ContentSubType{}).Where("ContentItemID in (?)", ContentItemIDs).Order("Sort asc").Find(&menus)
	return menus
}
func (service ContentService) FindContentSubTypesByContentItemID(ContentItemID uint64) []dao.ContentSubType {
	Orm := dao.Orm()
	var menus []dao.ContentSubType
	Orm.Model(dao.ContentSubType{}).Where("ContentItemID=? and ParentContentSubTypeID=0", ContentItemID).Order("Sort asc").Find(&menus)
	return menus
}
func (service ContentService) FindContentSubTypesByParentContentSubTypeID(ParentContentSubTypeID uint64) []dao.ContentSubType {
	Orm := dao.Orm()
	var menus []dao.ContentSubType
	Orm.Model(dao.ContentSubType{}).Where("ParentContentSubTypeID=?", ParentContentSubTypeID).Order("Sort asc").Find(&menus)
	return menus
}
func (service ContentService) FindContentSubTypesByContentItemIDAndParentContentSubTypeID(ContentItemID, ParentContentSubTypeID uint64) []dao.ContentSubType {
	Orm := dao.Orm()
	var menus []dao.ContentSubType
	Orm.Model(dao.ContentSubType{}).Where("ContentItemID=? and ParentContentSubTypeID=?", ContentItemID, ParentContentSubTypeID).Order("Sort asc").Find(&menus)
	return menus
}

//-----------------------------------------Content----------------------------------------------------------

func (service ContentService) GetContentItemByIDAndOID(ID, OID uint64) dao.ContentItem {
	Orm := dao.Orm()
	var menus dao.ContentItem

	Orm.Where("ID=? and OID=?", ID, OID).First(&menus)

	return menus
}
func (service ContentService) GetContentItemByID(ID uint64) dao.ContentItem {
	Orm := dao.Orm()
	var menus dao.ContentItem
	Orm.Where("ID=?", ID).First(&menus)
	return menus
}
func (service ContentService) GetContentSubTypeByID(ID uint64) dao.ContentSubType {
	Orm := dao.Orm()
	var menus dao.ContentSubType
	Orm.Where("ID=?", ID).First(&menus)
	return menus
}
func (service ContentService) ListContentItemByOID(OID uint64) []dao.ContentItem {
	Orm := dao.Orm()
	var menus []dao.ContentItem
	Orm.Model(dao.ContentItem{}).Where("OID=?", OID).Order("Sort").Order("UpdatedAt desc").Find(&menus)
	return menus
}
func (service ContentService) GetContentItemByNameAndOID(Name string, OID uint64) dao.ContentItem {
	Orm := dao.Orm()
	var menus dao.ContentItem

	Orm.Where("Name=? and OID=?", Name, OID).First(&menus)

	return menus
}

func (service ContentService) GetContentItemByType(Type dao.ContentTypeType, OID uint64) dao.ContentItem {
	Orm := dao.Orm()
	var menus dao.ContentItem
	Orm.Where("Type=? and OID=?", Type, OID).First(&menus)
	return menus
}

func (service ContentService) ListContentType() []dao.ContentType {
	Orm := dao.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	var list []dao.ContentType
	err := service.FindAll(Orm, &list)
	glog.Trace(err)
	return list
}
func (service ContentService) ListContentTypeByType(Type string) dao.ContentType {
	Orm := dao.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	var list dao.ContentType
	err := service.FindWhere(Orm, &list, "Type=?", Type)
	glog.Trace(err)
	return list
}
func (service ContentService) FindContentSubTypesByNameAndContentItemID(Name string, ContentItemID uint64) dao.ContentSubType {
	Orm := dao.Orm()
	var cst dao.ContentSubType
	Orm.Where("ContentItemID=? and Name=?", ContentItemID, Name).First(&cst)
	return cst
}

//-----------------------
func (service ContentService) AddSpiderContent(OID uint64, ContentName string, ContentSubTypeName string, Author, Title string, FromUrl string, Introduce string, Picture string, Content string, CreatedAt time.Time) {
	var article dao.Content
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
		service.Save(dao.Orm(), &content)

	}

	article.ContentItemID = content.ID
	contentSubType := service.FindContentSubTypesByNameAndContentItemID(ContentSubTypeName, content.ID)
	if contentSubType.ID == 0 {
		contentSubType.Name = ContentSubTypeName
		contentSubType.ContentItemID = content.ID
		service.Save(dao.Orm(), &contentSubType)
	}

	article.Author = Author
	article.ContentSubTypeID = contentSubType.ID
	service.AddContent(&article)

}

func (service ContentService) ChangeContent(article *dao.Content) error {

	return service.Save(dao.Orm(), article)
}

func (service ContentService) GetContentByTitle(Orm *gorm.DB, Title string) *dao.Content {
	article := &dao.Content{}
	err := Orm.Where("Title=?", Title).First(article).Error //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)
	return article
}
func (service ContentService) DelContent(ID uint64) error {
	err := service.Delete(dao.Orm(), &dao.Content{}, ID)
	return err
}

//获取ID，返回子类ID,包括本身
func (service ContentService) GetContentSubTypeAllIDByID(ContentItemID, ContentSubTypeID uint64) []uint64 {
	var IDList []uint64
	err := dao.Orm().Model(&dao.ContentSubType{}).Where("ContentItemID=? and (ID=? or ParentContentSubTypeID=?)", ContentItemID, ContentSubTypeID, ContentSubTypeID).Pluck("ID", &IDList).Error
	glog.Error(err)
	return IDList
}
func (service ContentService) FindContentByContentSubTypeID(ContentSubTypeID uint64) []dao.Content {
	var contentList []dao.Content
	err := service.FindWhere(dao.Orm(), &contentList, "ContentSubTypeID=?", ContentSubTypeID) //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)
	return contentList
}
func (service ContentService) FindContentByContentItemIDAndContentSubTypeID(ContentItemID uint64, ContentSubTypeID uint64) dao.Content {

	var content dao.Content
	if ContentItemID == 0 {
		glog.Trace("参数ContentItemID为0")
		return content
	}
	if ContentSubTypeID == 0 {
		glog.Trace("参数ContentSubTypeID为0")
		return content
	}

	service.FindWhere(dao.Orm(), &content, "ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeID) //SelectOne(user, "select * from User where Email=?", Email)

	return content
}

//ContentItemID
//ContentSubTypeID
//ContentSubTypeChildID
func (service ContentService) FindContentByTypeID(menusData *dao.MenusData, ContentItemID, ContentSubTypeID, ContentSubTypeChildID uint64) dao.Content {

	var content dao.Content
	if ContentItemID == 0 {
		glog.Trace("参数ContentItemID为0")
		return content
	}

	if ContentSubTypeID == 0 && ContentSubTypeChildID == 0 {

		dao.Orm().Model(&dao.Content{}).
			Where("ContentItemID=?", ContentItemID).
			Where("ContentSubTypeID=?", 0).
			Order("CreatedAt desc").Order("ID desc").First(&content)

		if content.ID > 0 {
			return content
		}

		if len(menusData.Top.SubType) > 0 {

			dao.Orm().Model(&dao.Content{}).
				Where("ContentItemID=?", ContentItemID).
				Where("ContentSubTypeID=?", menusData.Top.SubType[0].Item.ID).
				Order("CreatedAt desc").Order("ID desc").First(&content)
			if content.ID > 0 {
				ContentSubTypeID = menusData.Top.SubType[0].Item.ID
				menusData.SetCurrentMenus(ContentItemID, ContentSubTypeID, ContentSubTypeChildID)
				return content
			}

			if len(menusData.Top.SubType[0].SubType) > 0 {

				dao.Orm().Model(&dao.Content{}).
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
			dao.Orm().Model(&dao.Content{}).
				Where("ContentItemID=?", ContentItemID).
				Where("ContentSubTypeID=?", ContentSubTypeChildID).
				Order("CreatedAt desc").Order("ID desc").First(&content)

			return content

		} else if ContentSubTypeID > 0 {
			dao.Orm().Model(&dao.Content{}).
				Where("ContentItemID=?", ContentItemID).
				Where("ContentSubTypeID=?", ContentSubTypeID).
				Order("CreatedAt desc").Order("ID desc").First(&content)
			if content.ID > 0 {
				//ContentSubTypeID =menusData.Top.SubType[0].Item.ID
				//ContentSubTypeChildID=menusData.Top.SubType[0].SubType[0].Item.ID
				return content
			}

			if len(menusData.Sub.SubType) > 0 {
				dao.Orm().Model(&dao.Content{}).
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

	}

	return content
}
func (service ContentService) FindContentListByTypeID(menusData *dao.MenusData, ContentItemID, ContentSubTypeID, ContentSubTypeChildID uint64, _Page int, _Limit int) result.Pager {

	var pager result.Pager

	if ContentItemID == 0 {
		glog.Trace("参数ContentItemID为0")
		return pager
	}

	if ContentSubTypeID == 0 && ContentSubTypeChildID == 0 {
		if len(menusData.List) > 0 {

		}
		db := dao.Orm().Model(&dao.Content{}).Where("ContentItemID=?", ContentItemID).
			Order("CreatedAt desc").Order("ID desc")
		return dao.Paging(db, _Page, _Limit, dao.Content{})
	} else {

		if ContentSubTypeChildID > 0 {
			db := dao.Orm().Model(&dao.Content{}).Where("ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeChildID).
				Order("CreatedAt desc").Order("ID desc")
			return dao.Paging(db, _Page, _Limit, dao.Content{})
		} else {

			ContentSubTypeIDList := service.GetContentSubTypeAllIDByID(ContentItemID, ContentSubTypeID)

			db := dao.Orm().Model(&dao.Content{}).
				Where("ContentItemID=? and ContentSubTypeID in (?)", ContentItemID, ContentSubTypeIDList).
				Order("CreatedAt desc").Order("ID desc")
			return dao.Paging(db, _Page, _Limit, dao.Content{})
		}
	}

}

func (service ContentService) FindContentListForLeftRight(ContentItemID, ContentSubTypeID, ContentSubTypeChildID uint64, ContentID uint64, ContentCreatedAt time.Time) [2]dao.Content {

	var contentList [2]dao.Content
	if ContentItemID == 0 {
		glog.Trace("参数ContentItemID为0")
		return contentList
	}

	var ContentSubTypeIDList []uint64
	if ContentSubTypeID == 0 && ContentSubTypeChildID == 0 {
		ContentSubTypeIDList = []uint64{}
	} else {
		if ContentSubTypeChildID > 0 {
			ContentSubTypeIDList = service.GetContentSubTypeAllIDByID(ContentItemID, uint64(ContentSubTypeChildID))
		} else {
			ContentSubTypeIDList = service.GetContentSubTypeAllIDByID(ContentItemID, uint64(ContentSubTypeID))
		}

	}

	ContentSubTypeIDListStr := make([]string, 0)
	for index := range ContentSubTypeIDList {
		ContentSubTypeIDListStr = append(ContentSubTypeIDListStr, strconv.FormatUint(ContentSubTypeIDList[index], 10))
	}

	var whereSql = ""

	if len(ContentSubTypeIDList) > 0 {
		whereSql = fmt.Sprintf("ContentItemID=%v and ContentSubTypeID in (%v)", ContentItemID, strings.Join(ContentSubTypeIDListStr, ","))
	} else {
		whereSql = fmt.Sprintf("ContentItemID=%v", ContentItemID)
	}

	var left dao.Content
	var right dao.Content
	err := dao.Orm().Raw(`SELECT * FROM Content  WHERE `+whereSql+` and ID<>? and CreatedAt>=? ORDER BY CreatedAt,ID limit 1`, ContentID, ContentCreatedAt).Scan(&left).Error
	glog.Error(err)
	err = dao.Orm().Raw(`SELECT * FROM Content  WHERE `+whereSql+` and ID<>? and CreatedAt<=? ORDER BY CreatedAt desc,ID desc limit 1`, ContentID, ContentCreatedAt).Scan(&right).Error
	glog.Error(err)

	return [2]dao.Content{left, right}
}

func (service ContentService) GetContentByContentItemID(ContentItemID uint64) *dao.Content {
	article := &dao.Content{}
	dao.Orm().Where("ContentItemID=? and ContentSubTypeID=?", ContentItemID, 0).First(article)
	//service.ChangeMap(dao.Orm(), ID, &dao.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentByContentItemIDAndTitle(ContentItemID uint64, Title string) *dao.Content {
	article := &dao.Content{}
	dao.Orm().Where("ContentItemID=? and Title=?", ContentItemID, Title).First(article)
	//service.ChangeMap(dao.Orm(), ID, &dao.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentByContentItemIDAndContentSubTypeID(ContentItemID, ContentSubTypeID uint64) *dao.Content {
	article := &dao.Content{}
	dao.Orm().Where("ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeID).First(article)
	//service.ChangeMap(dao.Orm(), ID, &dao.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentByID(ID uint64) *dao.Content {
	article := &dao.Content{}
	err := service.Get(dao.Orm(), ID, article) //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)
	//service.ChangeMap(dao.Orm(), ID, &dao.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ContentService) GetContentAndAddLook(context *gweb.Context, ArticleID uint64) *dao.Content {

	article := &dao.Content{}
	err := service.Get(dao.Orm(), ArticleID, article) //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)

	if context.Session.Attributes.Get(strconv.Itoa(int(ArticleID))) == nil {
		context.Session.Attributes.Put(strconv.Itoa(int(ArticleID)), "Look")
		service.ChangeMap(dao.Orm(), ArticleID, &dao.Content{}, map[string]interface{}{"Look": article.Look + 1})

		if LookArticle, have := context.Data["LookArticle"]; have {

			if context.Session.Attributes.Get(play.SessionUser) != nil {
				user := context.Session.Attributes.Get(play.SessionUser).(*dao.User)
				err := service.Journal.AddScoreJournal(dao.Orm(),
					user.ID,
					"看文章送积分", "看文章/"+strconv.Itoa(int(article.ID)),
					play.ScoreJournal_Type_Look_Article, int64(LookArticle.(float64)), dao.KV{Key: "ArticleID", Value: article.ID})
				glog.Error(err)
			}
		}

	}
	return article
}

func (service ContentService) HaveContentByTitle(ContentItemID, ContentSubTypeID uint64, Title string) bool {
	Orm := dao.Orm()
	_article := &dao.Content{}
	Orm.Where("ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeID).Where("Title=?", Title).First(_article)
	if _article.ID == 0 {
		return false
	} else {
		return true
	}

}
func (service ContentService) AddContent(article *dao.Content) *result.ActionResult {

	as := &result.ActionResult{}
	Orm := dao.Orm()

	if article.ContentItemID == 0 {
		as.Code = result.ActionFail
		as.Message = "必须指定ContentItemID"
		return as
	}

	contentItem := service.GetContentItemByID(article.ContentItemID)
	switch contentItem.Type {

	case dao.ContentTypeContent:
		if article.ContentSubTypeID == 0 {
			//return &result.ActionResult{Code: result.ActionFail, Message: fmt.Sprintf("%v内容请指定类型", contentItem.Type)}
		} else {
			contentSubType := service.GetContentSubTypeByID(article.ContentSubTypeID)

			if contentSubType.ID == 0 {
				return &result.ActionResult{Code: result.ActionFail, Message: fmt.Sprintf("无效的类别%v", contentSubType.ID)}
			}

			content := service.GetContentByContentItemIDAndContentSubTypeID(article.ContentItemID, article.ContentSubTypeID)
			if content.ID > 0 && article.ID != content.ID {
				return &result.ActionResult{Code: result.ActionFail, Message: fmt.Sprintf("添加的内容与原内容冲突")}
			}
		}

	}

	//_article := &dao.Article{}
	//err := Orm.Where("ContentItemID=? and ContentSubTypeID=?", article.ContentItemID, article.ContentSubTypeID).Where("Title=?", article.Title).First(_article).Error
	//if _article.ID != 0 && _article.ID != article.ID {
	if false {
		as.Code = result.ActionFail
		as.Message = "添加失败，存在相同的标题"
	} else {
		articleID := article.ID
		var err error
		if articleID == 0 {
			err = service.Save(Orm, article) //self.dao.AddArticle(Orm, article)
		} else {
			err = service.ChangeMap(Orm, articleID, &dao.Content{}, map[string]interface{}{
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
			as.Code = result.ActionFail
			as.Message = err.Error()
		} else {
			as.Code = result.ActionOK
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
