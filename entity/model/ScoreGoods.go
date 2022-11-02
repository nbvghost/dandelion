package model

import (
	"errors"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/nbvghost/gpa/types"
)

//积分兑换产品
type ScoreGoods struct {
	types.Entity
	OID       types.PrimaryKey `gorm:"column:OID"`
	Name      string           `gorm:"column:Name"`
	Score     int              `gorm:"column:Score"`
	Price     uint             `gorm:"column:Price"`
	Image     string           `gorm:"column:Image"`
	Introduce string           `gorm:"column:Introduce"`
}

func (u *ScoreGoods) BeforeCreate(scope *gorm.DB) (err error) {
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
func (ScoreGoods) TableName() string {
	return "ScoreGoods"
}
