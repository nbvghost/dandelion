package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/gpa/types"
)

//满减
type FullCut struct {
	types.Entity
	OID       types.PrimaryKey `gorm:"column:OID"`
	Amount    uint             `gorm:"column:Amount"`
	CutAmount uint             `gorm:"column:CutAmount"`
}

func (u *FullCut) BeforeCreate(scope *gorm.DB) (err error) {
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
func (FullCut) TableName() string {
	return "FullCut"
}
