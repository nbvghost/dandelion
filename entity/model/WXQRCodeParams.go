package model

import (
	"github.com/nbvghost/dandelion/library/dao"
)

type WXQRCodeParams struct {
	dao.Entity
	CodeKey string `gorm:"column:CodeKey;not null;unique"`
	Params  string `gorm:"column:Params;not null"`
}

func (WXQRCodeParams) TableName() string {
	return "WXQRCodeParams"
}
