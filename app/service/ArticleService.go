package service

import (
	"dandelion/app/service/dao"

	"dandelion/app/play"
	"dandelion/app/util"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gweb"
)

type ArticleService struct {
	dao.BaseDao
	Content ContentService
	Journal JournalService
}

func (service ArticleService) AddSpiderArticle(OID uint64, ContentName string, ContentSubTypeName string, Author, Title string, FromUrl string, Introduce string, Thumbnail string, Content string, CreatedAt time.Time) {
	var article dao.Article
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

	//Thumbnail=tool.DownloadInternetImage(Thumbnail,"Mozilla/5.0 (Linux; Android 7.0; SLA-AL00 Build/HUAWEISLA-AL00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/6.2 TBS/044109 Mobile Safari/537.36 MicroMessenger/6.6.7.1321(0x26060739) NetType/WIFI Language/zh_CN",weixin_tmp_url)
	article.Thumbnail = Thumbnail
	article.Content = Content

	contentType := service.Content.ListContentTypeByType(play.ContentTypeArticles)

	content := service.Content.GetContentByNameAndOID(ContentName, OID)
	if content.ID == 0 {
		content.OID = OID
		content.Type = contentType.Type
		content.Name = ContentName
		content.ContentTypeID = contentType.ID
		service.Content.Save(dao.Orm(), &content)

	}

	article.ContentID = content.ID
	contentSubType := service.Content.FindContentSubTypesByNameAndContentID(ContentSubTypeName, content.ID)
	if contentSubType.ID == 0 {
		contentSubType.Name = ContentSubTypeName
		contentSubType.ContentID = content.ID
		service.Content.Save(dao.Orm(), &contentSubType)
	}

	article.Author = Author
	article.ContentSubTypeID = contentSubType.ID
	service.AddArticle(&article)

}
func (service ArticleService) DataTablesAction(context *gweb.Context) gweb.Result {
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*dao.Organization)
	Orm := dao.Orm()
	dts := &dao.Datatables{}
	util.RequestBodyToJSON(context.Request.Body, dts)
	draw, recordsTotal, recordsFiltered, list := service.DatatablesListOrder(Orm, dts, &[]dao.Article{}, 0)
	return &gweb.JsonResult{Data: map[string]interface{}{"data": list, "draw": draw, "recordsTotal": recordsTotal, "recordsFiltered": recordsFiltered}}
}
func (service ArticleService) ChangeArticle(article *dao.Article) error {

	return service.Save(dao.Orm(), article)
}

func (service ArticleService) GetArticleByTitle(Orm *gorm.DB, Title string) *dao.Article {
	article := &dao.Article{}
	err := Orm.Where("Title=?", Title).First(article).Error //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)
	return article
}
func (service ArticleService) DelArticle(ID uint64) error {
	err := service.Delete(dao.Orm(), &dao.Article{}, ID)
	return err
}
func (service ArticleService) FindArticleByContentSubTypeID(ContentSubTypeID uint64) []dao.Article {
	var articles []dao.Article

	err := service.FindWhere(dao.Orm(), &articles, "ContentSubTypeID=?", ContentSubTypeID) //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)
	return articles
}
func (service ArticleService) DeleteArticleAction(context *gweb.Context) gweb.Result {

	context.Request.ParseForm()
	fmt.Println(context.Request.FormValue("ID"))
	ID, _ := strconv.ParseUint(context.Request.FormValue("ID"), 10, 64)
	err := service.Delete(dao.Orm(), &dao.Article{}, ID)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "删除成功", nil)}
}
func (service ArticleService) GetArticleAction(context *gweb.Context) gweb.Result {
	ID, _ := strconv.ParseUint(context.PathParams["ID"], 10, 64)
	var article dao.Article
	err := service.Get(dao.Orm(), ID, &article)
	return &gweb.JsonResult{Data: (&dao.ActionStatus{}).SmartError(err, "OK", article)}
}
func (service ArticleService) GetArticle(ID uint64) *dao.Article {
	article := &dao.Article{}
	err := service.Get(dao.Orm(), ID, article) //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)
	//service.ChangeMap(dao.Orm(), ID, &dao.Article{}, map[string]interface{}{"Look": article.Look + 1})
	return article
}
func (service ArticleService) GetArticleAndAddLook(context *gweb.Context, ArticleID uint64) *dao.Article {

	article := &dao.Article{}
	err := service.Get(dao.Orm(), ArticleID, article) //SelectOne(user, "select * from User where Email=?", Email)
	glog.Error(err)

	if context.Session.Attributes.Get(strconv.Itoa(int(ArticleID))) == nil {
		context.Session.Attributes.Put(strconv.Itoa(int(ArticleID)), "Look")
		service.ChangeMap(dao.Orm(), ArticleID, &dao.Article{}, map[string]interface{}{"Look": article.Look + 1})

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
func (service ArticleService) SaveArticleAction(context *gweb.Context) gweb.Result {

	dts := &dao.Article{}
	util.RequestBodyToJSON(context.Request.Body, dts)

	fmt.Println(dts)

	as := service.AddArticle(dts)

	return &gweb.JsonResult{Data: as}
}
func (service ArticleService) HaveArticleByTitle(ContentID, ContentSubTypeID uint64, Title string) bool {
	Orm := dao.Orm()
	_article := &dao.Article{}
	Orm.Where("ContentID=? and ContentSubTypeID=?", ContentID, ContentSubTypeID).Where("Title=?", Title).First(_article)
	if _article.ID == 0 {
		return false
	} else {
		return true
	}

}
func (service ArticleService) AddArticle(article *dao.Article) *dao.ActionStatus {

	as := &dao.ActionStatus{}
	Orm := dao.Orm()

	_article := &dao.Article{}
	err := Orm.Where("ContentID=? and ContentSubTypeID=?", article.ContentID, article.ContentSubTypeID).Where("Title=?", article.Title).First(_article).Error
	if _article.ID != 0 && _article.ID != article.ID {
		as.Success = false
		as.Message = "添加失败，存在相同的标题"
	} else {
		//fmt.Println(article.Introduce)
		err = service.Save(Orm, article) //self.dao.AddArticle(Orm, article)
		if err != nil {
			glog.Error(err)
			as.Success = false
			as.Message = err.Error()
		} else {
			as.Success = true
			if article.ID != 0 {
				as.Message = "修改成功"
			} else {
				as.Message = "添加成功"
			}

		}
	}
	return as
}
