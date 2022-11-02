package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/gpa/types"
)

type ExpressTemplate struct {
	types.Entity
	OID      types.PrimaryKey `gorm:"column:OID"`
	Name     string           `gorm:"column:Name"`
	Drawee   string           `gorm:"column:Drawee"`             //付款人
	Type     string           `gorm:"column:Type"`               //KG  ITEM
	Template string           `gorm:"column:Template;type:text"` //json
	Free     string           `gorm:"column:Free;type:text"`     //json []
}

func (u *ExpressTemplate) BeforeCreate(scope *gorm.DB) (err error) {
	if u.OID == 0 {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		panic(errors.New(u.TableName() + ":OID不能为空"))

	}
	return nil
}
func (u ExpressTemplate) TableName() string {
	return "ExpressTemplate"
}
