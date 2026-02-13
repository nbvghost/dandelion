package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
	"gorm.io/gorm"
)

type ContentDao struct{}

func (ContentDao) GetContentByID(ctx context.Context, ID dao.PrimaryKey) *model.Content {
	article := dao.GetByPrimaryKey(db.GetDB(ctx), &model.Content{}, ID).(*model.Content) //SelectOne(user, "select * from User where Email=?", Email)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (ContentDao) HotViewList(ctx context.Context, OID, ContentItemID dao.PrimaryKey, count uint) []model.Content {
	Orm := db.GetDB(ctx)
	var result []model.Content
	d := Orm.Model(&model.Content{}).Where(map[string]interface{}{"OID": OID}).Where(`"ContentItemID"=?`, ContentItemID).Order(`"CountView" desc`).Limit(int(count))
	d.Find(&result)
	return result
}
func (ContentDao) HotLikeList(ctx context.Context, OID, ContentItemID dao.PrimaryKey, count uint) []model.Content {
	Orm := db.GetDB(ctx)
	var result []model.Content
	d := Orm.Model(&model.Content{}).Where(map[string]interface{}{"OID": OID}).Where(`"ContentItemID"=?`, ContentItemID).Order(`"CountLike" desc`).Limit(int(count))
	d.Find(&result)
	return result
}
func (ContentDao) SortList(ctx context.Context, OID, ContentItemID dao.PrimaryKey, sort string, sortMethod int, count uint) []model.Content {
	Orm := db.GetDB(ctx)
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

func (ContentDao) ChangeContentByField(ctx context.Context, id dao.PrimaryKey, field string, value any) error {

	return dao.UpdateByPrimaryKey(db.GetDB(ctx), &model.Content{}, id, map[string]any{field: value})
}
func (ContentDao) ChangeContent(ctx context.Context, article *model.Content) error {

	return dao.Save(db.GetDB(ctx), article)
}

func (ContentDao) GetContentByTitle(Orm *gorm.DB, OID dao.PrimaryKey, Title string) *model.Content {
	article := &model.Content{}
	err := Orm.Where(`"OID"=?`, OID).Where(`"Title"=?`, Title).First(article).Error //SelectOne(user, "select * from User where Email=?", Email)
	if err != nil {
		log.Println(err)
	}
	return article
}
func (ContentDao) DelContent(ctx context.Context, ID dao.PrimaryKey) error {
	err := dao.DeleteByPrimaryKey(db.GetDB(ctx), &model.Content{}, ID)
	return err
}
func (ContentDao) FindContentByContentSubTypeID(ctx context.Context, ContentSubTypeID dao.PrimaryKey) []dao.IEntity {
	//var contentList []model.Content
	//err := service.FindWhere(singleton.Orm(), &contentList, "ContentSubTypeID=?", ContentSubTypeID) //SelectOne(user, "select * from User where Email=?", Email)

	contentList := dao.Find(db.GetDB(ctx), &model.Content{}).Where(`"ContentSubTypeID"=?`, ContentSubTypeID).List()

	return contentList
}
func (ContentDao) FindContentByContentItemIDAndContentSubTypeID(ctx context.Context, ContentItemID dao.PrimaryKey, ContentSubTypeID dao.PrimaryKey) *model.Content {
	//service.FindWhere(singleton.Orm(), &content, "ContentItemID=? and ContentSubTypeID=?", ContentItemID, ContentSubTypeID) //SelectOne(user, "select * from User where Email=?", Email)
	content := dao.GetBy(db.GetDB(ctx), &model.Content{}, map[string]any{"ContentItemID": ContentItemID, "ContentSubTypeID": ContentSubTypeID}).(*model.Content)
	return content
}
func (m ContentDao) GetContentByContentItemID(ctx context.Context, ContentItemID dao.PrimaryKey) *model.Content {
	//article := &model.Content{}
	//db.GetDB(ctx).Where(map[string]interface{}{"ContentItemID": ContentItemID, "ContentSubTypeID": 0}).First(article)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return m.FindContentByContentItemIDAndContentSubTypeID(ctx, ContentItemID, 0)
}
func (ContentDao) HaveContentByTitle(ctx context.Context, ContentItemID, ContentSubTypeID uint, Title string) bool {
	Orm := db.GetDB(ctx)
	_article := &model.Content{}
	Orm.Where(`"ContentItemID"=? and "ContentSubTypeID"=?`, ContentItemID, ContentSubTypeID).Where(`"Title"=?`, Title).First(_article)
	if _article.ID == 0 {
		return false
	} else {
		return true
	}
}
func (ContentDao) FindContentByIDAndNum(ctx context.Context, contentItemIDList []dao.PrimaryKey, num int) []model.Content {
	Orm := db.GetDB(ctx)
	_articleList := make([]model.Content, 0)
	Orm.Where(`"ContentItemID" in ?`, contentItemIDList).Order(`"CreatedAt" desc`).Limit(num).Find(&_articleList)
	return _articleList
}
func (ContentDao) GetContentByUri(ctx context.Context, OID dao.PrimaryKey, uri string) *model.Content {
	Orm := db.GetDB(ctx)
	var item model.Content
	Orm.Model(model.Content{}).Where(map[string]interface{}{"OID": OID, "Uri": uri}).First(&item)
	return &item
}

func (ContentDao) GetContentByContentItemIDAndTitle(ctx context.Context, ContentItemID uint, Title string) *model.Content {
	article := &model.Content{}
	db.GetDB(ctx).Where(map[string]interface{}{
		"ContentItemID": ContentItemID,
		"Title":         Title,
	}).First(article)
	//service.ChangeMap(singleton.Orm(), ID, &model.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}

func (ContentDao) FindContentByFieldGroupID(ctx context.Context, oid dao.PrimaryKey, fieldGroupID dao.PrimaryKey) []*model.Content {
	var list []*model.Content
	db.GetDB(ctx).Model(model.Content{}).Where(`"OID"=? and "FieldGroupID"=?`, oid, fieldGroupID).Find(&list)
	return list
}
