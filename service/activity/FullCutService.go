package activity

import (
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/gpa/types"
	"gorm.io/gorm"
)

type FullCutService struct {
	model.BaseDao
}

func (service FullCutService) FindOrderByAmountDesc(DB *gorm.DB, OID types.PrimaryKey) []model.FullCut {
	var list []model.FullCut
	DB.Model(&model.FullCut{}).Where(`"OID"=?`, OID).Order(`"Amount" desc`).Find(&list)
	return list
}
func (service FullCutService) FindOrderByAmountASC(DB *gorm.DB, OID types.PrimaryKey) []model.FullCut {
	var list []model.FullCut
	DB.Model(&model.FullCut{}).Where(`"OID"=?`, OID).Order(`"Amount" asc`).Find(&list)
	return list
}
