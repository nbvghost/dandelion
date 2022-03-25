package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

//优惠券
type Voucher struct {
	base.BaseModel
	OID       types.PrimaryKey `gorm:"column:OID"`
	Name      string           `gorm:"column:Name"`
	Amount    uint             `gorm:"column:Amount"`
	Image     string           `gorm:"column:Image"`
	UseDay    int              `gorm:"column:UseDay"`
	Introduce string           `gorm:"column:Introduce"`
}

func (u *Voucher) BeforeCreate(scope *gorm.DB) (err error) {
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
func (Voucher) TableName() string {
	return "Voucher"
}
