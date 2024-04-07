package internal

import (
	"fmt"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"gorm.io/gorm"
	"log"
)

type ContentDao struct{}



func (ContentDao) GetContentByID(ID dao.PrimaryKey) *model.Content {
	article := dao.GetByPrimaryKey(db.Orm(), &model.Content{}, ID).(*model.Content) //SelectOne(user, "select * from User where Email=?", Email)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (ContentDao) HotViewList(OID, ContentItemID dao.PrimaryKey, count uint) []model.Content {
	Orm := db.Orm()
	var result []model.Content
	d := Orm.Model(&model.Content{}).Where(map[string]interface{}{"OID": OID}).Where(`"ContentItemID"=?`, ContentItemID).Order(`"CountView" desc`).Limit(int(count))
	d.Find(&result)
	return result
}
func (ContentDao) HotLikeList(OID, ContentItemID dao.PrimaryKey, count uint) []model.Content {
	Orm := db.Orm()
	var result []model.Content
	d := Orm.Model(&model.Content{}).Where(map[string]interface{}{"OID": OID}).Where(`"ContentItemID"=?`, ContentItemID).Order(`"CountLike" desc`).Limit(int(count))
	d.Find(&result)
	return result
}
func (ContentDao) SortList(OID, ContentItemID dao.PrimaryKey, sort string, sortMethod int, count uint) []model.Content {
	Orm := db.Orm()
	var result []model.Content
	d := Orm.Model(&model.Content{}).Where(map[string]interface{}{"OID": OID}).Where(`"ContentItemID"=?`, ContentItemID)
	if sortMethod >= 0 {
		d = d.Order(fmt.Sprintf(`"%s" asc`, sort))
	} else {
		d = d.Order(fmt.Sprintf(`"%s" desc`, sort))
	}
	d = d.Order(`"CreatedAt" desc`)
	d.Limit(int(count)).Find(&result)
	return result
}

func (ContentDao) ChangeContentByField(id dao.PrimaryKey, field string, value any) error {

	return dao.UpdateByPrimaryKey(db.Orm(), &model.Content{}, id, map[string]any{field: value})
}
func (ContentDao) ChangeContent(article *model.Content) error {

	return dao.Save(db.Orm(), article)
}

func (ContentDao) GetContentByTitle(Orm *gorm.DB, OID dao.PrimaryKey, Title string) *model.Content {
	article := &model.Content{}
	err := Orm.Where(`"OID"=?`, OID).Where(`"Title"=?`, Title).First(article).Error //SelectOne(user, "select * from User where Email=?", Email)
	if err != nil {
		log.Println(err)
	}
	return article
}
func (ContentDao) DelContent(ID dao.PrimaryKey) error {
	err := dao.DeleteByPrimaryKey(db.Orm(), &model.Content{}, ID)
	return err
}
func (ContentDao) FindContentByContentSubTypeID(ContentSubTypeID dao.PrimaryKey) []dao.IEntity {
	//var contentList []model.Content
	//err := service.FindWhere(singleton.Orm(), &contentList, "ContentSubTypeID=?", ContentSubTypeID) //SelectOne(user, "select * from User where Email=?", Email)

	contentList := dao.Find(db.Orm(), &model.Content{}).Where(`"ContentSubTypeID"=?`, ContentSubTypeID).List()

	return contentList
}
func (ContentDao) FindContentByContentItemIDAndContentSubTypeID(ContentItemID dao.PrimaryKey, ContentSubTypeID dao.PrimaryKey) *model.Content {
	//service.FindWhere(singleton.Orm(), &content, "ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeID) //SelectOne(user, "select * from User where Email=?", Email)
	content := dao.GetBy(db.Orm(), &model.Content{}, map[string]any{"ContentItemID": ContentItemID, "ContentSubTypeID": ContentSubTypeID}).(*model.Content)
	return content
}
func (m ContentDao) GetContentByContentItemID(ContentItemID dao.PrimaryKey) *model.Content {
	//article := &model.Content{}
	//db.Orm().Where(map[string]interface{}{"ContentItemID": ContentItemID, "ContentSubTypeID": 0}).First(article)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return m.FindContentByContentItemIDAndContentSubTypeID(ContentItemID, 0)
}
func (ContentDao) HaveContentByTitle(ContentItemID, ContentSubTypeID uint, Title string) bool {
	Orm := db.Orm()
	_article := &model.Content{}
	Orm.Where(`"ContentItemID"=? and "ContentSubTypeID"=?`, ContentItemID, ContentSubTypeID).Where(`"Title"=?`, Title).First(_article)
	if _article.ID == 0 {
		return false
	} else {
		return true
	}
}
func (ContentDao) FindContentByIDAndNum(contentItemIDList []dao.PrimaryKey, num int) []model.Content {
	Orm := db.Orm()
	_articleList := make([]model.Content, 0)
	Orm.Where(`"ContentItemID" in ?`, contentItemIDList).Order(`"CreatedAt" desc`).Limit(num).Find(&_articleList)
	return _articleList
}
func (ContentDao) GetContentByUri(OID dao.PrimaryKey, uri string) *model.Content {
	Orm := db.Orm()
	var item model.Content
	Orm.Model(model.Content{}).Where(map[string]interface{}{"OID": OID, "Uri": uri}).First(&item)
	return &item
}

func (ContentDao) GetContentByContentItemIDAndTitle(ContentItemID uint, Title string) *model.Content {
	article := &model.Content{}
	db.Orm().Where(map[string]interface{}{
		"ContentItemID": ContentItemID,
		"Title":         Title,
	}).First(article)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}

func (ContentDao) FindContentByFieldGroupID(oid dao.PrimaryKey, fieldGroupID dao.PrimaryKey) []*model.Content {
	var list []*model.Content
	db.Orm().Model(model.Content{}).Where(`"OID"=? and "FieldGroupID"=?`, oid, fieldGroupID).Find(&list)
	return list
}
