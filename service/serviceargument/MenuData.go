package serviceargument

import (
	"github.com/nbvghost/dandelion/entity/extends"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type CurrentMenuData struct {
	Menus      extends.Menus
	TypeID     dao.PrimaryKey
	TypeUri    string
	SubTypeID  dao.PrimaryKey
	SubTypeUri string
}

func NewProductMenusData(goodsType model.GoodsType, goodsTypeChild model.GoodsTypeChild) CurrentMenuData {
	md := CurrentMenuData{}
	md.TypeID = goodsType.ID
	md.SubTypeID = goodsTypeChild.ID
	md.TypeUri = goodsType.Uri
	md.SubTypeUri = goodsTypeChild.Uri
	return md
}
func NewMenusData(item model.ContentItem, subType model.ContentSubType) CurrentMenuData {
	md := CurrentMenuData{}
	md.TypeID = item.ID
	md.SubTypeID = subType.ID
	md.TypeUri = item.Uri
	md.SubTypeUri = subType.Uri
	return md
}
