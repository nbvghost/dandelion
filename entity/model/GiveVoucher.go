package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/gpa/types"
)

//条件增送优惠卷
type GiveVoucher struct {
	base.BaseModel
	OID           uint             `gorm:"column:OID"`
	ScoreMaxValue uint             `gorm:"column:ScoreMaxValue"`
	VoucherID     types.PrimaryKey `gorm:"column:VoucherID"`
}

func (u *GiveVoucher) BeforeCreate(scope *gorm.DB) (err error) {
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
func (GiveVoucher) TableName() string {
	return "GiveVoucher"
}
