package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb/tool"
)

type ArticleDao struct{}

func (ArticleDao) ChangeArticle(Orm *gorm.DB, article *Article) error {
	return Orm.Save(article).Error
}
func (ArticleDao) GetArticle(Orm *gorm.DB, ID uint64) *Article {
	article := &Article{}
	err := Orm.Where("ID=?", ID).First(article).Error //SelectOne(user, "select * from User where Email=?", Email)
	tool.CheckError(err)
	return article
}
func (ArticleDao) GetArticleByTitle(Orm *gorm.DB, Title string) *Article {
	article := &Article{}
	err := Orm.Where("Title=?", Title).First(article).Error //SelectOne(user, "select * from User where Email=?", Email)
	tool.CheckError(err)
	return article
}
func (ArticleDao) DelArticle(Orm *gorm.DB, ID uint64) error {
	err := Orm.Where("ID=?", ID).Delete(Article{}).Error
	return err
}
func (ArticleDao) FindArticleByCategoryID(Orm *gorm.DB, ID uint64) []Article {
	articles := []Article{}
	err := Orm.Where("CategoryID=?", ID).Find(&articles).Error //SelectOne(user, "select * from User where Email=?", Email)
	tool.CheckError(err)
	return articles
}
func (ArticleDao) AddArticle(Orm *gorm.DB, article *Article) error {

	err := Orm.Create(article).Error
	tool.CheckError(err)
	return err
}
