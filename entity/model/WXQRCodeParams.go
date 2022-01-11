package model

import "github.com/nbvghost/dandelion/entity/base"

type WXQRCodeParams struct {
	base.BaseModel
	CodeKey string `gorm:"column:CodeKey;not null;unique"`
	Params  string `gorm:"column:Params;not null"`
}

func (WXQRCodeParams) TableName() string {
	return "WXQRCodeParams"
}
