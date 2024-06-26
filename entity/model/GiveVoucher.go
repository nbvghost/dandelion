package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/library/dao"
)

// GiveVoucher 条件增送优惠卷
type GiveVoucher struct {
	dao.Entity
	OID           uint           `gorm:"column:OID"`
	ScoreMaxValue uint           `gorm:"column:ScoreMaxValue"`
	VoucherID     dao.PrimaryKey `gorm:"column:VoucherID"`
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
