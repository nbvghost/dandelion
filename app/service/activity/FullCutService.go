package activity

import (
	"github.com/jinzhu/gorm"
	"github.com/nbvghost/dandelion/app/service/dao"
)

type FullCutService struct {
	dao.BaseDao
}

func (service FullCutService) FindOrderByAmountDesc(DB *gorm.DB, OID uint64) []dao.FullCut {
	var list []dao.FullCut
	DB.Model(&dao.FullCut{}).Where("OID=?", OID).Order("Amount desc").Find(&list)
	return list
}
func (service FullCutService) FindOrderByAmountASC(DB *gorm.DB, OID uint64) []dao.FullCut {
	var list []dao.FullCut
	DB.Model(&dao.FullCut{}).Where("OID=?", OID).Order("Amount asc").Find(&list)
	return list
}
