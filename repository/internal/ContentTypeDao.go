package internal

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type ContentTypeDao struct{}

func (ContentTypeDao) ListContentType() []dao.IEntity {
	Orm := db.Orm()
	return dao.Find(Orm, &model.ContentType{}).List()
}
func (ContentTypeDao) ListContentTypeByType(Type string) *model.ContentType {
	//Orm := singleton.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//var list model.ContentType
	item := dao.GetBy(db.Orm(), &model.ContentType{}, map[string]any{"Type": Type}).(*model.ContentType) //service.FindWhere(Orm, &list, "Type=?", Type)
	return item
}
