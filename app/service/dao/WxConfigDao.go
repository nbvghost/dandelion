package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb/tool"
)

type WxConfigDao struct {
}

func (b WxConfigDao) GetWxConfig(DB *gorm.DB, CompanyID uint64) *WxConfig {
	item := &WxConfig{}
	err := DB.Where("CompanyID=?", CompanyID).First(item).Error
	tool.CheckError(err)

	if item.ID == 0 {
		item.CompanyID = CompanyID
		err = DB.Create(item).Error
		tool.CheckError(err)
		return item
	} else {
		return item
	}
}
func (b WxConfigDao) ChangeWxConfig(DB *gorm.DB, ID uint64, Value WxConfig) error {

	//item := b.GetWxConfig(DB, CompanyID)
	//item.V = Value
	return DB.Model(&WxConfig{}).Where("ID=?", ID).Updates(Value).Error
}
