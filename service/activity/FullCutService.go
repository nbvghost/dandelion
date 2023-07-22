package activity

import (
	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
)

type FullCutService struct {
	model.BaseDao
}

func (service FullCutService) FindOrderByAmountDesc(DB *gorm.DB, OID dao.PrimaryKey) []model.FullCut {
	var list []model.FullCut
	DB.Model(&model.FullCut{}).Where(`"OID"=?`, OID).Order(`"Amount" desc`).Find(&list)
	return list
}
func (service FullCutService) FindOrderByAmountASC(DB *gorm.DB, OID dao.PrimaryKey) []model.FullCut {
	var list []model.FullCut
	DB.Model(&model.FullCut{}).Where(`"OID"=?`, OID).Order(`"Amount" asc`).Find(&list)
	return list
}
