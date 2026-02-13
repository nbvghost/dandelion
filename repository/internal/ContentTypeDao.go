package internal

import (
	"context"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/db"
)

type ContentTypeDao struct{}

func (ContentTypeDao) ListContentType(ctx context.Context) []dao.IEntity {
	Orm := db.GetDB(ctx)
	return dao.Find(Orm, &model.ContentType{}).List()
}
func (ContentTypeDao) ListContentTypeByType(ctx context.Context, Type string) *model.ContentType {
	//Orm := singleton.Orm()
	//company := context.Session.Attributes.Get(play.SessionOrganization).(*model.Organization)
	//var list model.ContentType
	item := dao.GetBy(db.GetDB(ctx), &model.ContentType{}, map[string]any{"Type": Type}).(*model.ContentType) //service.FindWhere(Orm, &list, "Type=?", Type)
	return item
}
