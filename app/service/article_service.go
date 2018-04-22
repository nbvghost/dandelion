package service

import (
	"dandelion/app/service/dao"

	"github.com/nbvghost/gweb/tool"
)

type ArticleService struct {
	dao dao.ArticleDao
}

func (self ArticleService) ChangeArticle(article *dao.Article) error {
	return self.dao.ChangeArticle(Orm, article)
}
func (self ArticleService) GetArticle(ID uint64) *dao.Article {
	return self.dao.GetArticle(Orm, ID)
}
func (self ArticleService) GetArticleByTitle(Title string) *dao.Article {
	return self.dao.GetArticleByTitle(Orm, Title)
}
func (self ArticleService) DelArticle(ID uint64) error {
	return self.dao.DelArticle(Orm, ID)
}
func (self ArticleService) FindArticleByCategoryID(ID uint64) []dao.Article {
	return self.dao.FindArticleByCategoryID(Orm, ID)
}
func (self ArticleService) AddArticle(article *dao.Article) *dao.ActionStatus {

	as := &dao.ActionStatus{}

	_article := &dao.Article{}
	err := Orm.Where("Title=?", article.Title).First(_article).Error
	tool.CheckError(err)
	if _article.ID != 0 {
		as.Success = false
		as.Message = "添加失败，存在相同的标题"
	} else {
		err = self.dao.AddArticle(Orm, article)
		tool.CheckError(err)
		as.Success = true
		as.Message = "添加成功"
	}
	return as
}
