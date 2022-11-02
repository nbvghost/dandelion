package model

import (
	"github.com/nbvghost/gpa/types"
)

type WXQRCodeParams struct {
	types.Entity
	CodeKey string `gorm:"column:CodeKey;not null;unique"`
	Params  string `gorm:"column:Params;not null"`
}

func (WXQRCodeParams) TableName() string {
	return "WXQRCodeParams"
}
